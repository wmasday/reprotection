package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"reprotection/config"
	"reprotection/controllers"
	"reprotection/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Gagal load file .env")
	}


	port := os.Getenv("APP_PORT")
	if port == "" {
		log.Fatal("❌ APP_PORT belum diatur di file .env")
	}

	config.ConnectDB()
	log.Println("✅ Terhubung ke DB:", os.Getenv("DB_NAME"))

	// Static Resource
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Auth routes
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/auth", controllers.Auth)
	http.HandleFunc("/logout", controllers.Logout)

	// Change password routes
	http.HandleFunc("/change-password", controllers.ShowChangePassword)
	http.HandleFunc("/change-password/submit", controllers.ChangePassword)

	// Protected routes
	http.HandleFunc("/", middleware.Auth(controllers.Index))
	http.HandleFunc("/store", middleware.Auth(controllers.Store))
	http.HandleFunc("/delete", middleware.Auth(controllers.Delete))
	http.HandleFunc("/config/store", middleware.Auth(controllers.StoreConfig))
	http.HandleFunc("/detail", middleware.Auth(controllers.Detail))
	http.HandleFunc("/sync", middleware.Auth(controllers.Sync))

	// Running server
	log.Println("🚀 Server running at http://localhost:" + port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("❌ ListenAndServe error:", err)
	}
}
