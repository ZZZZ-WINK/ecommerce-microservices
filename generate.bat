@echo off

REM 生成用户服务
D:\download\protoc\bin\protoc --go_out=user-service --go_opt=paths=source_relative ^
    --go-grpc_out=user-service --go-grpc_opt=paths=source_relative ^
    common/proto/user.proto

REM 生成商品服务
D:\download\protoc\bin\protoc --go_out=product-service --go_opt=paths=source_relative ^
    --go-grpc_out=product-service --go-grpc_opt=paths=source_relative ^
    common/proto/product.proto

REM 生成订单服务
D:\download\protoc\bin\protoc --go_out=order-service --go_opt=paths=source_relative ^
    --go-grpc_out=order-service --go-grpc_opt=paths=source_relative ^
    common/proto/order.proto 