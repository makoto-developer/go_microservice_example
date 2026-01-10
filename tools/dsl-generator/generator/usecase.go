package generator

import (
	"fmt"
	"strings"

	"github.com/makoto-developer/go_microservice_example/tools/dsl-generator/parser"
)

type UsecaseGenerator struct {
	service *parser.Microservice
}

func NewUsecaseGenerator(service *parser.Microservice) *UsecaseGenerator {
	return &UsecaseGenerator{service: service}
}

func (g *UsecaseGenerator) GenerateUsecases() map[string]string {
	files := make(map[string]string)

	for _, usecase := range g.service.Usecases {
		fileName := toSnakeCase(usecase.Name) + ".go"
		content := g.generateUsecaseFile(usecase)
		files[fileName] = content
	}

	return files
}

func (g *UsecaseGenerator) generateUsecaseFile(usecase parser.Usecase) string {
	var sb strings.Builder

	sb.WriteString("package usecase\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("\t\"context\"\n")
	sb.WriteString("\t\"fmt\"\n")
	sb.WriteString("\n")
	sb.WriteString("\t\"github.com/google/uuid\"\n")
	sb.WriteString(")\n\n")

	// Input struct
	if len(usecase.Input) > 0 {
		sb.WriteString(g.generateInputStruct(usecase))
		sb.WriteString("\n")
	}

	// Output struct
	if len(usecase.Output) > 0 {
		sb.WriteString(g.generateOutputStruct(usecase))
		sb.WriteString("\n")
	}

	// Error types
	for _, errType := range usecase.Errors {
		sb.WriteString(g.generateErrorType(errType))
		sb.WriteString("\n")
	}

	// Interface
	sb.WriteString(g.generateUsecaseInterface(usecase))
	sb.WriteString("\n")

	// Implementation struct
	sb.WriteString(g.generateUsecaseImpl(usecase))
	sb.WriteString("\n")

	// Constructor
	sb.WriteString(g.generateUsecaseConstructor(usecase))
	sb.WriteString("\n")

	// Execute method
	sb.WriteString(g.generateExecuteMethod(usecase))

	return sb.String()
}

func (g *UsecaseGenerator) generateInputStruct(usecase parser.Usecase) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("// %sInput represents input for %s\n", usecase.Name, usecase.Name))
	sb.WriteString(fmt.Sprintf("type %sInput struct {\n", usecase.Name))

	for _, field := range usecase.Input {
		goType := g.mapDSLTypeToGo(field.Type)
		sb.WriteString(fmt.Sprintf("\t%s %s\n", toPascalCase(field.Name), goType))
	}

	sb.WriteString("}\n")

	return sb.String()
}

func (g *UsecaseGenerator) generateOutputStruct(usecase parser.Usecase) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("// %sOutput represents output for %s\n", usecase.Name, usecase.Name))
	sb.WriteString(fmt.Sprintf("type %sOutput struct {\n", usecase.Name))

	for _, field := range usecase.Output {
		goType := g.mapDSLTypeToGo(field.Type)
		sb.WriteString(fmt.Sprintf("\t%s %s\n", toPascalCase(field.Name), goType))
	}

	sb.WriteString("}\n")

	return sb.String()
}

func (g *UsecaseGenerator) generateErrorType(errType string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("// %sError represents %s error\n", errType, errType))
	sb.WriteString(fmt.Sprintf("type %sError struct {\n", errType))
	sb.WriteString("\tMessage string\n")
	sb.WriteString("}\n\n")

	sb.WriteString(fmt.Sprintf("func (e %sError) Error() string {\n", errType))
	sb.WriteString("\tif e.Message != \"\" {\n")
	sb.WriteString("\t\treturn e.Message\n")
	sb.WriteString("\t}\n")
	sb.WriteString(fmt.Sprintf("\treturn \"%s\"\n", toSnakeCase(errType)))
	sb.WriteString("}\n")

	return sb.String()
}

