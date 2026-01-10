package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Parser struct {
	lines      []string
	currentPos int
}

func NewParser(filePath string) (*Parser, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return &Parser{
		lines:      lines,
		currentPos: 0,
	}, nil
}

func (p *Parser) Parse() (*Microservice, error) {
	ms := &Microservice{
		Config: make(map[string]interface{}),
	}

	for p.currentPos < len(p.lines) {
		line := strings.TrimSpace(p.lines[p.currentPos])

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "//") {
			p.currentPos++
			continue
		}

		// Parse microservice declaration
		if strings.HasPrefix(line, "microservice ") {
			name := strings.TrimSuffix(strings.TrimSpace(strings.TrimPrefix(line, "microservice ")), "{")
			ms.Name = name
			p.currentPos++
			continue
		}

		// Parse version
		if strings.HasPrefix(line, "version:") {
			ms.Version = p.extractValue(line)
			p.currentPos++
			continue
		}

		// Parse entity
		if strings.HasPrefix(line, "entity ") {
			entity, err := p.parseEntity()
			if err != nil {
				return nil, err
			}
			ms.Entities = append(ms.Entities, entity)
			continue
		}

		// Parse value_object
		if strings.HasPrefix(line, "value_object ") {
			vo, err := p.parseValueObject()
			if err != nil {
				return nil, err
			}
			ms.ValueObjects = append(ms.ValueObjects, vo)
			continue
		}

		// Parse enum
		if strings.HasPrefix(line, "enum ") {
			enum, err := p.parseEnum()
			if err != nil {
				return nil, err
			}
			ms.Enums = append(ms.Enums, enum)
			continue
		}

		// Parse usecase
		if strings.HasPrefix(line, "usecase ") {
			usecase, err := p.parseUsecase()
			if err != nil {
				return nil, err
			}
			ms.Usecases = append(ms.Usecases, usecase)
			continue
		}

		// Parse grpc_service
		if strings.HasPrefix(line, "grpc_service") {
			grpcService, err := p.parseGRPCService()
			if err != nil {
				return nil, err
			}
			ms.GRPCService = grpcService
			continue
		}

		// Parse events
		if strings.HasPrefix(line, "events") {
			events, err := p.parseEvents()
			if err != nil {
				return nil, err
			}
			ms.Events = events
			continue
		}

		// Parse dependencies
		if strings.HasPrefix(line, "dependencies") {
			deps, err := p.parseDependencies()
			if err != nil {
				return nil, err
			}
			ms.Dependencies = deps
			continue
		}

		p.currentPos++
	}

	return ms, nil
}

func (p *Parser) parseEntity() (Entity, error) {
	line := strings.TrimSpace(p.lines[p.currentPos])
	nameMatch := regexp.MustCompile(`entity\s+(\w+)`).FindStringSubmatch(line)
	if len(nameMatch) < 2 {
		return Entity{}, fmt.Errorf("invalid entity declaration at line %d", p.currentPos)
	}

	entity := Entity{
		Name:   nameMatch[1],
		Fields: []Field{},
	}

	p.currentPos++

	// Parse fields
	for p.currentPos < len(p.lines) {
		line = strings.TrimSpace(p.lines[p.currentPos])

		if line == "}" {
			p.currentPos++
			break
		}

		if line == "" || strings.HasPrefix(line, "//") {
			p.currentPos++
			continue
		}

		field := p.parseField(line)
		entity.Fields = append(entity.Fields, field)
		p.currentPos++
	}

	return entity, nil
}

func (p *Parser) parseValueObject() (ValueObject, error) {
	line := strings.TrimSpace(p.lines[p.currentPos])
	nameMatch := regexp.MustCompile(`value_object\s+(\w+)`).FindStringSubmatch(line)
	if len(nameMatch) < 2 {
		return ValueObject{}, fmt.Errorf("invalid value_object declaration at line %d", p.currentPos)
	}

	vo := ValueObject{
		Name:   nameMatch[1],
		Fields: []Field{},
	}

	p.currentPos++

	for p.currentPos < len(p.lines) {
		line = strings.TrimSpace(p.lines[p.currentPos])

		if line == "}" {
			p.currentPos++
			break
		}

		if line == "" || strings.HasPrefix(line, "//") {
			p.currentPos++
			continue
		}

		field := p.parseField(line)
		vo.Fields = append(vo.Fields, field)
		p.currentPos++
	}

	return vo, nil
}

