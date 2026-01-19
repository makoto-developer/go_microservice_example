module github.com/makoto-developer/go_microservice_example/generated/review

go 1.25.0

require github.com/google/uuid v1.6.0

require (
	github.com/lib/pq v1.10.9 // indirect
	github.com/makoto-developer/go_microservice_example/proto v0.0.0-20260119152141-0e9092b6fa83 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251029180050-ab9386a59fda // indirect
	google.golang.org/grpc v1.78.0 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)

replace github.com/makoto-developer/go_microservice_example/proto => ../../proto
