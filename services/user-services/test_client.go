package main

import (
	"context"
	"log"
	pb "user-services/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	log.Println("Testing CreateNewUser...")
	createReq := &pb.CreateNewUserRequest{
		FullName:    "John Doe",
		Email:       "john.doe@example.com",
		Password:    "SecurePass123",
		CountryCode: "+1",
		PhoneNumber: "+1234567890",
		Role:        "USER",
	}

	createResp, err := client.CreateNewUser(context.Background(), createReq)
	if err != nil {
		log.Printf("CreateNewUser failed: %v", err)
	} else {
		log.Printf("CreateNewUser response: %v", createResp)
	}

	log.Println("\nTesting GetUserByEmail...")
	getReq := &pb.GetUserByEmailRequest{
		Email: "john.doe@example.com",
	}

	getResp, err := client.GetUserByEmail(context.Background(), getReq)
	if err != nil {
		log.Printf("GetUserByEmail failed: %v", err)
	} else {
		log.Printf("GetUserByEmail response: %v", getResp)
	}

	log.Println("\nTesting LoginUser...")
	loginReq := &pb.LoginUserRequest{
		Email:    "john.doe@example.com",
		Password: "SecurePass123",
	}

	loginResp, err := client.LoginUser(context.Background(), loginReq)
	if err != nil {
		log.Printf("LoginUser failed: %v", err)
	} else {
		log.Printf("LoginUser response: %v", loginResp)
	}

	log.Println("\nTesting ForgetPassword...")
	forgetReq := &pb.ForgetPasswordRequest{
		Email: "john.doe@example.com",
	}

	forgetResp, err := client.ForgetPassword(context.Background(), forgetReq)
	if err != nil {
		log.Printf("ForgetPassword failed: %v", err)
	} else {
		log.Printf("ForgetPassword response: %v", forgetResp)
	}

	log.Println("\nTest completed!")
}
