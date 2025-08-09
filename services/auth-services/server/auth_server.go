package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "auth-services/proto"
	userpb "user-services/proto"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	firebaseAuth      *auth.Client
	userServiceAddr   string
	jwtSecret         string
	userServiceConn   *grpc.ClientConn
	userServiceClient userpb.UserServiceClient
}

func NewAuthServer(firebaseCredentials, userServiceAddr, jwtSecret string) (*AuthServer, error) {
	opt := option.WithCredentialsFile(firebaseCredentials)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %v", err)
	}
	authClient, err := app.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting firebase auth client: %v", err)
	}
	conn, err := grpc.NewClient(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %v", err)
	}
	userClient := userpb.NewUserServiceClient(conn)

	return &AuthServer{
		firebaseAuth:      authClient,
		userServiceAddr:   userServiceAddr,
		jwtSecret:         jwtSecret,
		userServiceConn:   conn,
		userServiceClient: userClient,
	}, nil
}

func (s *AuthServer) getUserByEmail(ctx context.Context, email string) (*userpb.User, error) {
	req := &userpb.GetUserByEmailRequest{
		Email: email,
	}
	resp, err := s.userServiceClient.GetUserByEmail(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %v", err)
	}
	switch result := resp.Result.(type) {
	case *userpb.GetUserByEmialResponse_User:
		return result.User, nil
	case *userpb.GetUserByEmialResponse_Error:
		if result.Error.Code == 404 {
			return nil, nil
		}
		return nil, fmt.Errorf("user service error: %s", result.Error.Message)
	default:
		return nil, fmt.Errorf("unexpected response type")
	}
}

func (s *AuthServer) createNewUser(ctx context.Context, fullName, email, password string) (*userpb.User, error) {
	req := &userpb.CreateNewUserRequest{
		FullName: fullName,
		Email:    email,
		Password: password,
	}
	resp, err := s.userServiceClient.CreateNewUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create new user: %v", err)
	}
	switch result := resp.Result.(type) {
	case *userpb.CreateNewUserResponse_User:
		return result.User, nil
	case *userpb.CreateNewUserResponse_Error:
		return nil, fmt.Errorf("user service error: %s", result.Error.Message)
	default:
		return nil, fmt.Errorf("unexpected response type")
	}
}

func (s *AuthServer) Close() error {
	if s.userServiceConn != nil {
		return s.userServiceConn.Close()
	}
	return nil
}

func (s *AuthServer) AuthenticateWithFirebase(ctx context.Context, req *pb.FirebaseAuthRequest) (*pb.FirebaseAuthResponse, error) {
	token, err := s.firebaseAuth.VerifyIDToken(ctx, req.IdToken)
	if err != nil {
		return &pb.FirebaseAuthResponse{
			StatusCode: 401,
			Message:    "Invalid Firebase token get from user side",
			Result: &pb.FirebaseAuthResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      401,
					Message:   err.Error(),
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	user, err := s.getUserByEmail(ctx, token.Claims["email"].(string))
	if err != nil {
		return &pb.FirebaseAuthResponse{
			StatusCode: 500,
			Message:    "Error fetching user",
			Result: &pb.FirebaseAuthResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   err.Error(),
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	if user == nil {
		user, err = s.createNewUser(ctx, token.Claims["name"].(string), token.Claims["email"].(string), "")
		if err != nil {
			return &pb.FirebaseAuthResponse{
				StatusCode: 500,
				Message:    "Error creating user",
				Result: &pb.FirebaseAuthResponse_Error{
					Error: &pb.ErrorDetails{
						Code:      500,
						Message:   err.Error(),
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			}, nil
		}
	}
	accessToken, err := generateTokens(user.Id, user.Email, user.Role)
	if err != nil {
		return &pb.FirebaseAuthResponse{
			StatusCode: 500,
			Message:    "Error generating tokens",
			Result: &pb.FirebaseAuthResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   err.Error(),
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.FirebaseAuthResponse{
		StatusCode: 200,
		Message:    "Authentication successful",
		Result: &pb.FirebaseAuthResponse_Tokens{
			Tokens: &pb.AuthTokens{
				AccessToken: accessToken,
				ExpiresIn:   3600,
				User: &pb.UserInfo{
					Id:       user.Id,
					Email:    user.Email,
					FullName: user.FullName,
					Role:     user.Role,
				},
			},
		},
	}, nil
}

func getJWTSecret() string {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-super-secret-jwt-key-change-this-in-production"
		log.Println("Warning: Using default JWT secret. Set JWT_SECRET environment variable in production.")
	}
	return jwtSecret
}
func generateTokens(userID, email, role string) (string, error) {
	jwtSecret := getJWTSecret()
	if jwtSecret == "" {
		jwtSecret = "your-super-secret-jwt-key-change-this-in-production"
		log.Println("Warning: Using default JWT secret. Set JWT_SECRET environment variable in production.")
	}
	accessClaims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return accessTokenString, nil
}
