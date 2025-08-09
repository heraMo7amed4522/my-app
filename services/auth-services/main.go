package main

import (
	"log"
	"net"
	"os"

	pb "auth-services/proto"
	"auth-services/server"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	firebaseCredentials := os.Getenv("FIREBASE_CREDENTIALS_PATH")
	userServiceAddr := os.Getenv("USER_SERVICE_ADDR")
	jwtSecret := os.Getenv("JWT_SECRET")

	if firebaseCredentials == "" || userServiceAddr == "" || jwtSecret == "" {
		log.Fatal("Missing required environment variables")
	}
	authServer, err := server.NewAuthServer(firebaseCredentials, userServiceAddr, jwtSecret)
	if err != nil {
		log.Fatalf("Failed to create auth server: %v", err)
	}
	defer authServer.Close()
	pb.RegisterAuthServiceServer(s, authServer)
	log.Println("Auth service is running on port 50052...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
