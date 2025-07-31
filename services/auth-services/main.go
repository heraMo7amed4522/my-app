package main

import (
	"log"
	"net"
	"os"

	"auth-services/server"
	pb "auth-services/proto"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()

	// Get environment variables
	firebaseCredentials := os.Getenv("FIREBASE_CREDENTIALS_PATH")
	userServiceAddr := os.Getenv("USER_SERVICE_ADDR") // e.g., "localhost:50051"
	jwtSecret := os.Getenv("JWT_SECRET")

	if firebaseCredentials == "" || userServiceAddr == "" || jwtSecret == "" {
		log.Fatal("Missing required environment variables")
	}

	// Create auth server with gRPC connection to user-services
	authServer, err := server.NewAuthServer(firebaseCredentials, userServiceAddr, jwtSecret)
	if err != nil {
		log.Fatalf("Failed to create auth server: %v", err)
	}
	defer authServer.Close() // Clean up gRPC connection

	// Register the auth service with the gRPC server
	pb.RegisterAuthServiceServer(s, authServer)

	log.Println("Auth service is running on port 50052...")

	// Start the server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
