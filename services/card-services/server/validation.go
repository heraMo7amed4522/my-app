package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// UserServiceClient interface for user service communication
type UserServiceClient interface {
	ValidateToken(ctx context.Context, token string) (*TokenClaims, error)
	Close() error
}

// TokenClaims represents the validated token claims
type TokenClaims struct {
	UserID string
	Email  string
	Role   string
	Exp    int64
}

// GRPCUserServiceClient implements UserServiceClient using gRPC
type GRPCUserServiceClient struct {
	conn *grpc.ClientConn
}

// NewUserServiceClient creates a new user service client
func NewUserServiceClient() (UserServiceClient, error) {
	userServiceAddr := os.Getenv("USER_SERVICE_ADDR")
	if userServiceAddr == "" {
		userServiceAddr = "localhost:50051" // Default for development
	}

	conn, err := grpc.NewClient(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %v", err)
	}

	return &GRPCUserServiceClient{
		conn: conn,
	}, nil
}

// ValidateToken validates a JWT token by calling the user service
func (c *GRPCUserServiceClient) ValidateToken(ctx context.Context, token string) (*TokenClaims, error) {
	if token == "" {
		return nil, fmt.Errorf("token cannot be empty")
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// For development, we'll implement a simple validation
	// In production, this would call the actual user service gRPC method
	log.Printf("Validating token: %s...", token[:min(len(token), 20)])

	// Basic token format validation
	if len(token) < 10 {
		return nil, fmt.Errorf("invalid token format")
	}

	// Remove Bearer prefix if present
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// For development, simulate successful validation
	// In production, replace this with actual gRPC call to user service
	// Example:
	// userClient := userpb.NewUserServiceClient(c.conn)
	// response, err := userClient.ValidateToken(ctx, &userpb.ValidateTokenRequest{Token: token})
	// if err != nil {
	//     return nil, fmt.Errorf("user service validation failed: %v", err)
	// }
	// return &TokenClaims{
	//     UserID: response.Claims.UserId,
	//     Email:  response.Claims.Email,
	//     Role:   response.Claims.Role,
	//     Exp:    response.Claims.Exp,
	// }, nil

	// Mock validation for development
	return &TokenClaims{
		UserID: "550e8400-e29b-41d4-a716-446655440000", // Mock UUID
		Email:  "user@example.com",
		Role:   "USER",
		Exp:    time.Now().Add(time.Hour).Unix(),
	}, nil
}

// Close closes the gRPC connection
func (c *GRPCUserServiceClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// ValidateTokenMiddleware validates the token and extracts user claims
func ValidateTokenMiddleware(userClient UserServiceClient, token string) (*TokenClaims, error) {
	if token == "" {
		return nil, fmt.Errorf("authorization token is required")
	}

	// Remove "Bearer " prefix if present
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	claims, err := userClient.ValidateToken(context.Background(), token)
	if err != nil {
		return nil, fmt.Errorf("token validation failed: %v", err)
	}

	// Check if token is expired
	if claims.Exp < time.Now().Unix() {
		return nil, fmt.Errorf("token has expired")
	}

	return claims, nil
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}