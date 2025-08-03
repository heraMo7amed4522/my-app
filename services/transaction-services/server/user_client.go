package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	userpb "user-services/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// UserServiceClient interface for user service operations
type UserServiceClient interface {
	ValidateToken(ctx context.Context, token string) (*TokenClaims, error)
	Close() error
}

// TokenClaims represents the JWT token claims
type TokenClaims struct {
	UserID string
	Email  string
	Role   string
	Exp    int64
}

// userServiceClient implements UserServiceClient
type userServiceClient struct {
	conn   *grpc.ClientConn
	client userpb.UserServiceClient
}

// NewUserServiceClient creates a new user service client
func NewUserServiceClient() (UserServiceClient, error) {
	userServiceAddr := os.Getenv("USER_SERVICE_ADDR")
	if userServiceAddr == "" {
		userServiceAddr = "localhost:50051"
	}

	conn, err := grpc.Dial(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %v", err)
	}

	client := userpb.NewUserServiceClient(conn)

	return &userServiceClient{
		conn:   conn,
		client: client,
	}, nil
}

// ValidateToken validates a JWT token using the user service
func (c *userServiceClient) ValidateToken(ctx context.Context, token string) (*TokenClaims, error) {
	log.Printf("Validating token with user service")

	if token == "" {
		return nil, fmt.Errorf("token is empty")
	}

	// Remove "Bearer " prefix if present
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
	}

	// Make actual gRPC call to user service
	request := &userpb.ValidateTokenRequest{
		Token: token,
	}

	response, err := c.client.ValidateToken(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %v", err)
	}

	// Check if the response indicates success
	if response.StatusCode != 200 {
		if response.GetError() != nil {
			return nil, fmt.Errorf("token validation failed: %s", response.GetError().Message)
		}
		return nil, fmt.Errorf("token validation failed with status code: %d", response.StatusCode)
	}

	// Extract claims from successful response
	claims := response.GetClaims()
	if claims == nil {
		return nil, fmt.Errorf("no claims found in response")
	}

	return &TokenClaims{
		UserID: claims.UserId,
		Email:  claims.Email,
		Role:   claims.Role,
		Exp:    claims.Exp,
	}, nil
}

// Close closes the gRPC connection
func (c *userServiceClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// ValidateTokenMiddleware validates token and returns claims
func ValidateTokenMiddleware(userClient UserServiceClient, token string) (*TokenClaims, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	claims, err := userClient.ValidateToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("token validation failed: %v", err)
	}

	return claims, nil
}
