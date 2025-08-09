package server

import (
	"context"
	"fmt"
	"time"

	pb "chat-services/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GRPCEndpointsService handles all non-streaming gRPC endpoints
type GRPCEndpointsService struct {
	*ChatServer
	endpointMetrics *EndpointMetrics
}

// EndpointMetrics tracks endpoint performance
type EndpointMetrics struct {
	RequestCount    map[string]int64
	ResponseTimes   map[string]time.Duration
	ErrorCount      map[string]int64
	SuccessCount    map[string]int64
	LastRequestTime map[string]time.Time
}

// NewGRPCEndpointsService creates a new endpoints service
func NewGRPCEndpointsService(chatServer *ChatServer) *GRPCEndpointsService {
	return &GRPCEndpointsService{
		ChatServer: chatServer,
		endpointMetrics: &EndpointMetrics{
			RequestCount:    make(map[string]int64),
			ResponseTimes:   make(map[string]time.Duration),
			ErrorCount:      make(map[string]int64),
			SuccessCount:    make(map[string]int64),
			LastRequestTime: make(map[string]time.Time),
		},
	}
}

// ================================================ MESSAGE ENDPOINTS =======================================

// SendMessageEndpoint handles message sending with enhanced features
func (s *GRPCEndpointsService) SendMessageEndpoint(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	start := time.Now()
	endpoint := "SendMessage"
	s.recordRequest(endpoint)

	defer func() {
		s.recordResponseTime(endpoint, time.Since(start))
	}()

	// Validate request
	if err := s.validateSendMessageRequest(req); err != nil {
		s.recordError(endpoint)
		return s.createErrorResponse(codes.InvalidArgument, "Invalid request", err), nil
	}

	// Call the original SendMessage method
	resp, err := s.ChatServer.SendMessage(ctx, req)
	if err != nil {
		s.recordError(endpoint)
		return nil, err
	}

	s.recordSuccess(endpoint)
	return resp, nil
}

// EditMessageEndpoint handles message editing with validation
func (s *GRPCEndpointsService) EditMessageEndpoint(ctx context.Context, req *pb.EditMessageRequest) (*pb.EditMessageResponse, error) {
	start := time.Now()
	endpoint := "EditMessage"
	s.recordRequest(endpoint)

	defer func() {
		s.recordResponseTime(endpoint, time.Since(start))
	}()

	// Validate request
	if err := s.validateEditMessageRequest(req); err != nil {
		s.recordError(endpoint)
		return s.createEditMessageErrorResponse(codes.InvalidArgument, "Invalid request", err), nil
	}

	// Call the original EditMessage method
	resp, err := s.ChatServer.EditMessage(ctx, req)
	if err != nil {
		s.recordError(endpoint)
		return nil, err
	}

	s.recordSuccess(endpoint)
	return resp, nil
}

// DeleteMessageEndpoint handles message deletion with authorization
func (s *GRPCEndpointsService) DeleteMessageEndpoint(ctx context.Context, req *pb.DeleteMessageRequest) (*pb.DeleteMessageResponse, error) {
	start := time.Now()
	endpoint := "DeleteMessage"
	s.recordRequest(endpoint)

	defer func() {
		s.recordResponseTime(endpoint, time.Since(start))
	}()

	// Validate request
	if err := s.validateDeleteMessageRequest(req); err != nil {
		s.recordError(endpoint)
		return s.createDeleteMessageErrorResponse(codes.InvalidArgument, "Invalid request", err), nil
	}

	// Call the original DeleteMessage method
	resp, err := s.ChatServer.DeleteMessage(ctx, req)
	if err != nil {
		s.recordError(endpoint)
		return nil, err
	}

	s.recordSuccess(endpoint)
	return resp, nil
}

// ================================================ CHAT HISTORY ENDPOINTS =======================================

// GetChatHistoryEndpoint handles chat history retrieval with pagination
func (s *GRPCEndpointsService) GetChatHistoryEndpoint(ctx context.Context, req *pb.GetChatHistoryRequest) (*pb.GetChatHistoryResponse, error) {
	start := time.Now()
	endpoint := "GetChatHistory"
	s.recordRequest(endpoint)

	defer func() {
		s.recordResponseTime(endpoint, time.Since(start))
	}()

	// Validate pagination parameters
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 20 // Default limit
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	// Call the original GetChatHistory method
	resp, err := s.ChatServer.GetChatHistory(ctx, req)
	if err != nil {
		s.recordError(endpoint)
		return nil, err
	}

	s.recordSuccess(endpoint)
	return resp, nil
}

// SearchMessagesEndpoint handles message search with advanced filters
func (s *GRPCEndpointsService) SearchMessagesEndpoint(ctx context.Context, req *pb.SearchMessagesRequest) (*pb.SearchMessagesResponse, error) {
	start := time.Now()
	endpoint := "SearchMessages"
	s.recordRequest(endpoint)

	defer func() {
		s.recordResponseTime(endpoint, time.Since(start))
	}()

	// Validate search query
	if len(req.Query) < 2 {
		s.recordError(endpoint)
		return &pb.SearchMessagesResponse{
			StatusCode: 400,
			Message:    "Search query too short",
			Result: &pb.SearchMessagesResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      400,
					Message:   "Search query must be at least 2 characters",
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Call the original SearchMessages method
	resp, err := s.ChatServer.SearchMessages(ctx, req)
	if err != nil {
		s.recordError(endpoint)
		return nil, err
	}

	s.recordSuccess(endpoint)
	return resp, nil
}

// ================================================ GROUP MANAGEMENT ENDPOINTS =======================================

// CreateGroupEndpoint handles group creation with validation
func (s *GRPCEndpointsService) CreateGroupEndpoint(ctx context.Context, req *pb.CreateGroupRequest) (*pb.CreateGroupResponse, error) {
	start := time.Now()
	endpoint := "CreateGroup"
	s.recordRequest(endpoint)

	defer func() {
		s.recordResponseTime(endpoint, time.Since(start))
	}()

	// Validate group data
	if err := s.validateCreateGroupRequest(req); err != nil {
		s.recordError(endpoint)
		return s.createCreateGroupErrorResponse(codes.InvalidArgument, "Invalid request", err), nil
	}

	// Call the original CreateGroup method
	resp, err := s.ChatServer.CreateGroup(ctx, req)
	if err != nil {
		s.recordError(endpoint)
		return nil, err
	}

	s.recordSuccess(endpoint)
	return resp, nil
}

// JoinGroupEndpoint handles group joining with permissions check
func (s *GRPCEndpointsService) JoinGroupEndpoint(ctx context.Context, req *pb.JoinGroupRequest) (*pb.JoinGroupResponse, error) {
	start := time.Now()
	endpoint := "JoinGroup"
	s.recordRequest(endpoint)

	defer func() {
		s.recordResponseTime(endpoint, time.Since(start))
	}()

	// Call the original JoinGroup method
	resp, err := s.ChatServer.JoinGroup(ctx, req)
	if err != nil {
		s.recordError(endpoint)
		return nil, err
	}

	s.recordSuccess(endpoint)
	return resp, nil
}

// ================================================ CALL MANAGEMENT ENDPOINTS =======================================

// InitiateCallEndpoint handles call initiation with validation
func (s *GRPCEndpointsService) InitiateCallEndpoint(ctx context.Context, req *pb.InitiateCallRequest) (*pb.InitiateCallResponse, error) {
	start := time.Now()
	endpoint := "InitiateCall"
	s.recordRequest(endpoint)

	defer func() {
		s.recordResponseTime(endpoint, time.Since(start))
	}()

	// Validate call request
	if err := s.validateInitiateCallRequest(req); err != nil {
		s.recordError(endpoint)
		return s.createInitiateCallErrorResponse(codes.InvalidArgument, "Invalid request", err), nil
	}

	// Call the original InitiateCall method
	resp, err := s.ChatServer.InitiateCall(ctx, req)
	if err != nil {
		s.recordError(endpoint)
		return nil, err
	}

	s.recordSuccess(endpoint)
	return resp, nil
}

// AcceptCallEndpoint handles call acceptance
func (s *GRPCEndpointsService) AcceptCallEndpoint(ctx context.Context, req *pb.AcceptCallRequest) (*pb.AcceptCallResponse, error) {
	start := time.Now()
	endpoint := "AcceptCall"
	s.recordRequest(endpoint)

	defer func() {
		s.recordResponseTime(endpoint, time.Since(start))
	}()

	// Call the original AcceptCall method
	resp, err := s.ChatServer.AcceptCall(ctx, req)
	if err != nil {
		s.recordError(endpoint)
		return nil, err
	}

	s.recordSuccess(endpoint)
	return resp, nil
}

// ================================================ NOTIFICATION ENDPOINTS =======================================

// AddNotificationEndpoint handles notification creation
func (s *GRPCEndpointsService) AddNotificationEndpoint(ctx context.Context, req *pb.AddNotificationRequest) (*pb.AddNotificationResponse, error) {
	start := time.Now()
	endpoint := "AddNotification"
	s.recordRequest(endpoint)

	defer func() {
		s.recordResponseTime(endpoint, time.Since(start))
	}()

	// Validate notification data
	if err := s.validateAddNotificationRequest(req); err != nil {
		s.recordError(endpoint)
		return s.createAddNotificationErrorResponse(codes.InvalidArgument, "Invalid request", err), nil
	}

	// Call the original AddNotification method
	resp, err := s.ChatServer.AddNotification(ctx, req)
	if err != nil {
		s.recordError(endpoint)
		return nil, err
	}

	s.recordSuccess(endpoint)
	return resp, nil
}

// ================================================ VALIDATION METHODS =======================================

// validateSendMessageRequest validates send message request
func (s *GRPCEndpointsService) validateSendMessageRequest(req *pb.SendMessageRequest) error {
	if req.SenderId == "" {
		return fmt.Errorf("sender ID is required")
	}
	if req.ReceiverId == "" && req.GroupId == "" {
		return fmt.Errorf("either receiver ID or group ID is required")
	}
	if req.Content == "" {
		return fmt.Errorf("message content is required")
	}
	if !ValidateMessageContent(req.Content) {
		return fmt.Errorf("invalid message content")
	}
	return nil
}

// validateEditMessageRequest validates edit message request
func (s *GRPCEndpointsService) validateEditMessageRequest(req *pb.EditMessageRequest) error {
	if req.MessageId == "" {
		return fmt.Errorf("message ID is required")
	}
	if req.UserId == "" {
		return fmt.Errorf("user ID is required")
	}
	if req.NewContent == "" {
		return fmt.Errorf("new content is required")
	}
	if !ValidateMessageContent(req.NewContent) {
		return fmt.Errorf("invalid message content")
	}
	return nil
}

// validateDeleteMessageRequest validates delete message request
func (s *GRPCEndpointsService) validateDeleteMessageRequest(req *pb.DeleteMessageRequest) error {
	if req.MessageId == "" {
		return fmt.Errorf("message ID is required")
	}
	if req.UserId == "" {
		return fmt.Errorf("user ID is required")
	}
	return nil
}

// validateCreateGroupRequest validates create group request
func (s *GRPCEndpointsService) validateCreateGroupRequest(req *pb.CreateGroupRequest) error {
	if req.CreatorId == "" {
		return fmt.Errorf("creator ID is required")
	}
	if req.GroupName == "" {
		return fmt.Errorf("group name is required")
	}
	if !ValidateGroupName(req.GroupName) {
		return fmt.Errorf("invalid group name")
	}
	if req.GroupName != "" && !ValidateGroupDescription(req.GroupName) {
		return fmt.Errorf("invalid group description")
	}
	return nil
}

// validateInitiateCallRequest validates initiate call request
func (s *GRPCEndpointsService) validateInitiateCallRequest(req *pb.InitiateCallRequest) error {
	if req.CallerId == "" {
		return fmt.Errorf("caller ID is required")
	}
	if req.ReceiverId == "" && req.GroupId == "" {
		return fmt.Errorf("either receiver ID or group ID is required")
	}
	return nil
}

// validateAddNotificationRequest validates add notification request
func (s *GRPCEndpointsService) validateAddNotificationRequest(req *pb.AddNotificationRequest) error {
	if req.Notification.SenderId == "" {
		return fmt.Errorf("user ID is required")
	}
	if req.Notification.Title == "" {
		return fmt.Errorf("notification title is required")
	}
	if req.Notification.Content == "" {
		return fmt.Errorf("notification content is required")
	}
	return nil
}

// ================================================ ERROR RESPONSE HELPERS =======================================

// createErrorResponse creates a generic error response
func (s *GRPCEndpointsService) createErrorResponse(code codes.Code, message string, err error) *pb.SendMessageResponse {
	return &pb.SendMessageResponse{
		StatusCode: int32(code),
		Message:    message,
		Result: &pb.SendMessageResponse_Error{
			Error: &pb.ErrorMessage{
				Code:      int32(code),
				Message:   err.Error(),
				Timestamp: timestamppb.New(time.Now()),
			},
		},
	}
}

// createEditMessageErrorResponse creates an edit message error response
func (s *GRPCEndpointsService) createEditMessageErrorResponse(code codes.Code, message string, err error) *pb.EditMessageResponse {
	return &pb.EditMessageResponse{
		StatusCode: int32(code),
		Message:    message,
		Result: &pb.EditMessageResponse_Error{
			Error: &pb.ErrorMessage{
				Code:      int32(code),
				Message:   err.Error(),
				Timestamp: timestamppb.New(time.Now()),
			},
		},
	}
}

// createDeleteMessageErrorResponse creates a delete message error response
func (s *GRPCEndpointsService) createDeleteMessageErrorResponse(code codes.Code, message string, err error) *pb.DeleteMessageResponse {
	return &pb.DeleteMessageResponse{
		StatusCode: int32(code),
		Message:    message,
		Result: &pb.DeleteMessageResponse_Error{
			Error: &pb.ErrorMessage{
				Code:      int32(code),
				Message:   err.Error(),
				Timestamp: timestamppb.New(time.Now()),
			},
		},
	}
}

// createCreateGroupErrorResponse creates a create group error response
func (s *GRPCEndpointsService) createCreateGroupErrorResponse(code codes.Code, message string, err error) *pb.CreateGroupResponse {
	return &pb.CreateGroupResponse{
		StatusCode: int32(code),
		Message:    message,
		Result: &pb.CreateGroupResponse_Error{
			Error: &pb.ErrorMessage{
				Code:      int32(code),
				Message:   err.Error(),
				Timestamp: timestamppb.New(time.Now()),
			},
		},
	}
}

// createInitiateCallErrorResponse creates an initiate call error response
func (s *GRPCEndpointsService) createInitiateCallErrorResponse(code codes.Code, message string, err error) *pb.InitiateCallResponse {
	return &pb.InitiateCallResponse{
		StatusCode: int32(code),
		Message:    message,
		Result: &pb.InitiateCallResponse_Error{
			Error: &pb.ErrorMessage{
				Code:      int32(code),
				Message:   err.Error(),
				Timestamp: timestamppb.New(time.Now()),
			},
		},
	}
}

// createAddNotificationErrorResponse creates an add notification error response
func (s *GRPCEndpointsService) createAddNotificationErrorResponse(code codes.Code, message string, err error) *pb.AddNotificationResponse {
	return &pb.AddNotificationResponse{
		StatusCode: int32(code),
		Message:    message,
		Result: &pb.AddNotificationResponse_Error{
			Error: &pb.ErrorMessage{
				Code:      int32(code),
				Message:   err.Error(),
				Timestamp: timestamppb.New(time.Now()),
			},
		},
	}
}

// ================================================ METRICS METHODS =======================================

// recordRequest records a new request
func (s *GRPCEndpointsService) recordRequest(endpoint string) {
	s.endpointMetrics.RequestCount[endpoint]++
	s.endpointMetrics.LastRequestTime[endpoint] = time.Now()
}

// recordSuccess records a successful response
func (s *GRPCEndpointsService) recordSuccess(endpoint string) {
	s.endpointMetrics.SuccessCount[endpoint]++
}

// recordError records an error response
func (s *GRPCEndpointsService) recordError(endpoint string) {
	s.endpointMetrics.ErrorCount[endpoint]++
}

// recordResponseTime records response time for an endpoint
func (s *GRPCEndpointsService) recordResponseTime(endpoint string, duration time.Duration) {
	s.endpointMetrics.ResponseTimes[endpoint] = duration
}

// GetEndpointMetrics returns current endpoint metrics
func (s *GRPCEndpointsService) GetEndpointMetrics() *EndpointMetrics {
	return s.endpointMetrics
}

// GetHealthStatus returns the health status of the service
func (s *GRPCEndpointsService) GetHealthStatus() map[string]interface{} {
	return map[string]interface{}{
		"status":           "healthy",
		"timestamp":        time.Now(),
		"total_requests":   s.getTotalRequests(),
		"total_errors":     s.getTotalErrors(),
		"success_rate":     s.getSuccessRate(),
		"average_response": s.getAverageResponseTime(),
	}
}

// Helper methods for health status
func (s *GRPCEndpointsService) getTotalRequests() int64 {
	var total int64
	for _, count := range s.endpointMetrics.RequestCount {
		total += count
	}
	return total
}

func (s *GRPCEndpointsService) getTotalErrors() int64 {
	var total int64
	for _, count := range s.endpointMetrics.ErrorCount {
		total += count
	}
	return total
}

func (s *GRPCEndpointsService) getSuccessRate() float64 {
	totalRequests := s.getTotalRequests()
	if totalRequests == 0 {
		return 100.0
	}
	totalErrors := s.getTotalErrors()
	return float64(totalRequests-totalErrors) / float64(totalRequests) * 100.0
}

func (s *GRPCEndpointsService) getAverageResponseTime() time.Duration {
	var total time.Duration
	count := 0
	for _, duration := range s.endpointMetrics.ResponseTimes {
		total += duration
		count++
	}
	if count == 0 {
		return 0
	}
	return total / time.Duration(count)
}
