package main

import (
	"reprotection/config"
)

func main() {
	config.ConnectDB()
	config.Migrate()
}