func (g *UsecaseGenerator) generateUsecaseInterface(usecase parser.Usecase) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("// %sUsecase defines the interface for %s\n", usecase.Name, usecase.Name))
	sb.WriteString(fmt.Sprintf("type %sUsecase interface {\n", usecase.Name))

	inputType := "input"
	if len(usecase.Input) > 0 {
		inputType = fmt.Sprintf("input %sInput", usecase.Name)
	}

	outputType := "error"
	if len(usecase.Output) > 0 {
		outputType = fmt.Sprintf("(%sOutput, error)", usecase.Name)
	}

	sb.WriteString(fmt.Sprintf("\tExecute(ctx context.Context, %s) %s\n", inputType, outputType))
	sb.WriteString("}\n")

	return sb.String()
}

func (g *UsecaseGenerator) generateUsecaseImpl(usecase parser.Usecase) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("type %sUsecaseImpl struct {\n", toSnakeCase(usecase.Name)))
	sb.WriteString("\t// TODO: Add repository dependencies\n")
	sb.WriteString("}\n")

	return sb.String()
}

func (g *UsecaseGenerator) generateUsecaseConstructor(usecase parser.Usecase) string {
	var sb strings.Builder

	implName := toSnakeCase(usecase.Name) + "UsecaseImpl"
	sb.WriteString(fmt.Sprintf("// New%sUsecase creates a new instance\n", usecase.Name))
	sb.WriteString(fmt.Sprintf("func New%sUsecase() %sUsecase {\n", usecase.Name, usecase.Name))
	sb.WriteString(fmt.Sprintf("\treturn &%s{}\n", implName))
	sb.WriteString("}\n")

	return sb.String()
}

func (g *UsecaseGenerator) generateExecuteMethod(usecase parser.Usecase) string {
	var sb strings.Builder

	implName := toSnakeCase(usecase.Name) + "UsecaseImpl"

	inputParam := "_"
	if len(usecase.Input) > 0 {
		inputParam = fmt.Sprintf("input %sInput", usecase.Name)
	}

	returnType := "error"
	if len(usecase.Output) > 0 {
		returnType = fmt.Sprintf("(%sOutput, error)", usecase.Name)
	}

	sb.WriteString(fmt.Sprintf("// Execute executes %s\n", usecase.Name))
	sb.WriteString(fmt.Sprintf("func (u *%s) Execute(ctx context.Context, %s) %s {\n", implName, inputParam, returnType))

	// Generate basic implementation structure
	sb.WriteString("\t// TODO: Implement business logic\n")

	if usecase.Transaction {
		sb.WriteString("\t// TODO: Begin transaction\n")
	}

	if usecase.Saga != nil {
		sb.WriteString("\t// Saga pattern implementation\n")
		for i, step := range usecase.Saga.Steps {
			sb.WriteString(fmt.Sprintf("\t// Step %d: %s\n", i+1, step.Name))
		}
	}

	if len(usecase.Output) > 0 {
		sb.WriteString(fmt.Sprintf("\n\treturn %sOutput{}, nil\n", usecase.Name))
	} else {
		sb.WriteString("\n\treturn nil\n")
	}

	sb.WriteString("}\n")

	return sb.String()
}

func (g *UsecaseGenerator) mapDSLTypeToGo(dslType string) string {
	typeMap := map[string]string{
		"UUID":      "uuid.UUID",
		"string":    "string",
		"int":       "int",
		"boolean":   "bool",
		"timestamp": "time.Time",
		"decimal":   "decimal.Decimal",
	}

	if strings.HasPrefix(dslType, "list<") {
		innerType := strings.TrimSuffix(strings.TrimPrefix(dslType, "list<"), ">")
		return "[]" + g.mapDSLTypeToGo(innerType)
	}

	if mapped, ok := typeMap[dslType]; ok {
		return mapped
	}

	return dslType
}
