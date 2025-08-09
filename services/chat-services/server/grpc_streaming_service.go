package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	pb "chat-services/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GRPCStreamingService handles all real-time streaming operations
type GRPCStreamingService struct {
	*ChatStreamServer
	mutex sync.RWMutex

	// Live streaming connections
	liveStreamConnections map[string]*LiveStreamConnection
	streamingRooms        map[string]*StreamingRoom

	// Broadcasting channels
	broadcastChannels map[string]chan *pb.ChatStreamEnvelope

	// Stream metrics
	streamMetrics *StreamMetrics
}

// LiveStreamConnection represents a live streaming connection
type LiveStreamConnection struct {
	UserID       string
	RoomID       string
	Stream       pb.ChatService_ChatStreamServer
	Context      context.Context
	Cancel       context.CancelFunc
	ConnectedAt  time.Time
	LastActivity time.Time
	IsActive     bool
}

// StreamingRoom represents a streaming room with multiple participants
type StreamingRoom struct {
	RoomID         string
	OwnerID        string
	Participants   map[string]*LiveStreamConnection
	CreatedAt      time.Time
	IsActive       bool
	MaxViewers     int32
	CurrentViewers int32
	mutex          sync.RWMutex
}

// StreamMetrics tracks streaming performance
type StreamMetrics struct {
	ActiveStreams    int64
	TotalConnections int64
	BytesTransferred int64
	AverageLatency   time.Duration
	ErrorCount       int64
	mutex            sync.RWMutex
}

// NewGRPCStreamingService creates a new streaming service
func NewGRPCStreamingService(chatStreamServer *ChatStreamServer) *GRPCStreamingService {
	return &GRPCStreamingService{
		ChatStreamServer:      chatStreamServer,
		liveStreamConnections: make(map[string]*LiveStreamConnection),
		streamingRooms:        make(map[string]*StreamingRoom),
		broadcastChannels:     make(map[string]chan *pb.ChatStreamEnvelope),
		streamMetrics:         &StreamMetrics{},
	}
}

// StartLiveStream initiates a live streaming session
func (s *GRPCStreamingService) StartLiveStream(stream pb.ChatService_ChatStreamServer) error {
	ctx := stream.Context()

	// Validate user authentication
	claims, err := s.validateStreamRequest(ctx)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "Authentication failed: %v", err)
	}

	userID := claims.UserID
	log.Printf("Starting live stream for user: %s", userID)

	// Create stream connection
	connection := &LiveStreamConnection{
		UserID:       userID,
		Stream:       stream,
		Context:      ctx,
		ConnectedAt:  time.Now(),
		LastActivity: time.Now(),
		IsActive:     true,
	}

	// Register connection
	s.registerLiveStreamConnection(userID, connection)
	defer s.unregisterLiveStreamConnection(userID)

	// Handle incoming messages
	for {
		select {
		case <-ctx.Done():
			log.Printf("Live stream context cancelled for user: %s", userID)
			return ctx.Err()
		default:
			envelope, err := stream.Recv()
			if err == io.EOF {
				log.Printf("Live stream ended for user: %s", userID)
				return nil
			}
			if err != nil {
				log.Printf("Error receiving live stream message: %v", err)
				return err
			}

			// Process streaming message
			if err := s.processLiveStreamMessage(userID, envelope, stream); err != nil {
				log.Printf("Error processing live stream message: %v", err)
				return err
			}

			// Update activity timestamp
			connection.LastActivity = time.Now()
		}
	}
}

// processLiveStreamMessage processes incoming live stream messages
func (s *GRPCStreamingService) processLiveStreamMessage(userID string, envelope *pb.ChatStreamEnvelope, stream pb.ChatService_ChatStreamServer) error {
	switch payload := envelope.Payload.(type) {
	case *pb.ChatStreamEnvelope_Message:
		return s.handleLiveStreamChatMessage(userID, payload.Message, stream)
	case *pb.ChatStreamEnvelope_State:
		return s.handleLiveStreamStateMessage(userID, payload.State, stream)
	case *pb.ChatStreamEnvelope_Error:
		log.Printf("Live stream error from user %s: %s", userID, payload.Error.Message)
		return nil
	default:
		return fmt.Errorf("unknown live stream message type")
	}
}

// handleLiveStreamChatMessage handles chat messages in live stream
func (s *GRPCStreamingService) handleLiveStreamChatMessage(userID string, message *pb.ChatMessage, stream pb.ChatService_ChatStreamServer) error {
	// Add timestamp and user info
	message.Timestamp = timestamppb.New(time.Now())
	message.SenderId = userID

	// Broadcast to all viewers in the room
	if roomID := s.getUserStreamingRoom(userID); roomID != "" {
		s.broadcastToStreamingRoom(roomID, &pb.ChatStreamEnvelope{
			Payload: &pb.ChatStreamEnvelope_Message{
				Message: message,
			},
		}, userID)
	}

	return nil
}

// handleLiveStreamStateMessage handles state messages in live stream
func (s *GRPCStreamingService) handleLiveStreamStateMessage(userID string, state *pb.StateMessage, stream pb.ChatService_ChatStreamServer) error {
	switch state.Message {
	case "JOIN_ROOM":
		return s.joinStreamingRoom(userID, state.GetMessage(), stream)
	case "LEAVE_ROOM":
		return s.leaveStreamingRoom(userID, state.GetMessage())
	case "START_BROADCAST":
		return s.startBroadcast(userID, state.GetMessage())
	case "STOP_BROADCAST":
		return s.stopBroadcast(userID, state.GetMessage())
	default:
		log.Printf("Unknown live stream state type: %s", state.GetMessage())
		return nil
	}
}

