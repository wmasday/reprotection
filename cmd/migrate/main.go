package main

import (
	"reprotection/config"
)

func main() {
	config.ConnectDB()
	config.Migrate()
	// Remove Go-based config seeder, seeding is now handled in SQL migration
}
