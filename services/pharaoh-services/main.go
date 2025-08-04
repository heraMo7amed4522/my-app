package main

import (
	"log"
	"net"
	"os"

	pb "pharaoh-services/proto"
	"pharaoh-services/server"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Get port from environment or use default
	port := os.Getenv("PHARAOH_SERVICE_PORT")
	if port == "" {
		port = "50055"
	}

	// Create listener
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Initialize database connection
	pharaohServer, err := server.NewPharaohServer()
	if err != nil {
		log.Fatalf("Failed to create pharaoh server: %v", err)
	}
	defer pharaohServer.Close()

	// Register service
	pb.RegisterPharaohServiceServer(grpcServer, pharaohServer)

	log.Printf("Pharaoh gRPC server listening on port %s", port)

	// Start server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
