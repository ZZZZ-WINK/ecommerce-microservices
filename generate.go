package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// Define command line arguments
	serviceName := flag.String("service", "", "Specify the service to generate (user/product/order)")
	clean := flag.Bool("clean", false, "Clean old generated files")
	flag.Parse()

	// Get current working directory
	workDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Failed to get working directory: %v\n", err)
		os.Exit(1)
	}

	// Directory settings
	protoDir := filepath.Join(workDir, "common", "proto", "protos")
	genDir := filepath.Join(workDir, "common", "proto", "gen")

	// Clean generated files if needed
	if *clean {
		if *serviceName != "" {
			// Only clean the specified service's generated files
			serviceDir := filepath.Join(genDir, *serviceName)
			if err := os.RemoveAll(serviceDir); err != nil {
				fmt.Printf("Failed to clean directory %s: %v\n", serviceDir, err)
				os.Exit(1)
			}
			fmt.Printf("Cleaned generated files for service: %s\n", *serviceName)
		} else {
			// Clean all generated files
			if err := os.RemoveAll(genDir); err != nil {
				fmt.Printf("Failed to clean directory: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Cleaned all generated files.")
		}
	}

	// Determine which services to process
	var services []string
	if *serviceName != "" {
		services = []string{*serviceName}
	} else {
		services = []string{"user", "product", "order"}
	}

	// Create generation directories
	for _, service := range services {
		dir := filepath.Join(genDir, service)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Failed to create directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Get all proto files
	protoFiles, err := filepath.Glob(filepath.Join(protoDir, "*.proto"))
	if err != nil {
		fmt.Printf("Failed to find proto files: %v\n", err)
		os.Exit(1)
	}

	// Generate code for each proto file
	for _, protoFile := range protoFiles {
		fileName := filepath.Base(protoFile)
		fileServiceName := fileName[:len(fileName)-6] // Remove .proto suffix

		// If a specific service is specified, only process its proto file
		if *serviceName != "" && fileServiceName != *serviceName {
			continue
		}

		fmt.Printf("Generating code for %s ...\n", fileName)

		outputPath := filepath.Join("common", "proto", "gen", fileServiceName)

		// Generate Go code
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
			fmt.Printf("Failed to generate %s: %v\nOutput: %s\n", fileName, err, string(output))
			os.Exit(1)
		}
	}

	fmt.Println("Code generation completed!")
}
