package main

import (
	"log"
	"net"

	pb "feedback-services/proto"
	"feedback-services/server"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50058")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()

	// Create feedback server instance
	feedbackServer := server.NewFeedbackServer()

	// Register the feedback service with the gRPC server
	pb.RegisterFeedbackServiceServer(s, feedbackServer)

	log.Println("Feedback service is running on port 50058...")

	// Start the server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}