// joinStreamingRoom adds a user to a streaming room
func (s *GRPCStreamingService) joinStreamingRoom(userID, roomID string, stream pb.ChatService_ChatStreamServer) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Get or create room
	room, exists := s.streamingRooms[roomID]
	if !exists {
		room = &StreamingRoom{
			RoomID:       roomID,
			OwnerID:      userID,
			Participants: make(map[string]*LiveStreamConnection),
			CreatedAt:    time.Now(),
			IsActive:     true,
			MaxViewers:   1000, // Default max viewers
		}
		s.streamingRooms[roomID] = room
	}

	// Check room capacity
	if room.CurrentViewers >= room.MaxViewers {
		return status.Errorf(codes.ResourceExhausted, "Streaming room is full")
	}

	// Add user to room
	if connection, exists := s.liveStreamConnections[userID]; exists {
		connection.RoomID = roomID
		room.Participants[userID] = connection
		room.CurrentViewers++
	}

	// Notify other participants
	s.broadcastToStreamingRoom(roomID, &pb.ChatStreamEnvelope{
		Payload: &pb.ChatStreamEnvelope_State{
			State: &pb.StateMessage{
				//Type:      "USER_JOINED",
				Message:   userID,
				Timestamp: timestamppb.New(time.Now()),
			},
		},
	}, userID)

	log.Printf("User %s joined streaming room %s", userID, roomID)
	return nil
}

// leaveStreamingRoom removes a user from a streaming room
func (s *GRPCStreamingService) leaveStreamingRoom(userID, roomID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	room, exists := s.streamingRooms[roomID]
	if !exists {
		return fmt.Errorf("streaming room not found")
	}

	// Remove user from room
	delete(room.Participants, userID)
	room.CurrentViewers--

	// Clean up empty room
	if len(room.Participants) == 0 {
		delete(s.streamingRooms, roomID)
	}

	// Notify other participants
	s.broadcastToStreamingRoom(roomID, &pb.ChatStreamEnvelope{
		Payload: &pb.ChatStreamEnvelope_State{
			State: &pb.StateMessage{
				//Message:      "USER_LEFT",
				Message:   userID,
				Timestamp: timestamppb.New(time.Now()),
			},
		},
	}, userID)

	log.Printf("User %s left streaming room %s", userID, roomID)
	return nil
}

// broadcastToStreamingRoom broadcasts a message to all participants in a room
func (s *GRPCStreamingService) broadcastToStreamingRoom(roomID string, envelope *pb.ChatStreamEnvelope, excludeUserID string) {
	s.mutex.RLock()
	room, exists := s.streamingRooms[roomID]
	s.mutex.RUnlock()

	if !exists {
		return
	}

	room.mutex.RLock()
	defer room.mutex.RUnlock()

	for userID, connection := range room.Participants {
		if userID == excludeUserID {
			continue
		}

		if connection.IsActive {
			if err := connection.Stream.Send(envelope); err != nil {
				log.Printf("Failed to send message to user %s in room %s: %v", userID, roomID, err)
				// Mark connection as inactive
				connection.IsActive = false
			}
		}
	}
}

// registerLiveStreamConnection registers a new live stream connection
func (s *GRPCStreamingService) registerLiveStreamConnection(userID string, connection *LiveStreamConnection) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.liveStreamConnections[userID] = connection
	s.streamMetrics.TotalConnections++
	s.streamMetrics.ActiveStreams++

	log.Printf("Registered live stream connection for user: %s", userID)
}

// unregisterLiveStreamConnection removes a live stream connection
func (s *GRPCStreamingService) unregisterLiveStreamConnection(userID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if connection, exists := s.liveStreamConnections[userID]; exists {
		// Remove from streaming room if in one
		if connection.RoomID != "" {
			s.leaveStreamingRoom(userID, connection.RoomID)
		}

		delete(s.liveStreamConnections, userID)
		s.streamMetrics.ActiveStreams--
	}

	log.Printf("Unregistered live stream connection for user: %s", userID)
}

// getUserStreamingRoom gets the current streaming room for a user
func (s *GRPCStreamingService) getUserStreamingRoom(userID string) string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if connection, exists := s.liveStreamConnections[userID]; exists {
		return connection.RoomID
	}
	return ""
}

// startBroadcast starts broadcasting for a user
func (s *GRPCStreamingService) startBroadcast(userID, roomID string) error {
	log.Printf("Starting broadcast for user %s in room %s", userID, roomID)
	// Implementation for starting broadcast
	return nil
}

// stopBroadcast stops broadcasting for a user
func (s *GRPCStreamingService) stopBroadcast(userID, roomID string) error {
	log.Printf("Stopping broadcast for user %s in room %s", userID, roomID)
	// Implementation for stopping broadcast
	return nil
}

// GetStreamMetrics returns current streaming metrics
func (s *GRPCStreamingService) GetStreamMetrics() *StreamMetrics {
	s.streamMetrics.mutex.RLock()
	defer s.streamMetrics.mutex.RUnlock()

	return &StreamMetrics{
		ActiveStreams:    s.streamMetrics.ActiveStreams,
		TotalConnections: s.streamMetrics.TotalConnections,
		BytesTransferred: s.streamMetrics.BytesTransferred,
		AverageLatency:   s.streamMetrics.AverageLatency,
		ErrorCount:       s.streamMetrics.ErrorCount,
	}
}

// CleanupInactiveConnections removes inactive streaming connections
func (s *GRPCStreamingService) CleanupInactiveConnections(timeout time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now()
	for userID, connection := range s.liveStreamConnections {
		if now.Sub(connection.LastActivity) > timeout {
			log.Printf("Cleaning up inactive connection for user: %s", userID)
			if connection.Cancel != nil {
				connection.Cancel()
			}
			delete(s.liveStreamConnections, userID)
			s.streamMetrics.ActiveStreams--
		}
	}
}
