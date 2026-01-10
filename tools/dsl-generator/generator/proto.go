package generator

import (
	"fmt"
	"strings"

	"github.com/makoto-developer/go_microservice_example/tools/dsl-generator/parser"
)

type ProtoGenerator struct {
	service *parser.Microservice
}

func NewProtoGenerator(service *parser.Microservice) *ProtoGenerator {
	return &ProtoGenerator{service: service}
}

func (g *ProtoGenerator) GenerateProto() string {
	var sb strings.Builder

	// Header
	sb.WriteString("syntax = \"proto3\";\n\n")
	sb.WriteString(fmt.Sprintf("package %s.%s;\n\n", toSnakeCase(g.service.Name), g.service.Version))
	sb.WriteString(fmt.Sprintf("option go_package = \"github.com/makoto-developer/go_microservice_example/proto/%s/%s\";\n\n",
		toSnakeCase(g.service.Name), g.service.Version))

	// Import common types
	sb.WriteString("import \"google/protobuf/timestamp.proto\";\n")
	sb.WriteString("import \"google/protobuf/empty.proto\";\n\n")

	// Generate enums
	for _, enum := range g.service.Enums {
		sb.WriteString(g.generateProtoEnum(enum))
		sb.WriteString("\n")
	}

	// Generate messages for entities
	for _, entity := range g.service.Entities {
		sb.WriteString(g.generateProtoMessage(entity.Name, entity.Fields))
		sb.WriteString("\n")
	}

	// Generate request/response messages for usecases
	for _, usecase := range g.service.Usecases {
		if len(usecase.Input) > 0 {
			sb.WriteString(g.generateProtoMessage(usecase.Name+"Request", usecase.Input))
			sb.WriteString("\n")
		}
		if len(usecase.Output) > 0 {
			sb.WriteString(g.generateProtoMessage(usecase.Name+"Response", usecase.Output))
			sb.WriteString("\n")
		}
	}

	// Generate service definition
	if g.service.GRPCService != nil {
		sb.WriteString(g.generateProtoService())
	}

	return sb.String()
}

func (g *ProtoGenerator) generateProtoEnum(enum parser.Enum) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("enum %s {\n", enum.Name))
	sb.WriteString(fmt.Sprintf("\t%s_UNSPECIFIED = 0;\n", strings.ToUpper(toSnakeCase(enum.Name))))

	for i, value := range enum.Values {
		sb.WriteString(fmt.Sprintf("\t%s = %d;\n", value, i+1))
	}

	sb.WriteString("}\n")

	return sb.String()
}

func (g *ProtoGenerator) generateProtoMessage(name string, fields []parser.Field) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("message %s {\n", name))

	for i, field := range fields {
		protoType := g.mapGoTypeToProto(field.Type)
		fieldNum := i + 1
		sb.WriteString(fmt.Sprintf("\t%s %s = %d;\n", protoType, toSnakeCase(field.Name), fieldNum))
	}

	sb.WriteString("}\n")

	return sb.String()
}

func (g *ProtoGenerator) generateProtoService() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("service %s {\n", g.service.Name))

	for _, method := range g.service.GRPCService.Methods {
		sb.WriteString(fmt.Sprintf("\trpc %s(%s) returns (%s);\n",
			method.Name,
			method.Request,
			method.Response))
	}

	sb.WriteString("}\n")

	return sb.String()
}

func (g *ProtoGenerator) mapGoTypeToProto(goType string) string {
	typeMap := map[string]string{
		"UUID":      "string",
		"string":    "string",
		"int":       "int32",
		"boolean":   "bool",
		"timestamp": "google.protobuf.Timestamp",
		"decimal":   "string", // Decimal as string in proto
	}

	// Handle list types
	if strings.HasPrefix(goType, "list<") {
		innerType := strings.TrimSuffix(strings.TrimPrefix(goType, "list<"), ">")
		return "repeated " + g.mapGoTypeToProto(innerType)
	}

	if mapped, ok := typeMap[goType]; ok {
		return mapped
	}

	// Custom type (enum or message)
	return goType
}
