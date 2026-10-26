module github.com/makoto-developer/go_microservice_example/microservices/order

go 1.27

require (
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
	github.com/makoto-developer/go_microservice_example/microservices/inventory v0.0.0
	github.com/makoto-developer/go_microservice_example/microservices/notification v0.0.0
	github.com/makoto-developer/go_microservice_example/microservices/payment v0.0.0
	github.com/makoto-developer/go_microservice_example/microservices/shipping v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.84.0
	google.golang.org/protobuf v1.36.11
)

require (
	golang.org/x/net v0.61.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/text v0.43.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260706201446-f0a921348800 // indirect
)

replace github.com/makoto-developer/go_microservice_example/microservices/payment => ../payment

replace github.com/makoto-developer/go_microservice_example/microservices/shipping => ../shipping

replace github.com/makoto-developer/go_microservice_example/microservices/notification => ../notification

replace github.com/makoto-developer/go_microservice_example/microservices/inventory => ../inventory
