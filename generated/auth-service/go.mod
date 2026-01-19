module github.com/makoto-developer/go_microservice_example/generated/auth-service

go 1.25

require (
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
	github.com/makoto-developer/go_microservice_example/manual/auth v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.27.0
	google.golang.org/grpc v1.68.1
	google.golang.org/protobuf v1.36.1
)

require (
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240903143218-8af14fe29dc1 // indirect
)

replace github.com/makoto-developer/go_microservice_example/manual/auth => ./manual/auth
