package main

import (
	"fmt"
)

func main() {
	j := 0
	// In ra các số chẵn từ 1 đến 100
	fmt.Println("Các số chẵn từ 1 đến 100:")
	// Sử dụng vòng lặp for để lặp qua các số từ 1 đến 100
	for i := 1; i <= 100; i++ {
		if i%2 == 0 {
			fmt.Println(i)
			j = j + i
		}
	}
	fmt.Println("Tổng các số chẵn từ 1 đến 100 là:", j)
	fmt.Println("Kết thúc chương trình.")
}
