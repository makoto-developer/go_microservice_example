module github.com/makoto-developer/go_microservice_example/microservices/order

go 1.25

require (
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
	github.com/makoto-developer/go_microservice_example/microservices/payment v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.72.1
	google.golang.org/protobuf v1.36.5
)

require (
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
)

replace github.com/makoto-developer/go_microservice_example/microservices/payment => ../payment
