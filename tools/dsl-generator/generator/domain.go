package generator

import (
	"fmt"
	"strings"

	"github.com/makoto-developer/go_microservice_example/tools/dsl-generator/parser"
)

type DomainGenerator struct {
	service *parser.Microservice
}

func NewDomainGenerator(service *parser.Microservice) *DomainGenerator {
	return &DomainGenerator{service: service}
}

func (g *DomainGenerator) GenerateEntities() map[string]string {
	files := make(map[string]string)

	for _, entity := range g.service.Entities {
		fileName := toSnakeCase(entity.Name) + ".go"
		content := g.generateEntityFile(entity)
		files[fileName] = content
	}

	return files
}

func (g *DomainGenerator) generateEntityFile(entity parser.Entity) string {
	var sb strings.Builder

	// Package and imports
	sb.WriteString(fmt.Sprintf("package domain\n\n"))
	sb.WriteString(g.generateImports(entity.Fields))
	sb.WriteString("\n")

	// Struct definition
	sb.WriteString(fmt.Sprintf("// %s represents %s\n", entity.Name, entity.Name))
	sb.WriteString(fmt.Sprintf("type %s struct {\n", entity.Name))

	for _, field := range entity.Fields {
		sb.WriteString(g.generateStructField(field))
	}

	sb.WriteString("}\n\n")

	// Constructor
	sb.WriteString(g.generateConstructor(entity))

	return sb.String()
}

func (g *DomainGenerator) generateStructField(field parser.Field) string {
	goType := g.mapDSLTypeToGo(field.Type)
	jsonTag := toSnakeCase(field.Name)
	dbTag := toSnakeCase(field.Name)

	// Handle nullable fields
	if contains(field.Constraints, "nullable") {
		goType = "*" + goType
	}

	// JSON omitempty for nullable fields
	jsonOpts := ""
	if contains(field.Constraints, "nullable") {
		jsonOpts = ",omitempty"
	}

	// Hide password fields from JSON
	if strings.Contains(field.Name, "password") || strings.Contains(field.Name, "token") {
		jsonTag = "-"
	} else {
		jsonTag = jsonTag + jsonOpts
	}

	return fmt.Sprintf("\t%s %s `db:\"%s\" json:\"%s\"`\n",
		toPascalCase(field.Name),
		goType,
		dbTag,
		jsonTag,
	)
}

func (g *DomainGenerator) generateImports(fields []parser.Field) string {
	imports := make(map[string]bool)

	for _, field := range fields {
		switch field.Type {
		case "UUID":
			imports["github.com/google/uuid"] = true
		case "timestamp":
			imports["time"] = true
		case "decimal":
			imports["github.com/shopspring/decimal"] = true
		}
	}

	if len(imports) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("import (\n")
	for imp := range imports {
		sb.WriteString(fmt.Sprintf("\t\"%s\"\n", imp))
	}
	sb.WriteString(")\n")

	return sb.String()
}

func (g *DomainGenerator) generateConstructor(entity parser.Entity) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("// New%s creates a new %s instance\n", entity.Name, entity.Name))
	sb.WriteString(fmt.Sprintf("func New%s() *%s {\n", entity.Name, entity.Name))
	sb.WriteString(fmt.Sprintf("\treturn &%s{}\n", entity.Name))
	sb.WriteString("}\n")

	return sb.String()
}

func (g *DomainGenerator) GenerateEnums() map[string]string {
	files := make(map[string]string)

	for _, enum := range g.service.Enums {
		fileName := toSnakeCase(enum.Name) + ".go"
		content := g.generateEnumFile(enum)
		files[fileName] = content
	}

	return files
}

