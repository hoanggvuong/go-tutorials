package main

import (
	"fmt"
	"os"
	// Dùng để tạo file
)

func main() {
	// Ghi dữ liệu vào file
	// dataToWrite := []byte("Đây là dòng chữ sẽ được ghi vào file.\nGo thật tuyệt vời!")
	// err := ioutil.WriteFile("myfile.txt", dataToWrite, 0644) // 0644 là quyền hạn file (có thể đọc/ghi bởi chủ sở hữu, chỉ đọc bởi nhóm/người khác)
	// if err != nil {
	// 	fmt.Println("Lỗi khi ghi file:", err)
	// 	return
	// }
	// fmt.Println("Đã ghi thành công vào myfile.txt")

	// // Đọc dữ liệu từ file
	// dataRead, err := ioutil.ReadFile("myfile.txt")
	// if err != nil {
	// 	fmt.Println("Lỗi khi đọc file:", err)
	// 	return
	// }
	// fmt.Println("Nội dung đọc từ myfile.txt:\n", string(dataRead))

	//Xóa file (tùy chọn)
	err := os.Remove("myfile.txt")
	if err != nil {
		fmt.Println("Lỗi khi xóa file:", err)
	} else {
		fmt.Println("Đã xóa myfile.txt")
	}
}
