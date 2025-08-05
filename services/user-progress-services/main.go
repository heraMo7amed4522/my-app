package main

import (
	"log"
	"net"
	"os"

	pb "user-progress-services/proto"
	"user-progress-services/server"

	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50060"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()

	// Create user progress server instance
	userProgressServer := server.NewUserProgressServer()

	// Register the user progress service with the gRPC server
	pb.RegisterUserProgressServiceServer(s, userProgressServer)

	log.Printf("User Progress service is running on port %s...", port)

	// Start the server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}