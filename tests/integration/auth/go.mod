module github.com/makoto-developer/go_microservice_example/tests/integration/auth

go 1.27

replace github.com/makoto-developer/go_microservice_example/microservices/auth => ../../../microservices/auth

require (
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
	github.com/makoto-developer/go_microservice_example/microservices/auth v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.8.4
	google.golang.org/grpc v1.84.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.61.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/text v0.43.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260706201446-f0a921348800 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
