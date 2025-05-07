package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// 获取当前工作目录
	workDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("获取工作目录失败: %v\n", err)
		os.Exit(1)
	}

	// 目录设置
	protoDir := filepath.Join(workDir, "common", "proto", "protos")
	genDir := filepath.Join(workDir, "common", "proto", "gen")

	// 创建生成目录
	services := []string{"user", "product", "order"}
	for _, service := range services {
		dir := filepath.Join(genDir, service)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("创建目录失败: %v\n", err)
			os.Exit(1)
		}
	}

	// 获取所有proto文件
	protoFiles, err := filepath.Glob(filepath.Join(protoDir, "*.proto"))
	if err != nil {
		fmt.Printf("查找proto文件失败: %v\n", err)
		os.Exit(1)
	}

	// 为每个proto文件生成代码
	for _, protoFile := range protoFiles {
		fileName := filepath.Base(protoFile)
		serviceName := fileName[:len(fileName)-6] // 移除.proto后缀
		fmt.Printf("正在生成 %s 的代码...\n", fileName)

		outputPath := filepath.Join("common", "proto", "gen", serviceName)

		// 生成Go代码
		cmd := exec.Command("protoc",
			"--proto_path="+filepath.Join("common", "proto", "protos"),
			"--go_out="+outputPath,
			"--go_opt=paths=source_relative",
			"--go-grpc_out="+outputPath,
			"--go-grpc_opt=paths=source_relative",
			fileName)

		cmd.Dir = workDir
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("生成 %s 失败: %v\n输出: %s\n", fileName, err, string(output))
			os.Exit(1)
		}
	}

	fmt.Println("所有proto文件生成完成！")
}
