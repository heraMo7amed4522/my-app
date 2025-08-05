package main

import (
	"log"
	"net"

	pb "upload-services/proto"
	"upload-services/server"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50061")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()

	// Create upload server instance
	uploadServer := server.NewUploadServer()

	// Register the upload service with the gRPC server
	pb.RegisterUploadServiceServer(s, uploadServer)

	log.Println("Upload service is running on port 50061...")

	// Start the server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}