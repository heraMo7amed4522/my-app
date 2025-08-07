module chat-services

go 1.23.0

toolchain go1.24.5

require (
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
	google.golang.org/grpc v1.74.2
	google.golang.org/protobuf v1.36.6
	user-services v0.0.0-00010101000000-000000000000
)

replace user-services => ../user-services

require (
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250728155136-f173205681a0 // indirect
)
