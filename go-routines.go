package main

import (
	"fmt"
	"time"
)

func guiTinNhan(kenh chan string) {
	time.Sleep(1 * time.Second)        // Giả lập công việc
	kenh <- "Xin chàoti từ Goroutine!" // Gửi tin nhắn vào channel
}

func main() {
	// Tạo một channel có kiểu dữ liệu là string
	thongDiep := make(chan string)

	// Khởi tạo một Goroutine để gửi tin nhắn
	go guiTinNhan(thongDiep)

	// Nhận tin nhắn từ channel (chương trình chính sẽ chờ ở đây cho đến khi có tin nhắn)
	tinNhanNhanDuoc := <-thongDiep
	fmt.Println(tinNhanNhanDuoc)

	fmt.Println("Chương trình chính kết thúc.")
}
