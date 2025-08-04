// go-tutorials/json.go
package main

import (
	"encoding/json" // Import package encoding/json
	"fmt"
)

// Định nghĩa một struct để đại diện cho dữ liệu
type ServerConfig struct {
	Name    string `json:"server_name"` // Sử dụng tag 'json' để tùy chỉnh tên trường trong JSON
	IP      string `json:"ip_address"`
	Ports   []int  `json:"ports"`
	Enabled bool   `json:"is_enabled"`
}

func main() {
	// --- Mã hóa (Encoding) Go struct thành JSON ---
	config := ServerConfig{
		Name:    "Web_Server_01",
		IP:      "192.168.1.100",
		Ports:   []int{80, 443, 8080},
		Enabled: true,
	}

	// Chuyển đổi struct thành JSON byte slice
	jsonData, err := json.MarshalIndent(config, "", "  ") // MarshalIndent giúp format JSON dễ đọc hơn
	if err != nil {
		fmt.Println("Lỗi khi mã hóa JSON:", err)
		return
	}
	fmt.Println("Dữ liệu JSON đã mã hóa:\n", string(jsonData))

	// --- Giải mã (Decoding) JSON thành Go struct ---
	jsonString := `
	{
		"server_name": "DB_Server_02",
		"ip_address": "10.0.0.5",
		"ports": [3306, 5432],
		"is_enabled": false
	}`

	var decodedConfig ServerConfig // Khai báo một biến struct rỗng để chứa dữ liệu giải mã
	err = json.Unmarshal([]byte(jsonString), &decodedConfig)
	if err != nil {
		fmt.Println("Lỗi khi giải mã JSON:", err)
		return
	}
	fmt.Println("\nDữ liệu JSON đã giải mã thành struct Go:")
	fmt.Printf("Tên Server: %s\n", decodedConfig.Name)
	fmt.Printf("Địa chỉ IP: %s\n", decodedConfig.IP)
	fmt.Printf("Cổng: %v\n", decodedConfig.Ports)
	fmt.Printf("Kích hoạt: %t\n", decodedConfig.Enabled)
}
