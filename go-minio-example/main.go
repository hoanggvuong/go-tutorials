package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"         // Package chung của AWS SDK
	"github.com/aws/aws-sdk-go-v2/config"      // Package để tải cấu hình mặc định hoặc tùy chỉnh
	"github.com/aws/aws-sdk-go-v2/credentials" // Package để cung cấp thông tin xác thực
	"github.com/aws/aws-sdk-go-v2/service/s3"  // Package dành riêng cho dịch vụ S3
	// Không cần import "github.com/aws/aws-sdk-go-v2/aws/endpoints" nữa
	// vì EndpointResolverOptions đã được tích hợp vào config/aws.
)

// Các biến cấu hình MinIO của bạn
// THAY THẾ CÁC GIÁ TRỊ NÀY VỚI THÔNG TIN MINIO CỦA BẠN!
const (
	minioEndpoint   = "minio.example.com"      // Hoặc https://
	minioAccessKey  = "Y3APLQ"                    // Thay thế bằng Access Key của bạn
	minioSecretKey  = "BF6BOUqvzf9rP1" // Thay thế bằng Secret Key của bạn
	minioBucketName = "my-go-bucket"                             // Tên bucket chúng ta sẽ tạo/sử dụng
	minioRegion     = "hn-cmc"                                   // Vùng ảo, MinIO thường không quan tâm nhiều đến vùng
)

