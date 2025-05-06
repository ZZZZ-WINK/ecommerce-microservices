module user-service

go 1.21

require (
	golang.org/x/crypto v0.19.0
	google.golang.org/grpc v1.62.0
	google.golang.org/protobuf v1.32.0
	gorm.io/driver/mysql v1.5.4
	gorm.io/gorm v1.25.7
)

replace user-service/common/proto => ../common/proto

require (
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240123012728-ef4313101c80 // indirect
)
