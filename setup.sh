go mod init reprotection
go get github.com/go-sql-driver/mysql
go get golang.org/x/crypto/bcrypt
go get github.com/gorilla/sessions
go get github.com/joho/godotenv
go run cmd/migrate/main.go
go run cmd/create_user/main.go