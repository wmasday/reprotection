package main

import (
	"fmt"
	"reprotection/config"
	"reprotection/models"
)

func main() {
	config.ConnectDB()

	err := models.CreateUser("admin", "admin123")
	if err != nil {
		fmt.Println("❌ Gagal buat user:", err)
	} else {
		fmt.Println("✅ User admin berhasil dibuat")
	}
}