func (p *Parser) parseField(line string) Field {
	// Format: field_name: type constraints
	parts := strings.SplitN(line, ":", 2)
	if len(parts) < 2 {
		return Field{}
	}

	fieldName := strings.TrimSpace(parts[0])
	rest := strings.TrimSpace(parts[1])

	// Extract type and constraints
	tokens := strings.Fields(rest)
	if len(tokens) == 0 {
		return Field{}
	}

	fieldType := tokens[0]
	constraints := []string{}
	defaultValue := ""

	for i := 1; i < len(tokens); i++ {
		token := tokens[i]
		if strings.HasPrefix(token, "default(") {
			defaultValue = strings.TrimSuffix(strings.TrimPrefix(token, "default("), ")")
		} else {
			constraints = append(constraints, token)
		}
	}

	return Field{
		Name:         fieldName,
		Type:         fieldType,
		Constraints:  constraints,
		DefaultValue: defaultValue,
	}
}

func (p *Parser) parseEnum() (Enum, error) {
	line := strings.TrimSpace(p.lines[p.currentPos])
	nameMatch := regexp.MustCompile(`enum\s+(\w+)`).FindStringSubmatch(line)
	if len(nameMatch) < 2 {
		return Enum{}, fmt.Errorf("invalid enum declaration at line %d", p.currentPos)
	}

	enum := Enum{
		Name:   nameMatch[1],
		Values: []string{},
	}

	p.currentPos++

	for p.currentPos < len(p.lines) {
		line = strings.TrimSpace(p.lines[p.currentPos])

		if line == "}" {
			p.currentPos++
			break
		}

		if line == "" || strings.HasPrefix(line, "//") {
			p.currentPos++
			continue
		}

		// Parse enum values (comma-separated)
		values := strings.Split(line, ",")
		for _, val := range values {
			val = strings.TrimSpace(val)
			if val != "" {
				enum.Values = append(enum.Values, val)
			}
		}

		p.currentPos++
	}

	return enum, nil
}

func (p *Parser) parseUsecase() (Usecase, error) {
	line := strings.TrimSpace(p.lines[p.currentPos])
	nameMatch := regexp.MustCompile(`usecase\s+(\w+)`).FindStringSubmatch(line)
	if len(nameMatch) < 2 {
		return Usecase{}, fmt.Errorf("invalid usecase declaration at line %d", p.currentPos)
	}

	usecase := Usecase{
		Name:        nameMatch[1],
		Input:       []Field{},
		Output:      []Field{},
		Errors:      []string{},
		ExternalAPI: make(map[string]string),
		Validation:  make(map[string]string),
		CustomImpl:  []string{},
	}

	p.currentPos++

	for p.currentPos < len(p.lines) {
		line = strings.TrimSpace(p.lines[p.currentPos])

		if line == "}" {
			p.currentPos++
			break
		}

		if line == "" || strings.HasPrefix(line, "//") {
			p.currentPos++
			continue
		}

		// Parse input
		if strings.HasPrefix(line, "input:") {
			p.currentPos++
			usecase.Input = p.parseFieldBlock()
			continue
		}

		// Parse output
		if strings.HasPrefix(line, "output:") {
			p.currentPos++
			usecase.Output = p.parseFieldBlock()
			continue
		}

		// Parse errors
		if strings.HasPrefix(line, "errors:") {
			p.currentPos++
			usecase.Errors = p.parseErrorBlock()
			continue
		}

		// Parse saga
		if strings.HasPrefix(line, "saga:") {
			p.currentPos++
			saga := p.parseSagaBlock()
			usecase.Saga = &saga
			continue
		}

		// Parse transaction
		if strings.HasPrefix(line, "transaction:") {
			usecase.Transaction = strings.Contains(line, "true")
			p.currentPos++
			continue
		}

		p.currentPos++
	}

	return usecase, nil
}

func (p *Parser) parseFieldBlock() []Field {
	fields := []Field{}

	for p.currentPos < len(p.lines) {
		line := strings.TrimSpace(p.lines[p.currentPos])

		if line == "}" {
			break
		}

		if line == "" || strings.HasPrefix(line, "//") || line == "{" {
			p.currentPos++
			continue
		}

		field := p.parseField(line)
		fields = append(fields, field)
		p.currentPos++
	}

	return fields
}

func (p *Parser) parseErrorBlock() []string {
	errors := []string{}

	for p.currentPos < len(p.lines) {
		line := strings.TrimSpace(p.lines[p.currentPos])

		if line == "}" {
			break
		}

		if line == "" || strings.HasPrefix(line, "//") || line == "{" {
			p.currentPos++
			continue
		}

		// Parse error names (comma-separated)
		errorNames := strings.Split(line, ",")
		for _, err := range errorNames {
			err = strings.TrimSpace(err)
			if err != "" {
				errors = append(errors, err)
			}
		}

		p.currentPos++
	}

	return errors
}

