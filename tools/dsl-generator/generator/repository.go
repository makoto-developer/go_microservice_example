package generator

import (
	"fmt"
	"strings"

	"github.com/makoto-developer/go_microservice_example/tools/dsl-generator/parser"
)

type RepositoryGenerator struct {
	service *parser.Microservice
}

func NewRepositoryGenerator(service *parser.Microservice) *RepositoryGenerator {
	return &RepositoryGenerator{service: service}
}

func (g *RepositoryGenerator) GenerateRepositories() map[string]string {
	files := make(map[string]string)

	for _, entity := range g.service.Entities {
		fileName := toSnakeCase(entity.Name) + "_repository.go"
		content := g.generateRepositoryFile(entity)
		files[fileName] = content
	}

	return files
}

func (g *RepositoryGenerator) generateRepositoryFile(entity parser.Entity) string {
	var sb strings.Builder

	sb.WriteString("package domain\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("\t\"context\"\n")
	sb.WriteString("\n")
	sb.WriteString("\t\"github.com/google/uuid\"\n")
	sb.WriteString(")\n\n")

	// Repository interface
	sb.WriteString(fmt.Sprintf("// %sRepository defines repository interface for %s\n", entity.Name, entity.Name))
	sb.WriteString(fmt.Sprintf("type %sRepository interface {\n", entity.Name))

	// CRUD methods
	sb.WriteString(fmt.Sprintf("\t// Create creates a new %s\n", entity.Name))
	sb.WriteString(fmt.Sprintf("\tCreate(ctx context.Context, %s *%s) error\n\n", toSnakeCase(entity.Name), entity.Name))

	sb.WriteString(fmt.Sprintf("\t// FindByID finds %s by ID\n", entity.Name))
	sb.WriteString(fmt.Sprintf("\tFindByID(ctx context.Context, id uuid.UUID) (*%s, error)\n\n", entity.Name))

	sb.WriteString(fmt.Sprintf("\t// Update updates %s\n", entity.Name))
	sb.WriteString(fmt.Sprintf("\tUpdate(ctx context.Context, %s *%s) error\n\n", toSnakeCase(entity.Name), entity.Name))

	sb.WriteString(fmt.Sprintf("\t// Delete deletes %s by ID\n", entity.Name))
	sb.WriteString("\tDelete(ctx context.Context, id uuid.UUID) error\n\n")

	sb.WriteString(fmt.Sprintf("\t// List lists all %s\n", entity.Name))
	sb.WriteString(fmt.Sprintf("\tList(ctx context.Context, limit, offset int) ([]*%s, error)\n", entity.Name))

	// Find by unique fields
	for _, field := range entity.Fields {
		if contains(field.Constraints, "unique") && field.Name != "id" {
			sb.WriteString("\n")
			sb.WriteString(fmt.Sprintf("\t// FindBy%s finds %s by %s\n",
				toPascalCase(field.Name), entity.Name, field.Name))
			goType := g.mapDSLTypeToGo(field.Type)
			sb.WriteString(fmt.Sprintf("\tFindBy%s(ctx context.Context, %s %s) (*%s, error)\n",
				toPascalCase(field.Name), toSnakeCase(field.Name), goType, entity.Name))
		}
	}

	sb.WriteString("}\n")

	return sb.String()
}

func (g *RepositoryGenerator) mapDSLTypeToGo(dslType string) string {
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