func (g *DomainGenerator) generateEnumFile(enum parser.Enum) string {
	var sb strings.Builder

	sb.WriteString("package domain\n\n")

	// Type definition
	sb.WriteString(fmt.Sprintf("// %s represents %s type\n", enum.Name, enum.Name))
	sb.WriteString(fmt.Sprintf("type %s string\n\n", enum.Name))

	// Constants
	sb.WriteString(fmt.Sprintf("const (\n"))
	for _, value := range enum.Values {
		constName := enum.Name + toPascalCase(strings.ToLower(value))
		sb.WriteString(fmt.Sprintf("\t%s %s = \"%s\"\n", constName, enum.Name, value))
	}
	sb.WriteString(")\n\n")

	// Values function
	sb.WriteString(fmt.Sprintf("// %sValues returns all possible values\n", enum.Name))
	sb.WriteString(fmt.Sprintf("func %sValues() []%s {\n", enum.Name, enum.Name))
	sb.WriteString(fmt.Sprintf("\treturn []%s{\n", enum.Name))
	for _, value := range enum.Values {
		constName := enum.Name + toPascalCase(strings.ToLower(value))
		sb.WriteString(fmt.Sprintf("\t\t%s,\n", constName))
	}
	sb.WriteString("\t}\n}\n\n")

	// IsValid function
	sb.WriteString(fmt.Sprintf("// IsValid checks if the value is valid\n"))
	sb.WriteString(fmt.Sprintf("func (e %s) IsValid() bool {\n", enum.Name))
	sb.WriteString("\tswitch e {\n")
	for _, value := range enum.Values {
		constName := enum.Name + toPascalCase(strings.ToLower(value))
		sb.WriteString(fmt.Sprintf("\tcase %s:\n", constName))
	}
	sb.WriteString("\t\treturn true\n")
	sb.WriteString("\t}\n\treturn false\n}\n")

	return sb.String()
}

func (g *DomainGenerator) GenerateValueObjects() map[string]string {
	files := make(map[string]string)

	for _, vo := range g.service.ValueObjects {
		fileName := toSnakeCase(vo.Name) + ".go"
		content := g.generateValueObjectFile(vo)
		files[fileName] = content
	}

	return files
}

func (g *DomainGenerator) generateValueObjectFile(vo parser.ValueObject) string {
	var sb strings.Builder

	sb.WriteString("package domain\n\n")
	sb.WriteString(g.generateImports(vo.Fields))
	sb.WriteString("\n")

	sb.WriteString(fmt.Sprintf("// %s represents %s value object\n", vo.Name, vo.Name))
	sb.WriteString(fmt.Sprintf("type %s struct {\n", vo.Name))

	for _, field := range vo.Fields {
		sb.WriteString(g.generateStructField(field))
	}

	sb.WriteString("}\n\n")

	// Constructor
	sb.WriteString(fmt.Sprintf("// New%s creates a new %s instance\n", vo.Name, vo.Name))
	sb.WriteString(fmt.Sprintf("func New%s() *%s {\n", vo.Name, vo.Name))
	sb.WriteString(fmt.Sprintf("\treturn &%s{}\n", vo.Name))
	sb.WriteString("}\n")

	return sb.String()
}

func (g *DomainGenerator) mapDSLTypeToGo(dslType string) string {
	typeMap := map[string]string{
		"UUID":      "uuid.UUID",
		"string":    "string",
		"int":       "int",
		"boolean":   "bool",
		"timestamp": "time.Time",
		"decimal":   "decimal.Decimal",
		"json":      "map[string]interface{}",
	}

	// Handle list types
	if strings.HasPrefix(dslType, "list<") {
		innerType := strings.TrimSuffix(strings.TrimPrefix(dslType, "list<"), ">")
		return "[]" + g.mapDSLTypeToGo(innerType)
	}

	if mapped, ok := typeMap[dslType]; ok {
		return mapped
	}

	// Assume it's a custom type (enum or entity reference)
	return dslType
}

// Utility functions
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

func toPascalCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(string(part[0])) + part[1:]
		}
	}
	return strings.Join(parts, "")
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
