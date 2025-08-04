package main

import (
	"encoding/json" // Để làm việc với JSON
	"fmt"
	"log"      // Để ghi log
	"net/http" // Để tạo HTTP server
)

// User là một struct đại diện cho thông tin người dùng
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Hàm handler cho endpoint /users/{id}
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	// Lấy ID người dùng từ URL (ví dụ: /users/123 -> ID là "123")
	// Trong một ứng dụng thực tế, bạn sẽ dùng thư viện router để dễ dàng hơn
	userID := r.URL.Path[len("/users/"):] // Lấy phần sau /users/

	// Giả lập dữ liệu người dùng (trong thực tế sẽ truy vấn database)
	var user User
	if userID == "1" {
		user = User{ID: "1", Name: "Alice", Email: "alice@example.com"}
	} else if userID == "2" {
		user = User{ID: "2", Name: "Bob", Email: "bob@example.com"}
	} else {
		// Trả về lỗi 404 nếu không tìm thấy người dùng
		http.NotFound(w, r)
		return
	}

	// Thiết lập Content-Type là application/json
	w.Header().Set("Content-Type", "application/json")
	// Mã hóa struct User thành JSON và gửi về client
	json.NewEncoder(w).Encode(user)
}

func main() {
	// Định nghĩa endpoint /users/{id}
	// http.HandleFunc đăng ký một hàm handler cho một path cụ thể
	http.HandleFunc("/users/", getUserHandler)

	// Khởi động server HTTP trên cổng 8080
	port := ":8080"
	fmt.Printf("Microservice người dùng đang lắng nghe trên cổng %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil)) // log.Fatal sẽ in lỗi và thoát nếu server không thể khởi động
}
