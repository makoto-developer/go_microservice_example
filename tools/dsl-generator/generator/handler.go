package generator

import (
	"fmt"
	"strings"

	"github.com/makoto-developer/go_microservice_example/tools/dsl-generator/parser"
)

type HandlerGenerator struct {
	service *parser.Microservice
}

func NewHandlerGenerator(service *parser.Microservice) *HandlerGenerator {
	return &HandlerGenerator{service: service}
}

func (g *HandlerGenerator) GenerateHandler() string {
	var sb strings.Builder

	sb.WriteString("package handler\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("\t\"context\"\n")
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("\tpb \"github.com/makoto-developer/go_microservice_example/proto/%s/%s\"\n",
		toSnakeCase(g.service.Name), g.service.Version))
	sb.WriteString(fmt.Sprintf("\t\"github.com/makoto-developer/go_microservice_example/generated/%s/usecase\"\n",
		toSnakeCase(g.service.Name)))
	sb.WriteString(")\n\n")

	// Handler struct
	sb.WriteString(g.generateHandlerStruct())
	sb.WriteString("\n")

	// Constructor
	sb.WriteString(g.generateHandlerConstructor())
	sb.WriteString("\n")

	// Generate handler methods for each gRPC method
	if g.service.GRPCService != nil {
		for _, method := range g.service.GRPCService.Methods {
			sb.WriteString(g.generateHandlerMethod(method))
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func (g *HandlerGenerator) generateHandlerStruct() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("// %sHandler implements gRPC handler\n", g.service.Name))
	sb.WriteString(fmt.Sprintf("type %sHandler struct {\n", g.service.Name))
	sb.WriteString(fmt.Sprintf("\tpb.Unimplemented%sServer\n", g.service.Name))

	// Add usecase dependencies
	for _, uc := range g.service.Usecases {
		fieldName := toSnakeCase(uc.Name) + "Usecase"
		sb.WriteString(fmt.Sprintf("\t%s usecase.%sUsecase\n", fieldName, uc.Name))
	}

	sb.WriteString("}\n")

	return sb.String()
}

func (g *HandlerGenerator) generateHandlerConstructor() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("// New%sHandler creates a new handler instance\n", g.service.Name))
	sb.WriteString(fmt.Sprintf("func New%sHandler(\n", g.service.Name))

	// Parameters
	for _, uc := range g.service.Usecases {
		fieldName := toSnakeCase(uc.Name) + "Usecase"
		sb.WriteString(fmt.Sprintf("\t%s usecase.%sUsecase,\n", fieldName, uc.Name))
	}

	sb.WriteString(fmt.Sprintf(") *%sHandler {\n", g.service.Name))
	sb.WriteString(fmt.Sprintf("\treturn &%sHandler{\n", g.service.Name))

	// Initialize fields
	for _, uc := range g.service.Usecases {
		fieldName := toSnakeCase(uc.Name) + "Usecase"
		sb.WriteString(fmt.Sprintf("\t\t%s: %s,\n", fieldName, fieldName))
	}

	sb.WriteString("\t}\n}\n")

	return sb.String()
}

func (g *HandlerGenerator) generateHandlerMethod(method parser.GRPCMethod) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("// %s handles %s RPC\n", method.Name, method.Name))
	sb.WriteString(fmt.Sprintf("func (h *%sHandler) %s(\n", g.service.Name, method.Name))
	sb.WriteString("\tctx context.Context,\n")
	sb.WriteString(fmt.Sprintf("\treq *pb.%s,\n", method.Request))
	sb.WriteString(fmt.Sprintf(") (*pb.%s, error) {\n", method.Response))

	sb.WriteString("\t// TODO: Implement handler logic\n")
	sb.WriteString("\t// 1. Convert request to usecase input\n")
	sb.WriteString("\t// 2. Execute usecase\n")
	sb.WriteString("\t// 3. Convert usecase output to response\n")

	sb.WriteString(fmt.Sprintf("\n\treturn &pb.%s{}, nil\n", method.Response))
	sb.WriteString("}\n")

	return sb.String()
}
