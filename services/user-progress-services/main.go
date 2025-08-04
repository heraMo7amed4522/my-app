package main

import (
	"log"
	"net"

	pb "user-progress-services/proto"
	"user-progress-services/server"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50057")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()

	// Create user progress server instance
	userProgressServer := server.NewUserProgressServer()

	// Register the user progress service with the gRPC server
	pb.RegisterUserProgressServiceServer(s, userProgressServer)

	log.Println("User Progress service is running on port 50057...")

	// Start the server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}