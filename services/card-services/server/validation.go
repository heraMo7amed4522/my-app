package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	userpb "user-services/proto"

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
		userServiceAddr = "localhost:50051"
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
		return nil, fmt.Errorf("token can not be empity")
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	log.Printf("Validating token: %s...", token[:min(len(token), 20)])
	if len(token) < 10 {
		return nil, fmt.Errorf("invaild token format")
	}
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}
	userClient := userpb.NewUserServiceClient(c.conn)
	response, err := userClient.ValidateToken(ctx, &userpb.ValidateTokenRequest{Token: token})
	if err != nil {
		return nil, fmt.Errorf("user Service Error: %v", err)
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("token ValidationError: %s", response.Message)
	}
	claims := response.GetClaims()
	if claims == nil {
		return nil, fmt.Errorf("no Clamis found in response")
	}
	return &TokenClaims{
		UserID: claims.UserId,
		Email:  claims.Email,
		Role:   claims.Role,
		Exp:    claims.Exp,
	}, nil
}
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
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}
	claims, err := userClient.ValidateToken(context.Background(), token)
	if err != nil {
		return nil, fmt.Errorf("token validation failed: %v", err)
	}
	if claims.Exp < time.Now().Unix() {
		return nil, fmt.Errorf("token has expired")
	}
	return claims, nil
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
