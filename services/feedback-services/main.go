package main

import (
	"log"
	"net"
	"os"

	pb "feedback-services/proto"
	"feedback-services/server"

	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50059"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()

	// Create feedback server instance
	feedbackServer := server.NewFeedbackServer()

	// Register the feedback service with the gRPC server
	pb.RegisterFeedbackServiceServer(s, feedbackServer)

	log.Printf("Feedback service is running on port %s...", port)

	// Start the server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}