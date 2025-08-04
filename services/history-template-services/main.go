package main

import (
	"log"
	"net"

	pb "history-template-services/proto"
	"history-template-services/server"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50055")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()

	// Create history template server instance
	historyTemplateServer := server.NewHistoryTemplateServer()

	// Register the history template service with the gRPC server
	pb.RegisterHistoryTemplateServiceServer(s, historyTemplateServer)

	log.Println("History Template service is running on port 50055...")

	// Start the server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}