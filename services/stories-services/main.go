package main

import (
	"log"
	"net"

	pb "stories-services/proto"
	"stories-services/server"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50058")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()

	// Create stories server instance
	storiesServer := server.NewStoriesServer()

	// Register the stories service with the gRPC server
	pb.RegisterStoriesServiceServer(s, storiesServer)

	log.Println("Stories service is running on port 50058...")

	// Start the server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}