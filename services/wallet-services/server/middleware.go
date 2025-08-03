package server

import (
	"context"
	"log"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthInterceptor validates tokens for protected endpoints
func AuthInterceptor(userClient UserServiceClient) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Skip validation for health checks or other public endpoints if needed
		// For now, all wallet endpoints require authentication
		
		// Extract metadata from context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "Missing metadata")
		}

		// Get authorization header
		auth := md.Get("authorization")
		if len(auth) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "Missing authorization header")
		}

		// Extract token from Bearer format
		token := auth[0]
		if strings.HasPrefix(token, "Bearer ") {
			token = strings.TrimPrefix(token, "Bearer ")
		}

		// Validate token with user service
		claims, err := userClient.ValidateToken(ctx, token)
		if err != nil {
			log.Printf("Token validation failed: %v", err)
			return nil, status.Errorf(codes.Unauthenticated, "Invalid token")
		}

		// Add user claims to context for use in handlers
		ctx = context.WithValue(ctx, "userClaims", claims)
		ctx = context.WithValue(ctx, "userID", claims.UserID)
		ctx = context.WithValue(ctx, "userEmail", claims.Email)
		ctx = context.WithValue(ctx, "userRole", claims.Role)

		// Call the handler with the enriched context
		return handler(ctx, req)
	}
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return "", status.Errorf(codes.Unauthenticated, "User ID not found in context")
	}
	return userID, nil
}

// GetUserClaimsFromContext extracts user claims from context
func GetUserClaimsFromContext(ctx context.Context) (*TokenClaims, error) {
	claims, ok := ctx.Value("userClaims").(*TokenClaims)
	if !ok || claims == nil {
		return nil, status.Errorf(codes.Unauthenticated, "User claims not found in context")
	}
	return claims, nil
}