func main() {
	// 1. Tạo một EndpointResolver tùy chỉnh để trỏ đến MinIO
	// Các tùy chọn xử lý endpoint giờ được cung cấp trực tiếp trong config.LoadDefaultAWSConfig
	customResolver := func(svc, region string, opt ...interface{}) (config.EndpointResolver, error) {
		if svc == s3.ServiceID {
			return config.EndpointResolverWithOptionsFunc(func(options ...interface{}) (aws.Endpoint, error) {
				ep, ok := options[0].(aws.Endpoint)
				if !ok {
					return aws.Endpoint{}, fmt.Errorf("invalid endpoint option")
				}
				return ep, nil
			}), nil
		}
		return nil, fmt.Errorf("unsupported service: %s", svc)
	}

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(minioAccessKey, minioSecretKey, "")),
		config.WithEndpointResolver(customResolver),
		config.WithRegion(minioRegion),
	)
	if err != nil {
		log.Fatalf("Không thể tải cấu hình AWS SDK: %v", err)
	}

	// 3. Tạo S3 client
	s3Client := s3.NewFromConfig(cfg)
	fmt.Println("Đã kết nối thành công tới MinIO.")

	// --- Bắt đầu các thao tác MinIO ---

	// A. Liệt kê các Bucket hiện có (để kiểm tra kết nối)
	fmt.Println("\n--- Liệt kê các Buckets ---")
	listBucketsOutput, err := s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatalf("Lỗi khi liệt kê buckets: %v", err)
	}
	if len(listBucketsOutput.Buckets) == 0 {
		fmt.Println("Không có bucket nào được tìm thấy.")
	} else {
		for _, bucket := range listBucketsOutput.Buckets {
			fmt.Printf("- %s (Tạo lúc: %s)\n", *bucket.Name, bucket.CreationDate.Format("2006-01-02 15:04:05"))
		}
	}

	// B. Tạo một Bucket mới
	fmt.Println("\n--- Tạo Bucket mới ---")
	_, err = s3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(minioBucketName),
		// Nếu MinIO của bạn hỗ trợ region, bạn có thể thêm:
		// CreateBucketConfiguration: &types.CreateBucketConfiguration{
		// 	LocationConstraint: types.BucketLocationConstraint(minioRegion),
		// },
	})
	if err != nil {
		if strings.Contains(err.Error(), "BucketAlreadyOwnedByYou") || strings.Contains(err.Error(), "BucketAlreadyExists") {
			fmt.Printf("Bucket '%s' đã tồn tại hoặc bạn đã sở hữu nó.\n", minioBucketName)
		} else {
			log.Fatalf("Lỗi khi tạo bucket '%s': %v", minioBucketName, err)
		}
	} else {
		fmt.Printf("Đã tạo thành công bucket: '%s'\n", minioBucketName)
	}

	// C. Tải lên một Đối tượng (File)
	fmt.Println("\n--- Tải lên Đối tượng (File) ---")
	objectKey := "my-first-go-object.txt"
	fileContent := "Hello from Go and MinIO! This is my first object."
	uploadInput := &s3.PutObjectInput{
		Bucket: aws.String(minioBucketName),
		Key:    aws.String(objectKey),
		Body:   strings.NewReader(fileContent), // Đọc từ một string
		// Bạn có thể thêm ContentType nếu biết:
		// ContentType: aws.String("text/plain"),
	}
	_, err = s3Client.PutObject(context.TODO(), uploadInput)
	if err != nil {
		log.Fatalf("Lỗi khi tải lên đối tượng '%s': %v", objectKey, err)
	}
	fmt.Printf("Đã tải lên thành công đối tượng '%s' vào bucket '%s'.\n", objectKey, minioBucketName)

	// D. Liệt kê các Đối tượng trong Bucket
	fmt.Println("\n--- Liệt kê Đối tượng trong Bucket ---")
	listObjectsInput := &s3.ListObjectsV2Input{
		Bucket: aws.String(minioBucketName),
	}
	listObjectsOutput, err := s3Client.ListObjectsV2(context.TODO(), listObjectsInput)
	if err != nil {
		log.Fatalf("Lỗi khi liệt kê đối tượng trong bucket '%s': %v", minioBucketName, err)
	}
	if len(listObjectsOutput.Contents) == 0 {
		fmt.Println("Không có đối tượng nào trong bucket.")
	} else {
		for _, object := range listObjectsOutput.Contents {
			fmt.Printf("- %s (Kích thước: %d bytes, LastModified: %s)\n", *object.Key, object.Size, object.LastModified.Format("2006-01-02 15:04:05"))
		}
	}

	// E. Tải xuống một Đối tượng (File)
	fmt.Println("\n--- Tải xuống Đối tượng (File) ---")
	downloadInput := &s3.GetObjectInput{
		Bucket: aws.String(minioBucketName),
		Key:    aws.String(objectKey),
	}
	downloadOutput, err := s3Client.GetObject(context.TODO(), downloadInput)
	if err != nil {
		log.Fatalf("Lỗi khi tải xuống đối tượng '%s': %v", objectKey, err)
	}
	defer downloadOutput.Body.Close() // Đảm bảo đóng body stream

	// Đọc nội dung file
	downloadedContent, err := os.ReadAll(downloadOutput.Body)
	if err != nil {
		log.Fatalf("Lỗi khi đọc nội dung tải xuống: %v", err)
	}
	fmt.Printf("Nội dung tải xuống từ '%s':\n%s\n", objectKey, string(downloadedContent))

	// F. Xóa một Đối tượng (File)
	fmt.Println("\n--- Xóa Đối tượng (File) ---")
	deleteObjectInput := &s3.DeleteObjectInput{
		Bucket: aws.String(minioBucketName),
		Key:    aws.String(objectKey),
	}
	_, err = s3Client.DeleteObject(context.TODO(), deleteObjectInput)
	if err != nil {
		log.Fatalf("Lỗi khi xóa đối tượng '%s': %v", objectKey, err)
	}
	fmt.Printf("Đã xóa thành công đối tượng '%s'.\n", objectKey)

	fmt.Println("\n--- Các thao tác MinIO đã hoàn tất. ---")
}
func main() {
	// 1. Tạo một custom resolver để trỏ đến endpoint MinIO của bạn
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == s3.ServiceID {
			return aws.Endpoint{
				URL:           minioEndpoint,
				SigningRegion: minioRegion,
				Source:        aws.EndpointSourceCustom,
				// Tùy chọn nếu MinIO của bạn dùng HTTP (không HTTPS) hoặc chứng chỉ tự ký
				// Insecure: true,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	// 2. Tải cấu hình AWS SDK
	cfg, err := config.LoadDefaultAWSConfig(context.TODO(),
		// Cung cấp credentials MinIO
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(minioAccessKey, minioSecretKey, "")),
		// Sử dụng custom resolver cho MinIO endpoint
		config.WithEndpointResolverWithOptions(customResolver),
		// Đặt region
		config.WithRegion(minioRegion),
		// Tùy chọn nếu MinIO của bạn dùng HTTP (không HTTPS) hoặc chứng chỉ tự ký
		// config.WithHTTPClient(&http.Client{Transport: &http.Transport{
		// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		// }}),
	)
	if err != nil {
		log.Fatalf("Không thể tải cấu hình AWS SDK: %v", err)
	}

	// 3. Tạo S3 client
	s3Client := s3.NewFromConfig(cfg)
	fmt.Println("Đã kết nối thành công tới MinIO.")

	// --- Bắt đầu các thao tác MinIO ---

	// A. Liệt kê các Bucket hiện có (để kiểm tra kết nối)
	fmt.Println("\n--- Liệt kê các Buckets ---")
	listBucketsOutput, err := s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatalf("Lỗi khi liệt kê buckets: %v", err)
	}
	if len(listBucketsOutput.Buckets) == 0 {
		fmt.Println("Không có bucket nào được tìm thấy.")
	} else {
		for _, bucket := range listBucketsOutput.Buckets {
			fmt.Printf("- %s (Tạo lúc: %s)\n", *bucket.Name, bucket.CreationDate.Format("2006-01-02 15:04:05"))
		}
	}

	// B. Tạo một Bucket mới
	fmt.Println("\n--- Tạo Bucket mới ---")
	_, err = s3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(minioBucketName),
	})
	if err != nil {
		if strings.Contains(err.Error(), "BucketAlreadyOwnedByYou") || strings.Contains(err.Error(), "BucketAlreadyExists") {
			fmt.Printf("Bucket '%s' đã tồn tại hoặc bạn đã sở hữu nó.\n", minioBucketName)
		} else {
			log.Fatalf("Lỗi khi tạo bucket '%s': %v", minioBucketName, err)
		}
	} else {
		fmt.Printf("Đã tạo thành công bucket: '%s'\n", minioBucketName)
	}

	// C. Tải lên một Đối tượng (File)
	fmt.Println("\n--- Tải lên Đối tượng (File) ---")
	objectKey := "my-first-go-object.txt"
	fileContent := "Hello from Go and MinIO! This is my first object."
	uploadInput := &s3.PutObjectInput{
		Bucket: aws.String(minioBucketName),
		Key:    aws.String(objectKey),
		Body:   strings.NewReader(fileContent), // Đọc từ một string
	}
	_, err = s3Client.PutObject(context.TODO(), uploadInput)
	if err != nil {
		log.Fatalf("Lỗi khi tải lên đối tượng '%s': %v", objectKey, err)
	}
	fmt.Printf("Đã tải lên thành công đối tượng '%s' vào bucket '%s'.\n", objectKey, minioBucketName)

	// D. Liệt kê các Đối tượng trong Bucket
	fmt.Println("\n--- Liệt kê Đối tượng trong Bucket ---")
	listObjectsInput := &s3.ListObjectsV2Input{
		Bucket: aws.String(minioBucketName),
	}
	listObjectsOutput, err := s3Client.ListObjectsV2(context.TODO(), listObjectsInput)
	if err != nil {
		log.Fatalf("Lỗi khi liệt kê đối tượng trong bucket '%s': %v", minioBucketName, err)
	}
	if len(listObjectsOutput.Contents) == 0 {
		fmt.Println("Không có đối tượng nào trong bucket.")
	} else {
		for _, object := range listObjectsOutput.Contents {
			fmt.Printf("- %s (Kích thước: %d bytes, LastModified: %s)\n", *object.Key, object.Size, object.LastModified.Format("2006-01-02 15:04:05"))
		}
	}

	// E. Tải xuống một Đối tượng (File)
	fmt.Println("\n--- Tải xuống Đối tượng (File) ---")
	downloadInput := &s3.GetObjectInput{
		Bucket: aws.String(minioBucketName),
		Key:    aws.String(objectKey),
	}
	downloadOutput, err := s3Client.GetObject(context.TODO(), downloadInput)
	if err != nil {
		log.Fatalf("Lỗi khi tải xuống đối tượng '%s': %v", objectKey, err)
	}
	defer downloadOutput.Body.Close() // Đảm bảo đóng body stream

	// Đọc nội dung file
	downloadedContent, err := os.ReadAll(downloadOutput.Body)
	if err != nil {
		log.Fatalf("Lỗi khi đọc nội dung tải xuống: %v", err)
	}
	fmt.Printf("Nội dung tải xuống từ '%s':\n%s\n", objectKey, string(downloadedContent))

	// F. Xóa một Đối tượng (File)
	fmt.Println("\n--- Xóa Đối tượng (File) ---")
	deleteObjectInput := &s3.DeleteObjectInput{
		Bucket: aws.String(minioBucketName),
		Key:    aws.String(objectKey),
	}
	_, err = s3Client.DeleteObject(context.TODO(), deleteObjectInput)
	if err != nil {
		log.Fatalf("Lỗi khi xóa đối tượng '%s': %v", objectKey, err)
	}
	fmt.Printf("Đã xóa thành công đối tượng '%s'.\n", objectKey)

	fmt.Println("\n--- Các thao tác MinIO đã hoàn tất. ---")
}
