package parser

// Microservice represents the root DSL definition
type Microservice struct {
	Name         string
	Version      string
	Entities     []Entity
	ValueObjects []ValueObject
	Enums        []Enum
	Usecases     []Usecase
	GRPCService  *GRPCService
	Events       *Events
	Dependencies *Dependencies
	Config       map[string]interface{}
}

// Entity represents a domain entity
type Entity struct {
	Name   string
	Fields []Field
}

// ValueObject represents a value object
type ValueObject struct {
	Name   string
	Fields []Field
}

// Field represents a field in an entity or value object
type Field struct {
	Name         string
	Type         string
	Constraints  []string
	DefaultValue string
}

// Enum represents an enumeration type
type Enum struct {
	Name   string
	Values []string
}

// Usecase represents a business use case
type Usecase struct {
	Name        string
	Input       []Field
	Output      []Field
	Errors      []string
	Saga        *Saga
	Transaction bool
	RealTime    *RealTime
	ExternalAPI map[string]string
	Validation  map[string]string
	CustomImpl  []string
}

// Saga represents a saga pattern definition
type Saga struct {
	Steps        []SagaStep
	Compensates  map[string]string
}

// SagaStep represents a single step in a saga
type SagaStep struct {
	Name    string
	Service string
}

// RealTime represents real-time communication settings
type RealTime struct {
	Channel string
	Event   string
}

// GRPCService represents gRPC service definition
type GRPCService struct {
	Methods []GRPCMethod
}

// GRPCMethod represents a single gRPC method
type GRPCMethod struct {
	Name    string
	Request string
	Response string
}

// Events represents event definitions
type Events struct {
	Publish   []Event
	Subscribe []Event
}

// Event represents a single event
type Event struct {
	Name   string
	Fields []Field
}

// Dependencies represents service dependencies
type Dependencies struct {
	Database  string
	Cache     string
	Messaging string
	Services  map[string]string
}
