package server

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "user-services/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserServiceClient interface {
	ValidateToken(ctx context.Context, token string) (*TokenClaims, error)
	GetUserByEmail(ctx context.Context, email string) (*pb.User, error)
}

type userServiceClient struct {
	conn   *grpc.ClientConn
	client pb.UserServiceClient
}

type TokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Exp    int64  `json:"exp"`
}

func NewUserServiceClient() (UserServiceClient, error) {
	userServiceAddr := getEnv("USER_SERVICE_ADDR", "localhost:50051")

	conn, err := grpc.Dial(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %v", err)
	}

	return &userServiceClient{
		conn:   conn,
		client: pb.NewUserServiceClient(conn),
	}, nil
}

func (c *userServiceClient) ValidateToken(ctx context.Context, token string) (*TokenClaims, error) {
	log.Printf("Validating token with user service: %s", token)

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

func (c *userServiceClient) GetUserByEmail(ctx context.Context, email string) (*pb.User, error) {
	log.Printf("Getting user by email from user service: %s", email)

	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	req := &pb.GetUserByEmailRequest{
		Email: email,
	}

	resp, err := c.client.GetUserByEmail(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %v", err)
	}

	if resp.StatusCode != 200 {
		if errorDetails := resp.GetError(); errorDetails != nil {
			return nil, fmt.Errorf("get user failed: %s", errorDetails.Message)
		}
		return nil, fmt.Errorf("get user failed with status: %d", resp.StatusCode)
	}

	user := resp.GetUser()
	if user == nil {
		return nil, fmt.Errorf("no user returned from user service")
	}

	return user, nil
}

func (c *userServiceClient) Close() error {
	return c.conn.Close()
}

func ValidateTokenMiddleware(userClient UserServiceClient, token string) (*TokenClaims, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	claims, err := userClient.ValidateToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("token validation failed: %v", err)
	}

	return claims, nil
}
