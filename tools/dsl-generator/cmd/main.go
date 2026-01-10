package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/makoto-developer/go_microservice_example/tools/dsl-generator/generator"
	"github.com/makoto-developer/go_microservice_example/tools/dsl-generator/parser"
)

func main() {
	var (
		inputFile  string
		outputDir  string
		serviceName string
	)

	flag.StringVar(&inputFile, "input", "", "Path to DSL definition file (.model)")
	flag.StringVar(&outputDir, "output", "", "Output directory for generated code")
	flag.StringVar(&serviceName, "service", "", "Service name (e.g., auth-service)")
	flag.Parse()

	if inputFile == "" || outputDir == "" || serviceName == "" {
		fmt.Println("Usage: dsl-generator -input <dsl-file> -output <output-dir> -service <service-name>")
		os.Exit(1)
	}

	// Parse DSL file
	fmt.Printf("Parsing DSL file: %s\n", inputFile)
	p, err := parser.NewParser(inputFile)
	if err != nil {
		fmt.Printf("Error creating parser: %v\n", err)
		os.Exit(1)
	}

	service, err := p.Parse()
	if err != nil {
		fmt.Printf("Error parsing DSL: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully parsed service: %s (version %s)\n", service.Name, service.Version)

	// Create output directory structure
	domainDir := filepath.Join(outputDir, serviceName, "domain")
	usecaseDir := filepath.Join(outputDir, serviceName, "usecase")
	handlerDir := filepath.Join(outputDir, serviceName, "handler")
	protoDir := filepath.Join(outputDir, "..", "proto", serviceName, service.Version)

	for _, dir := range []string{domainDir, usecaseDir, handlerDir, protoDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dir, err)
			os.Exit(1)
		}
	}

	// Generate domain layer
	fmt.Println("\nGenerating domain layer...")
	domainGen := generator.NewDomainGenerator(service)

	// Generate entities
	entities := domainGen.GenerateEntities()
	for fileName, content := range entities {
		filePath := filepath.Join(domainDir, fileName)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			fmt.Printf("Error writing file %s: %v\n", filePath, err)
			os.Exit(1)
		}
		fmt.Printf("  ✓ Generated %s\n", fileName)
	}

	// Generate enums
	enums := domainGen.GenerateEnums()
	for fileName, content := range enums {
		filePath := filepath.Join(domainDir, fileName)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			fmt.Printf("Error writing file %s: %v\n", filePath, err)
			os.Exit(1)
		}
		fmt.Printf("  ✓ Generated %s\n", fileName)
	}

	// Generate value objects
	valueObjects := domainGen.GenerateValueObjects()
	for fileName, content := range valueObjects {
		filePath := filepath.Join(domainDir, fileName)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			fmt.Printf("Error writing file %s: %v\n", filePath, err)
			os.Exit(1)
		}
		fmt.Printf("  ✓ Generated %s\n", fileName)
	}

	// Generate repositories
	repoGen := generator.NewRepositoryGenerator(service)
	repositories := repoGen.GenerateRepositories()
	for fileName, content := range repositories {
		filePath := filepath.Join(domainDir, fileName)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			fmt.Printf("Error writing file %s: %v\n", filePath, err)
			os.Exit(1)
		}
		fmt.Printf("  ✓ Generated %s\n", fileName)
	}

	// Generate usecase layer
	fmt.Println("\nGenerating usecase layer...")
	usecaseGen := generator.NewUsecaseGenerator(service)
	usecases := usecaseGen.GenerateUsecases()
	for fileName, content := range usecases {
		filePath := filepath.Join(usecaseDir, fileName)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			fmt.Printf("Error writing file %s: %v\n", filePath, err)
			os.Exit(1)
		}
		fmt.Printf("  ✓ Generated %s\n", fileName)
	}

	// Generate handler layer
	fmt.Println("\nGenerating handler layer...")
	handlerGen := generator.NewHandlerGenerator(service)
	handlerContent := handlerGen.GenerateHandler()
	handlerFile := filepath.Join(handlerDir, "grpc_handler.go")
	if err := os.WriteFile(handlerFile, []byte(handlerContent), 0644); err != nil {
		fmt.Printf("Error writing handler file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  ✓ Generated grpc_handler.go\n")

	// Generate Protocol Buffers
	fmt.Println("\nGenerating Protocol Buffers...")
	protoGen := generator.NewProtoGenerator(service)
	protoContent := protoGen.GenerateProto()
	protoFile := filepath.Join(protoDir, serviceName+".proto")
	if err := os.WriteFile(protoFile, []byte(protoContent), 0644); err != nil {
		fmt.Printf("Error writing proto file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  ✓ Generated %s.proto\n", serviceName)

	// Generate go.mod
	fmt.Println("\nGenerating go.mod...")
	goModContent := fmt.Sprintf(`module github.com/makoto-developer/go_microservice_example/generated/%s

go 1.25

require (
	github.com/google/uuid v1.6.0
	github.com/shopspring/decimal v1.4.0
	google.golang.org/grpc v1.68.1
	google.golang.org/protobuf v1.36.1
)
`, serviceName)
	goModFile := filepath.Join(outputDir, serviceName, "go.mod")
	if err := os.WriteFile(goModFile, []byte(goModContent), 0644); err != nil {
		fmt.Printf("Error writing go.mod: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  ✓ Generated go.mod\n")

	fmt.Printf("\n✨ Successfully generated code for %s\n", serviceName)
	fmt.Printf("   Output directory: %s\n", filepath.Join(outputDir, serviceName))
	fmt.Printf("   Proto file: %s\n", protoFile)
}
