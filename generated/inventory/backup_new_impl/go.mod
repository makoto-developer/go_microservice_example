module github.com/makoto-developer/go_microservice_example/generated/inventory

go 1.25

require (
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
	github.com/makoto-developer/go_microservice_example/proto v0.0.0
	google.golang.org/grpc v1.72.1
	google.golang.org/protobuf v1.36.5
)

replace github.com/makoto-developer/go_microservice_example/proto => ../../proto