func (p *Parser) parseSagaBlock() Saga {
	saga := Saga{
		Steps:       []SagaStep{},
		Compensates: make(map[string]string),
	}

	for p.currentPos < len(p.lines) {
		line := strings.TrimSpace(p.lines[p.currentPos])

		if line == "}" {
			break
		}

		if line == "" || strings.HasPrefix(line, "//") || line == "{" {
			p.currentPos++
			continue
		}

		// Parse step or compensate
		if strings.HasPrefix(line, "step") {
			stepName := p.extractValue(line)
			saga.Steps = append(saga.Steps, SagaStep{Name: stepName})
		} else if strings.HasPrefix(line, "compensate_") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				val := strings.TrimSpace(parts[1])
				saga.Compensates[key] = val
			}
		}

		p.currentPos++
	}

	return saga
}

func (p *Parser) parseGRPCService() (*GRPCService, error) {
	grpcService := &GRPCService{
		Methods: []GRPCMethod{},
	}

	p.currentPos++

	for p.currentPos < len(p.lines) {
		line := strings.TrimSpace(p.lines[p.currentPos])

		if line == "}" {
			p.currentPos++
			break
		}

		if line == "" || strings.HasPrefix(line, "//") || line == "{" {
			p.currentPos++
			continue
		}

		// Parse rpc method: rpc MethodName(Request) returns (Response)
		if strings.HasPrefix(line, "rpc ") {
			method := p.parseGRPCMethod(line)
			grpcService.Methods = append(grpcService.Methods, method)
		}

		p.currentPos++
	}

	return grpcService, nil
}

func (p *Parser) parseGRPCMethod(line string) GRPCMethod {
	rpcRegex := regexp.MustCompile(`rpc\s+(\w+)\((\w+)\)\s+returns\s+\((\w+)\)`)
	matches := rpcRegex.FindStringSubmatch(line)

	if len(matches) < 4 {
		return GRPCMethod{}
	}

	return GRPCMethod{
		Name:     matches[1],
		Request:  matches[2],
		Response: matches[3],
	}
}

func (p *Parser) parseEvents() (*Events, error) {
	events := &Events{
		Publish:   []Event{},
		Subscribe: []Event{},
	}

	p.currentPos++

	for p.currentPos < len(p.lines) {
		line := strings.TrimSpace(p.lines[p.currentPos])

		if line == "}" {
			p.currentPos++
			break
		}

		if line == "" || strings.HasPrefix(line, "//") || line == "{" {
			p.currentPos++
			continue
		}

		// Parse publish or subscribe
		if strings.HasPrefix(line, "publish ") {
			event := p.parseEvent(line, "publish")
			events.Publish = append(events.Publish, event)
		} else if strings.HasPrefix(line, "subscribe ") {
			event := p.parseEvent(line, "subscribe")
			events.Subscribe = append(events.Subscribe, event)
		}

		p.currentPos++
	}

	return events, nil
}

func (p *Parser) parseEvent(line, prefix string) Event {
	nameMatch := regexp.MustCompile(prefix + `\s+(\w+)`).FindStringSubmatch(line)
	if len(nameMatch) < 2 {
		return Event{}
	}

	event := Event{
		Name:   nameMatch[1],
		Fields: []Field{},
	}

	if strings.Contains(line, "{") {
		p.currentPos++
		event.Fields = p.parseFieldBlock()
	}

	return event
}

func (p *Parser) parseDependencies() (*Dependencies, error) {
	deps := &Dependencies{
		Services: make(map[string]string),
	}

	p.currentPos++

	for p.currentPos < len(p.lines) {
		line := strings.TrimSpace(p.lines[p.currentPos])

		if line == "}" {
			p.currentPos++
			break
		}

		if line == "" || strings.HasPrefix(line, "//") || line == "{" {
			p.currentPos++
			continue
		}

		// Parse dependency
		if strings.HasPrefix(line, "database:") {
			deps.Database = p.extractValue(line)
		} else if strings.HasPrefix(line, "cache:") {
			deps.Cache = p.extractValue(line)
		} else if strings.HasPrefix(line, "messaging:") {
			deps.Messaging = p.extractValue(line)
		}

		p.currentPos++
	}

	return deps, nil
}

func (p *Parser) extractValue(line string) string {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) < 2 {
		return ""
	}
	val := strings.TrimSpace(parts[1])
	val = strings.Trim(val, `"`)
	return val
}
