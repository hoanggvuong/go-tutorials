package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/segmentio/kafka-go"
)

// Cấu hình Kafka
const (
	kafkaBroker = "localhost:9092"
	topic       = "test-topic"
)

// Cấu trúc dữ liệu cho tin nhắn
type Message struct {
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
}

var kafkaWriter *kafka.Writer

func init() {
	// Khởi tạo Kafka Writer một lần khi chương trình bắt đầu
	kafkaWriter = &kafka.Writer{
		Addr:         kafka.TCP(kafkaBroker),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		BatchTimeout: 5 * time.Millisecond,
	}
}

func main() {
	defer kafkaWriter.Close()

	// Tạo một topic nếu nó chưa tồn tại (chỉ dùng cho mục đích ví dụ)
	createTopic()

	// Định nghĩa route cho API
	http.HandleFunc("/publish", publishHandler)

	fmt.Println("Kafka API Gateway đang chạy trên cổng 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Hàm handler để xử lý yêu cầu POST
func publishHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Chỉ chấp nhận phương thức POST", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Không thể đọc body của request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var msg Message
	if err := json.Unmarshal(body, &msg); err != nil {
		http.Error(w, "Dữ liệu JSON không hợp lệ", http.StatusBadRequest)
		return
	}

	// Gửi tin nhắn đến Kafka
	err = kafkaWriter.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(msg.Key),
		Value: (msg.Value),
	})
	if err != nil {
		log.Printf("Lỗi khi gửi tin nhắn đến Kafka: %v", err)
		http.Error(w, "Lỗi nội bộ khi gửi tin nhắn", http.StatusInternalServerError)
		return
	}

	// Trả về phản hồi thành công
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"status": "success", "message": "Tin nhắn đã được gửi đến Kafka"}
	json.NewEncoder(w).Encode(response)
}

// Hàm tiện ích để tạo topic
func createTopic() {
	conn, err := kafka.Dial("tcp", kafkaBroker)
	if err != nil {
		log.Fatalf("Không thể kết nối đến Kafka để tạo topic: %v", err)
	}
	defer conn.Close()

	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     3,
		ReplicationFactor: 1,
	}

	// Tạo topic nếu nó chưa tồn tại
	if err := conn.CreateTopics(topicConfig); err != nil {
		log.Printf("Cảnh báo: Không thể tạo topic. Lỗi: %v", err)
	} else {
		fmt.Printf("Đã tạo thành công topic '%s' với %d partitions.\n", topic, topicConfig.NumPartitions)
	}
}
