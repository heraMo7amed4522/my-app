package server

import (
	"context"
	"fmt"
	"os"
	"time"

	pb "user-services/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserServiceClient interface {
	ValidateToken(ctx context.Context, token string) (*TokenClaims, error)
	Close() error
}

type TokenClaims struct {
	UserID string
	Email  string
	Role   string
	Exp    int64
}

type GRPCUserServiceClient struct {
	conn   *grpc.ClientConn
	client pb.UserServiceClient
}

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

func (c *GRPCUserServiceClient) ValidateToken(ctx context.Context, token string) (*TokenClaims, error) {
	if token == "" {
		return nil, fmt.Errorf("token is required")
	}
	req := &pb.ValidateTokenRequest{
		Token: token,
	}
	resp, err := c.client.ValidateToken(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %v", err)
	}
	if resp.StatusCode != 200 {
		if errorDetails := resp.GetError(); errorDetails != nil {
			return nil, fmt.Errorf("token validation failed: %s", errorDetails.Message)
		}
		return nil, fmt.Errorf("token validation failed with status: %d", resp.StatusCode)
	}
	claims := resp.GetClaims()
	if claims == nil {
		return nil, fmt.Errorf("no claims returned from user service")
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
