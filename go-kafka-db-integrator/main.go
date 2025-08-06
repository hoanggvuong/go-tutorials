package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/segmentio/kafka-go"
)

// Cấu hình Kafka
const (
	kafkaBroker = "localhost:9092"
	topic       = "test-topic"
	groupID     = "db-integrator-group"
)

// Cấu hình Database
const (
	dbHost     = "localhost"
	dbPort     = 5432
	dbUser     = "user"
	dbPassword = "password"
	dbName     = "mydatabase"
)

// Cấu trúc dữ liệu cho tin nhắn/bản ghi đơn hàng
type OrderMessage struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

func main() {
	// 1. Kết nối đến database PostgreSQL
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Không thể kết nối đến database: %v", err)
	}
	defer db.Close()

	// Kiểm tra kết nối
	err = db.Ping()
	if err != nil {
		log.Fatalf("Không thể ping database: %v", err)
	}
	fmt.Println("Đã kết nối thành công đến PostgreSQL!")

	// 2. Tạo bảng `orders` nếu chưa tồn tại
	createTable(db)

	// 3. Khởi tạo Kafka Reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker},
		Topic:   topic,
		GroupID: groupID,
	})
	defer reader.Close()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Lắng nghe tín hiệu dừng chương trình
	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
		<-sigchan
		fmt.Println("\nNhận tín hiệu dừng, đang đóng Consumer...")
		cancel()
	}()

	fmt.Println("Bắt đầu lắng nghe tin nhắn từ Kafka...")

	// 4. Đọc tin nhắn và chèn vào database
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Consumer đã đóng.")
			return
		default:
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					continue
				}
				log.Fatalf("Lỗi khi đọc tin nhắn: %v", err)
			}

			// Giải mã tin nhắn JSON
			var orderMsg OrderMessage
			err = json.Unmarshal(msg.Value, &orderMsg)
			if err != nil {
				log.Printf("Bỏ qua tin nhắn không hợp lệ (không phải JSON): %v", err)
				continue
			}
			fmt.Printf("Consumer nhận tin nhắn: ID=%s, Value=%s\n", orderMsg.ID, orderMsg.Value)
			// Chèn dữ liệu vào database
			err = insertOrder(db, orderMsg)
			if err != nil {
				log.Printf("Lỗi khi chèn đơn hàng vào database: %v", err)
			} else {
				fmt.Printf("Đã lưu đơn hàng %s vào database.\n", orderMsg.ID)
			}
		}
	}
}

// Hàm tiện ích để tạo bảng
func createTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS orders (
		id VARCHAR(255) PRIMARY KEY,
		value TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Không thể tạo bảng 'orders': %v", err)
	}
	fmt.Println("Đã tạo bảng 'orders' (nếu chưa tồn tại).")
}

// Hàm tiện ích để chèn dữ liệu
func insertOrder(db *sql.DB, order OrderMessage) error {
	_, err := db.Exec("INSERT INTO orders (id, value) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING", order.ID, order.Value)
	return err
}
