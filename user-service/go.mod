module user-service

go 1.23

toolchain go1.24.2

require (
	common v0.0.0
	github.com/joho/godotenv v1.5.1
	golang.org/x/crypto v0.33.0
	google.golang.org/grpc v1.72.0
	gorm.io/driver/mysql v1.5.4
	gorm.io/gorm v1.25.7
)

replace common => ../common

require (
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
	google.golang.org/protobuf v1.36.5 // indirect
)
