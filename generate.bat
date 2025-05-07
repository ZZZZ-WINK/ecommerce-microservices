@echo off
echo 正在生成proto文件...

:: 检查protoc是否安装
where protoc >nul 2>nul
if %errorlevel% neq 0 (
    echo 错误: 未找到protoc，请先安装protoc
    exit /b 1
)

:: 检查go是否安装
where go >nul 2>nul
if %errorlevel% neq 0 (
    echo 错误: 未找到go，请先安装go
    exit /b 1
)

:: 安装必要的go插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

:: 运行生成脚本
go run generate.go

echo 生成完成！ 