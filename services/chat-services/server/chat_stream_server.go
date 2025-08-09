package server

import (
	"context"
	"io"
	"log"
	"sync"
	"time"

	pb "chat-services/proto"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// StreamConnection represents an active streaming connection
type StreamConnection struct {
	UserID string
	Stream pb.ChatService_ChatStreamServer
	Ctx    context.Context
	Cancel context.CancelFunc
}

// ChatStreamServer handles all streaming operations
type ChatStreamServer struct {
	*ChatServer // Embed the existing ChatServer for database operations
	connections map[string]*StreamConnection
	mutex       sync.RWMutex
	channels    map[string][]string // userID -> list of connected userIDs for broadcasting

	// New subscription maps
	lastMessageSubscriptions map[string]chan *pb.ChatMessage
	userStatusSubscriptions  map[string]chan *pb.UserStatus
	userStatusWatchers       map[string][]string // userID -> list of subscribers watching this user

	// New subscription maps for enhanced functionality
	notificationSubscriptions  map[string]chan *pb.NotificationUpdate
	callUpdatesSubscriptions   map[string]chan *pb.CallUpdate
	callSignalingSubscriptions map[string]map[string]chan *pb.CallSignalingMessage // userID -> callID -> channel
	videoCallConnections       map[string]pb.ChatService_VideoCallStreamServer
	phoneCallConnections       map[string]pb.ChatService_PhoneCallStreamServer
	webrtcService              *WebRTCService
}

// Video call message handler
func (s *ChatStreamServer) handleVideoCallMessage(userID string, envelope *pb.VideoCallStreamEnvelope, stream pb.ChatService_VideoCallStreamServer) {
	switch payload := envelope.Payload.(type) {
	case *pb.VideoCallStreamEnvelope_VideoData:
		s.handleVideoData(userID, payload.VideoData, stream)
	case *pb.VideoCallStreamEnvelope_AudioData:
		s.handleAudioData(userID, payload.AudioData, stream)
	case *pb.VideoCallStreamEnvelope_Control:
		s.handleCallControl(userID, payload.Control, stream)
	case *pb.VideoCallStreamEnvelope_State:
		s.handleCallState(userID, payload.State, stream)
	case *pb.VideoCallStreamEnvelope_Error:
		log.Printf("Video call error from user %s: %s", userID, payload.Error.Message)
	}
}

// Phone call message handler
func (s *ChatStreamServer) handlePhoneCallMessage(userID string, envelope *pb.PhoneCallStreamEnvelope, stream pb.ChatService_PhoneCallStreamServer) {
	switch payload := envelope.Payload.(type) {
	case *pb.PhoneCallStreamEnvelope_AudioData:
		s.handleAudioData(userID, payload.AudioData, nil)
	case *pb.PhoneCallStreamEnvelope_Control:
		s.handleCallControl(userID, payload.Control, nil)
	case *pb.PhoneCallStreamEnvelope_State:
		s.handleCallState(userID, payload.State, nil)
	case *pb.PhoneCallStreamEnvelope_Error:
		log.Printf("Phone call error from user %s: %s", userID, payload.Error.Message)
	}
}

// Handle video data
func (s *ChatStreamServer) handleVideoData(userID string, videoData *pb.VideoCallData, stream pb.ChatService_VideoCallStreamServer) {
	log.Printf("Received video data from user %s for call %s", userID, videoData.CallId)
	s.broadcastVideoData(videoData, userID)
}

// Handle audio data
func (s *ChatStreamServer) handleAudioData(userID string, audioData *pb.AudioCallData, stream pb.ChatService_VideoCallStreamServer) {
	log.Printf("Received audio data from user %s for call %s", userID, audioData.CallId)
	s.broadcastAudioData(audioData, userID)
}

// broadcastCallControl broadcasts call control messages to all participants in a call
func (s *ChatStreamServer) broadcastCallControl(control *pb.CallControlMessage, senderID string) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Get call participants (excluding sender)
	participants := s.getCallParticipants(control.CallId, senderID)

	// Broadcast to video call connections
	for _, participantID := range participants {
		if stream, exists := s.videoCallConnections[participantID]; exists {
			envelope := &pb.VideoCallStreamEnvelope{
				Payload: &pb.VideoCallStreamEnvelope_Control{
					Control: control,
				},
			}
			if err := stream.Send(envelope); err != nil {
				log.Printf("Failed to send call control to user %s: %v", participantID, err)
			}
		}

		// Also broadcast to phone call connections
		if stream, exists := s.phoneCallConnections[participantID]; exists {
			envelope := &pb.PhoneCallStreamEnvelope{
				Payload: &pb.PhoneCallStreamEnvelope_Control{
					Control: control,
				},
			}
			if err := stream.Send(envelope); err != nil {
				log.Printf("Failed to send call control to user %s: %v", participantID, err)
			}
		}
	}
}

// updateCallStateFromControl updates call state based on control messages
func (s *ChatStreamServer) updateCallStateFromControl(control *pb.CallControlMessage) {
	// Update call state in database based on control type
	switch control.ControlType {
	case pb.CallControlType(pb.CallStatus_CALL_ENDED):
		// Update participant video status
		log.Printf("User turned video off in call %s", control.CallId)
	case pb.CallControlType(pb.CallStatus_CALL_REJECTED):
		// Update call participants
		log.Printf("User left call %s", control.CallId)
		// TODO: Implement database update logic
	default:
		log.Printf("Unknown control type: %s", control.ControlType)
	}
}

// broadcastCallState broadcasts call state updates to all participants
func (s *ChatStreamServer) broadcastCallState(state *pb.CallStateMessage, senderID string) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Get call participants (excluding sender)
	participants := s.getCallParticipants(state.CallId, senderID)

	// Broadcast to video call connections
	for _, participantID := range participants {
		if stream, exists := s.videoCallConnections[participantID]; exists {
			envelope := &pb.VideoCallStreamEnvelope{
				Payload: &pb.VideoCallStreamEnvelope_State{
					State: state,
				},
			}
			if err := stream.Send(envelope); err != nil {
				log.Printf("Failed to send call state to user %s: %v", participantID, err)
			}
		}

		// Also broadcast to phone call connections
		if stream, exists := s.phoneCallConnections[participantID]; exists {
			envelope := &pb.PhoneCallStreamEnvelope{
				Payload: &pb.PhoneCallStreamEnvelope_State{
					State: state,
				},
			}
			if err := stream.Send(envelope); err != nil {
				log.Printf("Failed to send call state to user %s: %v", participantID, err)
			}
		}
	}

	// Also broadcast to call updates subscribers
	if ch, exists := s.callUpdatesSubscriptions[state.CallId]; exists {
		update := &pb.CallUpdate{
			CallId:     state.CallId,
			UpdateType: pb.CallUpdateType(state.Status),
			Timestamp:  state.Timestamp,
			// Add other relevant fields as needed
		}
		select {
		case ch <- update:
		default:
			log.Printf("Failed to send call update for call %s", state.CallId)
		}
	}
}

// updateCallInDatabase updates call information in the database
func (s *ChatStreamServer) updateCallInDatabase(state *pb.CallStateMessage) {
	// Update call state in database
	switch state.Status {
	case pb.CallStatus_CALL_INITIATED:
		// Update call start time
		log.Printf("Call %s started at %v", state.CallId, state.Timestamp)
		// TODO: Implement database update logic
		// Example: s.ChatServer.UpdateCallStatus(state.CallId, "ACTIVE", state.Timestamp)
	case pb.CallStatus_CALL_ENDED:
		// Update call end time
		log.Printf("Call %s ended at %v", state.CallId, state.Timestamp)
		// TODO: Implement database update logic
		// Example: s.ChatServer.UpdateCallStatus(state.CallId, "ENDED", state.Timestamp)
	case pb.CallStatus_CALL_BUSY:
		// Update call pause status
		log.Printf("Call %s paused at %v", state.CallId, state.Timestamp)
		// TODO: Implement database update logic
	case pb.CallStatus_CALL_ACCEPTED:
		// Update call resume status
		log.Printf("Call %s resumed at %v", state.CallId, state.Timestamp)
		// TODO: Implement database update logic
	default:
		log.Printf("Unknown call status: %s for call %s", state.Status, state.CallId)
	}

	// You may also want to update participant information
	if len(state.Participants) > 0 {
		log.Printf("Updating participants for call %s: %v", state.CallId, state.Participants)
		// TODO: Implement participant update logic
	}
}

// Handle call control
func (s *ChatStreamServer) handleCallControl(userID string, control *pb.CallControlMessage, stream pb.ChatService_VideoCallStreamServer) {
	log.Printf("Received call control from user %s: %s", userID, control.ControlType)
	s.broadcastCallControl(control, userID)
	s.updateCallStateFromControl(control)
}

// Handle call state
func (s *ChatStreamServer) handleCallState(userID string, state *pb.CallStateMessage, stream pb.ChatService_VideoCallStreamServer) {
	log.Printf("Received call state from user %s: %s", userID, state.Status)
	s.broadcastCallState(state, userID)
	s.updateCallInDatabase(state)
}

// Notification subscription methods
func (s *ChatStreamServer) registerNotificationSubscription(userID string, ch chan *pb.NotificationUpdate) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.notificationSubscriptions[userID] = ch
}
func (s *ChatStreamServer) unregisterNotificationSubscription(userID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if ch, exists := s.notificationSubscriptions[userID]; exists {
		close(ch)
		delete(s.notificationSubscriptions, userID)
	}
}

// Broadcast video data to call participants
func (s *ChatStreamServer) broadcastVideoData(videoData *pb.VideoCallData, senderID string) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Get call participants (excluding sender)
	participants := s.getCallParticipants(videoData.CallId, senderID)

	envelope := &pb.VideoCallStreamEnvelope{
		Payload: &pb.VideoCallStreamEnvelope_VideoData{
			VideoData: videoData,
		},
	}

	for _, participantID := range participants {
		if stream, exists := s.videoCallConnections[participantID]; exists {
			go func(stream pb.ChatService_VideoCallStreamServer) {
				if err := stream.Send(envelope); err != nil {
					log.Printf("Failed to broadcast video data to participant %s: %v", participantID, err)
				}
			}(stream)
		}
	}
}

// Broadcast audio data to call participants
func (s *ChatStreamServer) broadcastAudioData(audioData *pb.AudioCallData, senderID string) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Get call participants (excluding sender)
	participants := s.getCallParticipants(audioData.CallId, senderID)

	// Broadcast to video call streams
	videoEnvelope := &pb.VideoCallStreamEnvelope{
		Payload: &pb.VideoCallStreamEnvelope_AudioData{
			AudioData: audioData,
		},
	}

	// Broadcast to phone call streams
	phoneEnvelope := &pb.PhoneCallStreamEnvelope{
		Payload: &pb.PhoneCallStreamEnvelope_AudioData{
			AudioData: audioData,
		},
	}

	for _, participantID := range participants {
		// Try video call connection first
		if stream, exists := s.videoCallConnections[participantID]; exists {
			go func(stream pb.ChatService_VideoCallStreamServer) {
				if err := stream.Send(videoEnvelope); err != nil {
					log.Printf("Failed to broadcast audio data to video participant %s: %v", participantID, err)
				}
			}(stream)
		} else if stream, exists := s.phoneCallConnections[participantID]; exists {
			// Fallback to phone call connection
			go func(stream pb.ChatService_PhoneCallStreamServer) {
				if err := stream.Send(phoneEnvelope); err != nil {
					log.Printf("Failed to broadcast audio data to phone participant %s: %v", participantID, err)
				}
			}(stream)
		}
	}
}

// Get call participants excluding sender
func (s *ChatStreamServer) getCallParticipants(callID, excludeUserID string) []string {
	// This should query the database for call participants
	// For now, returning mock data
	return []string{} // Implement based on your database schema
}

// Send initial notifications
func (s *ChatStreamServer) sendInitialNotifications(userID string, stream pb.ChatService_SubscribeToNotificationsServer) {
	// Query database for recent unread notifications
	// This is a placeholder implementation
	log.Printf("Sending initial notifications to user %s", userID)
}

// Broadcast notification to subscribers
func (s *ChatStreamServer) broadcastNotification(notification *pb.NotificationUpdate) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if ch, exists := s.notificationSubscriptions[notification.SenderId]; exists {
		select {
		case ch <- notification:
		default:
			log.Printf("Notification channel full for user %s", notification.SenderId)
		}
	}
}

// Call updates subscription methods
func (s *ChatStreamServer) registerCallUpdatesSubscription(userID string, callIDs []string, ch chan *pb.CallUpdate) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.callUpdatesSubscriptions[userID] = ch
}

func (s *ChatStreamServer) unregisterCallUpdatesSubscription(userID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if ch, exists := s.callUpdatesSubscriptions[userID]; exists {
		close(ch)
		delete(s.callUpdatesSubscriptions, userID)
	}
}

// Call signaling subscription methods
func (s *ChatStreamServer) registerCallSignalingSubscription(userID, callID string, ch chan *pb.CallSignalingMessage) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.callSignalingSubscriptions[userID] == nil {
		s.callSignalingSubscriptions[userID] = make(map[string]chan *pb.CallSignalingMessage)
	}
	s.callSignalingSubscriptions[userID][callID] = ch
}
func (s *ChatStreamServer) unregisterCallSignalingSubscription(userID, callID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if userCalls, exists := s.callSignalingSubscriptions[userID]; exists {
		if ch, exists := userCalls[callID]; exists {
			close(ch)
			delete(userCalls, callID)
		}
		if len(userCalls) == 0 {
			delete(s.callSignalingSubscriptions, userID)
		}
	}
}

// Video call connection methods
func (s *ChatStreamServer) registerVideoCallConnection(userID string, stream pb.ChatService_VideoCallStreamServer) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.videoCallConnections[userID] = stream
}
func (s *ChatStreamServer) unregisterVideoCallConnection(userID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.videoCallConnections, userID)
}

// Phone call connection methods
func (s *ChatStreamServer) registerPhoneCallConnection(userID string, stream pb.ChatService_PhoneCallStreamServer) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.phoneCallConnections[userID] = stream
}
func (s *ChatStreamServer) unregisterPhoneCallConnection(userID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.phoneCallConnections, userID)
}

// NewChatStreamServer creates a new streaming server
func NewChatStreamServer() *ChatStreamServer {
	return &ChatStreamServer{
		ChatServer:                 NewChatServer(),
		connections:                make(map[string]*StreamConnection),
		channels:                   make(map[string][]string),
		lastMessageSubscriptions:   make(map[string]chan *pb.ChatMessage),
		userStatusSubscriptions:    make(map[string]chan *pb.UserStatus),
		userStatusWatchers:         make(map[string][]string),
		notificationSubscriptions:  make(map[string]chan *pb.NotificationUpdate),
		callUpdatesSubscriptions:   make(map[string]chan *pb.CallUpdate),
		callSignalingSubscriptions: make(map[string]map[string]chan *pb.CallSignalingMessage),
		videoCallConnections:       make(map[string]pb.ChatService_VideoCallStreamServer),
		phoneCallConnections:       make(map[string]pb.ChatService_PhoneCallStreamServer),
	}
}

// ChatStream implements the main bidirectional streaming RPC
func (s *ChatStreamServer) ChatStream(stream pb.ChatService_ChatStreamServer) error {
	ctx := stream.Context()

	// Extract and validate token
	claims, err := s.validateStreamRequest(ctx)
	if err != nil {
		return s.sendStreamError(stream, 401, "Unauthorized", err.Error())
	}

	userID := claims.UserID
	log.Printf("User %s connected to ChatStream", userID)

	// Create stream context with cancellation
	streamCtx, cancel := context.WithCancel(ctx)
	connection := &StreamConnection{
		UserID: userID,
		Stream: stream,
		Ctx:    streamCtx,
		Cancel: cancel,
	}

	// Register connection
	s.registerConnection(userID, connection)
	defer s.unregisterConnection(userID)
	defer cancel()

	// Send welcome message
	welcomeMsg := &pb.ChatStreamEnvelope{
		Payload: &pb.ChatStreamEnvelope_State{
			State: &pb.StateMessage{
				StatusCode: 200,
				Message:    "Connected to ChatStream",
				Timestamp:  timestamppb.New(time.Now()),
			},
		},
	}
	if err := stream.Send(welcomeMsg); err != nil {
		return err
	}

	// Handle incoming messages
	for {
		select {
		case <-streamCtx.Done():
			return streamCtx.Err()
		default:
			envelope, err := stream.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				log.Printf("Stream receive error: %v", err)
				return err
			}

			// Process the received message
			go s.handleStreamMessage(userID, envelope, stream)
		}
	}
}

// handleStreamMessage processes different types of streaming messages
func (s *ChatStreamServer) handleStreamMessage(userID string, envelope *pb.ChatStreamEnvelope, stream pb.ChatService_ChatStreamServer) {
	switch payload := envelope.Payload.(type) {
	case *pb.ChatStreamEnvelope_Message:
		// Handle chat message - this will route to SendMessage, EditMessage, etc.
		s.handleChatMessage(userID, payload.Message, stream)
	case *pb.ChatStreamEnvelope_Request:
		// Handle operation requests with explicit operation types
		s.handleStreamRequest(userID, payload.Request, stream)
	case *pb.ChatStreamEnvelope_State:
		// Handle state messages (typing indicators, presence, etc.)
		s.handleStateMessage(userID, payload.State, stream)
	default:
		s.sendStreamError(stream, 400, "Bad Request", "Unknown message type")
	}
}

// Add new function to handle stream requests with operation types
func (s *ChatStreamServer) handleStreamRequest(userID string, request *pb.StreamRequest, stream pb.ChatService_ChatStreamServer) {
	stream.Context()

	switch request.Operation {
	case pb.StreamOperation_SEND_MESSAGE:
		if req, ok := request.Request.(*pb.StreamRequest_SendMessage); ok {
			s.StreamSendMessage(userID, req.SendMessage, stream)
		} else {
			s.sendStreamError(stream, 400, "Bad Request", "Invalid SendMessage request")
		}

	case pb.StreamOperation_EDIT_MESSAGE:
		if req, ok := request.Request.(*pb.StreamRequest_EditMessage); ok {
			s.StreamEditMessage(userID, req.EditMessage, stream)
		} else {
			s.sendStreamError(stream, 400, "Bad Request", "Invalid EditMessage request")
		}

	case pb.StreamOperation_DELETE_MESSAGE:
		if req, ok := request.Request.(*pb.StreamRequest_DeleteMessage); ok {
			s.StreamDeleteMessage(userID, req.DeleteMessage, stream)
		} else {
			s.sendStreamError(stream, 400, "Bad Request", "Invalid DeleteMessage request")
		}

	case pb.StreamOperation_GET_CHAT_HISTORY:
		if req, ok := request.Request.(*pb.StreamRequest_GetChatHistory); ok {
			s.StreamGetChatHistory(userID, req.GetChatHistory, stream)
		} else {
			s.sendStreamError(stream, 400, "Bad Request", "Invalid GetChatHistory request")
		}

	case pb.StreamOperation_MARK_AS_READ:
		if req, ok := request.Request.(*pb.StreamRequest_MarkAsRead); ok {
			s.StreamMarkAsRead(userID, req.MarkAsRead, stream)
		} else {
			s.sendStreamError(stream, 400, "Bad Request", "Invalid MarkAsRead request")
		}

	case pb.StreamOperation_SEARCH_MESSAGES:
		if req, ok := request.Request.(*pb.StreamRequest_SearchMessages); ok {
			s.StreamSearchMessages(userID, req.SearchMessages, stream)
		} else {
			s.sendStreamError(stream, 400, "Bad Request", "Invalid SearchMessages request")
		}

	case pb.StreamOperation_FORWARD_MESSAGE:
		if req, ok := request.Request.(*pb.StreamRequest_ForwardMessage); ok {
			s.StreamForwardMessage(userID, req.ForwardMessage, stream)
		} else {
			s.sendStreamError(stream, 400, "Bad Request", "Invalid ForwardMessage request")
		}

	case pb.StreamOperation_PIN_MESSAGE:
		if req, ok := request.Request.(*pb.StreamRequest_PinMessage); ok {
			s.StreamPinMessage(userID, req.PinMessage, stream)
		} else {
			s.sendStreamError(stream, 400, "Bad Request", "Invalid PinMessage request")
		}

	case pb.StreamOperation_UNPIN_MESSAGE:
		if req, ok := request.Request.(*pb.StreamRequest_UnpinMessage); ok {
			s.StreamUnpinMessage(userID, req.UnpinMessage, stream)
		} else {
			s.sendStreamError(stream, 400, "Bad Request", "Invalid UnpinMessage request")
		}

	case pb.StreamOperation_GET_PINNED_MESSAGES:
		if req, ok := request.Request.(*pb.StreamRequest_GetPinnedMessages); ok {
			s.StreamGetPinnedMessages(userID, req.GetPinnedMessages, stream)
		} else {
			s.sendStreamError(stream, 400, "Bad Request", "Invalid GetPinnedMessages request")
		}

	default:
		s.sendStreamError(stream, 400, "Bad Request", "Unknown operation type")
	}
}

// handleChatMessage routes chat messages to appropriate handlers
func (s *ChatStreamServer) handleChatMessage(userID string, message *pb.ChatMessage, stream pb.ChatService_ChatStreamServer) {
	ctx := stream.Context()
	req := &pb.SendMessageRequest{
		SenderId:   message.SenderId,
		ReceiverId: message.ReceiverId,
		GroupId:    message.GroupId,
		Content:    message.Content,
		Type:       message.Type,
		IsGroup:    message.IsGroup,
	}
	// Call the existing SendMessage method
	response, err := s.ChatServer.SendMessage(ctx, req)
	if err != nil {
		s.sendStreamError(stream, 500, "Internal Server Error", err.Error())
		return
	}
	// Convert response to stream envelope
	var envelope *pb.ChatStreamEnvelope
	if response.StatusCode == 200 {
		if savedMsg, ok := response.Result.(*pb.SendMessageResponse_SavedMessage); ok {
			envelope = &pb.ChatStreamEnvelope{
				Payload: &pb.ChatStreamEnvelope_Message{
					Message: savedMsg.SavedMessage,
				},
			}
			// Broadcast to relevant users
			s.broadcastMessage(savedMsg.SavedMessage)
		}
	} else {
		if errorMsg, ok := response.Result.(*pb.SendMessageResponse_Error); ok {
			envelope = &pb.ChatStreamEnvelope{
				Payload: &pb.ChatStreamEnvelope_Error{
					Error: errorMsg.Error,
				},
			}
		}
	}

	// Send response back to client
	if envelope != nil {
		if err := stream.Send(envelope); err != nil {
			log.Printf("Failed to send stream response: %v", err)
		}
	}
}

// handleStateMessage processes state messages like typing indicators
func (s *ChatStreamServer) handleStateMessage(userID string, state *pb.StateMessage, stream pb.ChatService_ChatStreamServer) {
	// Handle typing indicators, presence updates, etc.
	log.Printf("State message from %s: %s", userID, state.Message)

	// Echo back the state message for now
	envelope := &pb.ChatStreamEnvelope{
		Payload: &pb.ChatStreamEnvelope_State{
			State: state,
		},
	}

	if err := stream.Send(envelope); err != nil {
		log.Printf("Failed to send state response: %v", err)
	}
}

// broadcastMessage sends a message to all relevant connected users
func (s *ChatStreamServer) broadcastMessage(message *pb.ChatMessage) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Determine who should receive this message
	var targetUserIDs []string
	if message.IsGroup {
		// For group messages, broadcast to all group members
		// You'd need to implement getGroupMembers method
		targetUserIDs = s.getGroupMembers(message.GroupId)
	} else {
		// For direct messages, send to receiver
		targetUserIDs = []string{message.ReceiverId}
	}

	// Send to all target users
	envelope := &pb.ChatStreamEnvelope{
		Payload: &pb.ChatStreamEnvelope_Message{
			Message: message,
		},
	}

	for _, targetUserID := range targetUserIDs {
		if connection, exists := s.connections[targetUserID]; exists {
			go func(conn *StreamConnection) {
				if err := conn.Stream.Send(envelope); err != nil {
					log.Printf("Failed to broadcast to user %s: %v", conn.UserID, err)
				}
			}(connection)
		}
	}
}

// registerConnection adds a new streaming connection
func (s *ChatStreamServer) registerConnection(userID string, connection *StreamConnection) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.connections[userID] = connection
	log.Printf("Registered streaming connection for user: %s", userID)
	// Broadcast that user is now online
	s.broadcastUserStatusUpdate(userID, true, timestamppb.New(time.Now()))
}

// unregisterConnection removes a streaming connection
func (s *ChatStreamServer) unregisterConnection(userID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.connections, userID)
	log.Printf("Unregistered streaming connection for user: %s", userID)
	// Broadcast that user is now offline
	s.broadcastUserStatusUpdate(userID, false, timestamppb.New(time.Now()))
}

// validateStreamRequest validates the streaming request
func (s *ChatStreamServer) validateStreamRequest(ctx context.Context) (*TokenClaims, error) {
	return s.ChatServer.validateRequest(ctx)
}

// sendStreamError sends an error message through the stream
func (s *ChatStreamServer) sendStreamError(stream pb.ChatService_ChatStreamServer, code int32, message, details string) error {
	errorEnvelope := &pb.ChatStreamEnvelope{
		Payload: &pb.ChatStreamEnvelope_Error{
			Error: &pb.ErrorMessage{
				Code:      code,
				Message:   message,
				Details:   []string{details},
				Timestamp: timestamppb.New(time.Now()),
			},
		},
	}
	return stream.Send(errorEnvelope)
}

// getGroupMembers returns the list of user IDs in a group
func (s *ChatStreamServer) getGroupMembers(groupID string) []string {
	// This should query the database for group members
	// For now, returning empty slice
	return []string{}
}

// StreamSendMessage handles SendMessage through streaming
func (s *ChatStreamServer) StreamSendMessage(userID string, req *pb.SendMessageRequest, stream pb.ChatService_ChatStreamServer) {
	ctx := stream.Context()
	response, err := s.ChatServer.SendMessage(ctx, req)
	if err != nil {
		s.sendStreamError(stream, 500, "Internal Server Error", err.Error())
		return
	}

	// Handle response and broadcast
	s.handleSendMessageResponse(response, stream)
}

// StreamEditMessage handles EditMessage through streaming
func (s *ChatStreamServer) StreamEditMessage(userID string, req *pb.EditMessageRequest, stream pb.ChatService_ChatStreamServer) {
	ctx := stream.Context()
	response, err := s.ChatServer.EditMessage(ctx, req)
	if err != nil {
		s.sendStreamError(stream, 500, "Internal Server Error", err.Error())
		return
	}

	// Handle response and broadcast
	s.handleEditMessageResponse(response, stream)
}

// StreamDeleteMessage handles DeleteMessage through streaming
func (s *ChatStreamServer) StreamDeleteMessage(userID string, req *pb.DeleteMessageRequest, stream pb.ChatService_ChatStreamServer) {
	ctx := stream.Context()
	response, err := s.ChatServer.DeleteMessage(ctx, req)
	if err != nil {
		s.sendStreamError(stream, 500, "Internal Server Error", err.Error())
		return
	}

	// Handle response and broadcast
	s.handleDeleteMessageResponse(response, stream)
}

func (s *ChatStreamServer) StreamGetChatHistory(userID string, req *pb.GetChatHistoryRequest, stream pb.ChatService_ChatStreamServer) {
	ctx := stream.Context()
	response, err := s.ChatServer.GetChatHistory(ctx, req)
	if err != nil {
		s.sendStreamError(stream, 500, "Internal Server Error", err.Error())
		return
	}
	s.handleGetChatHistoryResponse(response, stream)
}

func (s *ChatStreamServer) StreamMarkAsRead(userID string, req *pb.ReadReceiptRequest, stream pb.ChatService_ChatStreamServer) {
	ctx := stream.Context()
	response, err := s.ChatServer.MarkAsRead(ctx, req)
	if err != nil {
		s.sendStreamError(stream, 500, "Internal Server Error", err.Error())
		return
	}
	s.handleMarkAsReadResponse(response, stream)
}

func (s *ChatStreamServer) StreamSearchMessages(userID string, req *pb.SearchMessagesRequest, stream pb.ChatService_ChatStreamServer) {
	ctx := stream.Context()
	response, err := s.ChatServer.SearchMessages(ctx, req)
	if err != nil {
		s.sendStreamError(stream, 500, "Internal Server Error", err.Error())
		return
	}
	s.handleSearchMessagesResponse(response, stream)
}

func (s *ChatStreamServer) StreamForwardMessage(userID string, req *pb.ForwardMessageRequest, stream pb.ChatService_ChatStreamServer) {
	ctx := stream.Context()
	response, err := s.ChatServer.ForwardMessage(ctx, req)
	if err != nil {
		s.sendStreamError(stream, 500, "Internal Server Error", err.Error())
		return
	}
	s.handleForwardMessageResponse(response, stream)
}

func (s *ChatStreamServer) StreamPinMessage(userID string, req *pb.PinMessageRequest, stream pb.ChatService_ChatStreamServer) {
	ctx := stream.Context()
	response, err := s.ChatServer.PinMessage(ctx, req)
	if err != nil {
		s.sendStreamError(stream, 500, "Internal Server Error", err.Error())
		return
	}
	s.handlePinMessageResponse(response, stream)
}

func (s *ChatStreamServer) StreamUnpinMessage(userID string, req *pb.UnpinMessageRequest, stream pb.ChatService_ChatStreamServer) {
	ctx := stream.Context()
	response, err := s.ChatServer.UnpinMessage(ctx, req)
	if err != nil {
		s.sendStreamError(stream, 500, "Internal Server Error", err.Error())
		return
	}
	s.handleUnpinMessageResponse(response, stream)
}

func (s *ChatStreamServer) StreamGetPinnedMessages(userID string, req *pb.GetPinnedMessagesRequest, stream pb.ChatService_ChatStreamServer) {
	ctx := stream.Context()
	response, err := s.ChatServer.GetPinnedMessages(ctx, req)
	if err != nil {
		s.sendStreamError(stream, 500, "Internal Server Error", err.Error())
		return
	}
	s.handleGetPinnedMessagesResponse(response, stream)
}

func (s *ChatStreamServer) handleSendMessageResponse(response *pb.SendMessageResponse, stream pb.ChatService_ChatStreamServer) {
	var envelope *pb.ChatStreamEnvelope
	if response.StatusCode == 200 {
		if savedMsg, ok := response.Result.(*pb.SendMessageResponse_SavedMessage); ok {
			envelope = &pb.ChatStreamEnvelope{
				Payload: &pb.ChatStreamEnvelope_Message{
					Message: savedMsg.SavedMessage,
				},
			}
			s.broadcastMessage(savedMsg.SavedMessage)
			s.broadcastLastMessageUpdate(savedMsg.SavedMessage)
		}
	} else {
		if errorMsg, ok := response.Result.(*pb.SendMessageResponse_Error); ok {
			envelope = &pb.ChatStreamEnvelope{
				Payload: &pb.ChatStreamEnvelope_Error{
					Error: errorMsg.Error,
				},
			}
		}
	}

	if envelope != nil {
		stream.Send(envelope)
	}
}

func (s *ChatStreamServer) handleEditMessageResponse(response *pb.EditMessageResponse, stream pb.ChatService_ChatStreamServer) {
	var envelope *pb.ChatStreamEnvelope
	if response.StatusCode == 200 {
		if updatedMsg, ok := response.Result.(*pb.EditMessageResponse_UpdatedMessage); ok {
			envelope = &pb.ChatStreamEnvelope{
				Payload: &pb.ChatStreamEnvelope_Message{
					Message: updatedMsg.UpdatedMessage,
				},
			}
			s.broadcastMessage(updatedMsg.UpdatedMessage)
		}
	} else {
		if errorMsg, ok := response.Result.(*pb.EditMessageResponse_Error); ok {
			envelope = &pb.ChatStreamEnvelope{
				Payload: &pb.ChatStreamEnvelope_Error{
					Error: errorMsg.Error,
				},
			}
		}
	}

	if envelope != nil {
		stream.Send(envelope)
	}
}

func (s *ChatStreamServer) handleDeleteMessageResponse(response *pb.DeleteMessageResponse, stream pb.ChatService_ChatStreamServer) {
	var envelope *pb.ChatStreamEnvelope
	if response.StatusCode == 200 {
		if status, ok := response.Result.(*pb.DeleteMessageResponse_Status); ok {
			envelope = &pb.ChatStreamEnvelope{
				Payload: &pb.ChatStreamEnvelope_State{
					State: status.Status,
				},
			}
		}
	} else {
		if errorMsg, ok := response.Result.(*pb.DeleteMessageResponse_Error); ok {
			envelope = &pb.ChatStreamEnvelope{
				Payload: &pb.ChatStreamEnvelope_Error{
					Error: errorMsg.Error,
				},
			}
		}
	}

	if envelope != nil {
		stream.Send(envelope)
	}
}

func (s *ChatStreamServer) handleGetChatHistoryResponse(response *pb.GetChatHistoryResponse, stream pb.ChatService_ChatStreamServer) {
	if response.StatusCode == 200 {
		if messages, ok := response.Result.(*pb.GetChatHistoryResponse_Messages); ok {
			for _, message := range messages.Messages.Messages {
				envelope := &pb.ChatStreamEnvelope{
					Payload: &pb.ChatStreamEnvelope_Message{
						Message: message,
					},
				}
				stream.Send(envelope)
			}
		}
	} else {
		if errorMsg, ok := response.Result.(*pb.GetChatHistoryResponse_Error); ok {
			envelope := &pb.ChatStreamEnvelope{
				Payload: &pb.ChatStreamEnvelope_Error{
					Error: errorMsg.Error,
				},
			}
			stream.Send(envelope)
		}
	}
}

func (s *ChatStreamServer) handleMarkAsReadResponse(response *pb.ReadReceiptResponse, stream pb.ChatService_ChatStreamServer) {
	var envelope *pb.ChatStreamEnvelope
	if response.StatusCode == 200 {
		if status, ok := response.Result.(*pb.ReadReceiptResponse_Status); ok {
			envelope = &pb.ChatStreamEnvelope{
				Payload: &pb.ChatStreamEnvelope_State{
					State: status.Status,
				},
			}
		}
	} else {
		if errorMsg, ok := response.Result.(*pb.ReadReceiptResponse_Error); ok {
			envelope = &pb.ChatStreamEnvelope{
				Payload: &pb.ChatStreamEnvelope_Error{
					Error: errorMsg.Error,
				},
			}
		}
	}

	if envelope != nil {
		stream.Send(envelope)
	}
}

func (s *ChatStreamServer) handleSearchMessagesResponse(response *pb.SearchMessagesResponse, stream pb.ChatService_ChatStreamServer) {
	if response.StatusCode == 200 {
		if messages, ok := response.Result.(*pb.SearchMessagesResponse_Messages); ok {
			for _, message := range messages.Messages.Messages {
				envelope := &pb.ChatStreamEnvelope{
					Payload: &pb.ChatStreamEnvelope_Message{
						Message: message,
					},
				}
				stream.Send(envelope)
			}
		}
	} else {
		if errorMsg, ok := response.Result.(*pb.SearchMessagesResponse_Error); ok {
			envelope := &pb.ChatStreamEnvelope{
				Payload: &pb.ChatStreamEnvelope_Error{
					Error: errorMsg.Error,
				},
			}
			stream.Send(envelope)
		}
	}
}

func (s *ChatStreamServer) handleForwardMessageResponse(response *pb.ForwardMessageResponse, stream pb.ChatService_ChatStreamServer) {
	if response.StatusCode == 200 {
		if messages, ok := response.Result.(*pb.ForwardMessageResponse_ForwardedMessages); ok {
			for _, message := range messages.ForwardedMessages.Messages {
				s.broadcastMessage(message)
			}
			envelope := &pb.ChatStreamEnvelope{
				Payload: &pb.ChatStreamEnvelope_State{
					State: &pb.StateMessage{
						StatusCode: 200,
						Message:    "Messages forwarded successfully",
						Timestamp:  timestamppb.New(time.Now()),
					},
				},
			}
			stream.Send(envelope)
		}
	} else {
		if errorMsg, ok := response.Result.(*pb.ForwardMessageResponse_Error); ok {
			envelope := &pb.ChatStreamEnvelope{
				Payload: &pb.ChatStreamEnvelope_Error{
					Error: errorMsg.Error,
				},
			}
			stream.Send(envelope)
		}
	}
}

func (s *ChatStreamServer) handlePinMessageResponse(response *pb.PinMessageResponse, stream pb.ChatService_ChatStreamServer) {
	var envelope *pb.ChatStreamEnvelope
	if response.StatusCode == 200 {
		if status, ok := response.Result.(*pb.PinMessageResponse_Status); ok {
			envelope = &pb.ChatStreamEnvelope{
				Payload: &pb.ChatStreamEnvelope_State{
					State: status.Status,
				},
			}
		}
	} else {
		if errorMsg, ok := response.Result.(*pb.PinMessageResponse_Error); ok {
			envelope = &pb.ChatStreamEnvelope{
				Payload: &pb.ChatStreamEnvelope_Error{
					Error: errorMsg.Error,
				},
			}
		}
	}
	if envelope != nil {
		stream.Send(envelope)
	}
}

func (s *ChatStreamServer) handleUnpinMessageResponse(response *pb.UnpinMessageResponse, stream pb.ChatService_ChatStreamServer) {
	var envelope *pb.ChatStreamEnvelope
	if response.StatusCode == 200 {
		if status, ok := response.Result.(*pb.UnpinMessageResponse_Status); ok {
			envelope = &pb.ChatStreamEnvelope{
				Payload: &pb.ChatStreamEnvelope_State{
					State: status.Status,
				},
			}
		}
	} else {
		if errorMsg, ok := response.Result.(*pb.UnpinMessageResponse_Error); ok {
			envelope = &pb.ChatStreamEnvelope{
				Payload: &pb.ChatStreamEnvelope_Error{
					Error: errorMsg.Error,
				},
			}
		}
	}
	if envelope != nil {
		stream.Send(envelope)
	}
}

func (s *ChatStreamServer) handleGetPinnedMessagesResponse(response *pb.GetPinnedMessagesResponse, stream pb.ChatService_ChatStreamServer) {
	if response.StatusCode == 200 {
		if messages, ok := response.Result.(*pb.GetPinnedMessagesResponse_PinnedMessages); ok {
			for _, message := range messages.PinnedMessages.Messages {
				envelope := &pb.ChatStreamEnvelope{
					Payload: &pb.ChatStreamEnvelope_Message{
						Message: message,
					},
				}
				stream.Send(envelope)
			}
		}
	} else {
		if errorMsg, ok := response.Result.(*pb.GetPinnedMessagesResponse_Error); ok {
			envelope := &pb.ChatStreamEnvelope{
				Payload: &pb.ChatStreamEnvelope_Error{
					Error: errorMsg.Error,
				},
			}
			stream.Send(envelope)
		}
	}
}

func (s *ChatStreamServer) SubscribeToLastMessages(req *pb.LastMessageStreamRequest, stream pb.ChatService_SubscribeToLastMessagesServer) error {
	ctx := stream.Context()
	claims, err := s.validateStreamRequest(ctx)
	if err != nil {
		return err
	}
	userID := claims.UserID
	lastMessageChan := make(chan *pb.ChatMessage, 100)
	s.registerLastMessageSubscription(userID, lastMessageChan)
	defer s.unregisterLastMessageSubscription(userID)

	if err := s.sendInitialLastMessages(userID, stream); err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case message := <-lastMessageChan:
			if err := stream.Send(message); err != nil {
				log.Printf("Failed to send last message update: %v", err)
				return err
			}
		}
	}
}

func (s *ChatStreamServer) registerLastMessageSubscription(userID string, ch chan *pb.ChatMessage) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.lastMessageSubscriptions == nil {
		s.lastMessageSubscriptions = make(map[string]chan *pb.ChatMessage)
	}
	s.lastMessageSubscriptions[userID] = ch
}

func (s *ChatStreamServer) unregisterLastMessageSubscription(userID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if ch, exists := s.lastMessageSubscriptions[userID]; exists {
		close(ch)
		delete(s.lastMessageSubscriptions, userID)
	}
}

func (s *ChatStreamServer) sendInitialLastMessages(userID string, stream pb.ChatService_SubscribeToLastMessagesServer) error {
	req := &pb.GetLastMessagesRequest{
		UserId: userID,
	}
	response, err := s.ChatServer.GetLastMessages(stream.Context(), req)
	if err != nil {
		return err
	}
	if response.StatusCode == 200 {
		if messages, ok := response.Result.(*pb.GetLastMessagesResponse_LastMessages); ok {
			for _, message := range messages.LastMessages.Messages {
				if err := stream.Send(message); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (s *ChatStreamServer) broadcastLastMessageUpdate(message *pb.ChatMessage) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	var targetUserIDs []string
	if message.IsGroup {
		targetUserIDs = s.getGroupMembers(message.GroupId)
	} else {
		targetUserIDs = []string{message.SenderId, message.ReceiverId}
	}
	for _, userID := range targetUserIDs {
		if ch, exists := s.lastMessageSubscriptions[userID]; exists {
			select {
			case ch <- message:
			default:
				log.Printf("Last message channel full for user %s", userID)
			}
		}
	}
}

func (s *ChatStreamServer) SubscribeToUserStatus(req *pb.UserStatusSubscriptionRequest, stream pb.ChatService_SubscribeToUserStatusServer) error {
	ctx := stream.Context()
	claims, err := s.validateStreamRequest(ctx)
	if err != nil {
		return err
	}
	subscriberID := claims.UserID
	statusChan := make(chan *pb.UserStatus, 100)
	s.registerUserStatusSubscription(subscriberID, req.UserIds, statusChan)
	defer s.unregisterUserStatusSubscription(subscriberID)
	if err := s.sendInitialUserStatuses(req.UserIds, stream); err != nil {
		log.Printf("Failed to send initial user statuses: %v", err)
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case status := <-statusChan:
			if err := stream.Send(status); err != nil {
				log.Printf("Failed to send status update: %v", err)
				return err
			}
		}
	}
}

func (s *ChatStreamServer) registerUserStatusSubscription(subscriberID string, userIDs []string, ch chan *pb.UserStatus) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.userStatusSubscriptions == nil {
		s.userStatusSubscriptions = make(map[string]chan *pb.UserStatus)
	}
	if s.userStatusWatchers == nil {
		s.userStatusWatchers = make(map[string][]string)
	}
	s.userStatusSubscriptions[subscriberID] = ch
	for _, userID := range userIDs {
		s.userStatusWatchers[userID] = append(s.userStatusWatchers[userID], subscriberID)
	}
}
func (s *ChatStreamServer) unregisterUserStatusSubscription(subscriberID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if ch, exists := s.userStatusSubscriptions[subscriberID]; exists {
		close(ch)
		delete(s.userStatusSubscriptions, subscriberID)
	}
	for userID, subscribers := range s.userStatusWatchers {
		for i, sub := range subscribers {
			if sub == subscriberID {
				s.userStatusWatchers[userID] = append(subscribers[:i], subscribers[i+1:]...)
				break
			}
		}
		if len(s.userStatusWatchers[userID]) == 0 {
			delete(s.userStatusWatchers, userID)
		}
	}
}

func (s *ChatStreamServer) sendInitialUserStatuses(userIDs []string, stream pb.ChatService_SubscribeToUserStatusServer) error {
	req := &pb.UserStatusRequest{
		UserIds: userIDs,
	}
	response, err := s.ChatServer.GetUserStatus(stream.Context(), req)
	if err != nil {
		return err
	}
	if response.StatusCode == 200 {
		if statuses, ok := response.Result.(*pb.UserStatusResponse_Statuses); ok {
			for _, status := range statuses.Statuses.Statuses {
				if err := stream.Send(status); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (s *ChatStreamServer) broadcastUserStatusUpdate(userID string, isOnline bool, lastSeen *timestamppb.Timestamp) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	status := &pb.UserStatus{
		UserId:   userID,
		IsOnline: isOnline,
		LastSeen: lastSeen,
	}
	if subscribers, exists := s.userStatusWatchers[userID]; exists {
		for _, subscriberID := range subscribers {
			if ch, exists := s.userStatusSubscriptions[subscriberID]; exists {
				select {
				case ch <- status:
				default:
					log.Printf("Status channel full for subscriber %s", subscriberID)
				}
			}
		}
	}
}

func (s *ChatStreamServer) SubscribeToNotifications(req *pb.SubscribeToNotificationsRequest, stream pb.ChatService_SubscribeToNotificationsServer) error {
	ctx := stream.Context()
	claims, err := s.validateStreamRequest(ctx)
	if err != nil {
		return err
	}
	userID := claims.UserID
	notificationChan := make(chan *pb.NotificationUpdate, 100)
	s.registerNotificationSubscription(userID, notificationChan)
	defer s.unregisterNotificationSubscription(userID)
	go s.sendInitialNotifications(userID, stream)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case notification := <-notificationChan:
			if err := stream.Send(notification); err != nil {
				log.Printf("Failed to send notification to user %s: %v", userID, err)
				return err
			}
		}
	}
}

func (s *ChatStreamServer) VideoCallStream(stream pb.ChatService_VideoCallStreamServer) error {
	ctx := stream.Context()
	claims, err := s.validateStreamRequest(ctx)
	if err != nil {
		return err
	}
	userID := claims.UserID
	callCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	s.registerVideoCallConnection(userID, stream)
	defer s.unregisterVideoCallConnection(userID)
	for {
		select {
		case <-callCtx.Done():
			return callCtx.Err()
		default:
			envelope, err := stream.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				log.Printf("Video call stream receive error: %v", err)
				return err
			}

			go s.handleVideoCallMessage(userID, envelope, stream)
		}
	}
}

func (s *ChatStreamServer) PhoneCallStream(stream pb.ChatService_PhoneCallStreamServer) error {
	ctx := stream.Context()
	claims, err := s.validateStreamRequest(ctx)
	if err != nil {
		return err
	}
	userID := claims.UserID
	callCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	s.registerPhoneCallConnection(userID, stream)
	defer s.unregisterPhoneCallConnection(userID)
	for {
		select {
		case <-callCtx.Done():
			return callCtx.Err()
		default:
			envelope, err := stream.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				log.Printf("Phone call stream receive error: %v", err)
				return err
			}
			go s.handlePhoneCallMessage(userID, envelope, stream)
		}
	}
}

func (s *ChatStreamServer) SubscribeToCallUpdates(req *pb.CallSubscriptionRequest, stream pb.ChatService_SubscribeToCallUpdatesServer) error {
	ctx := stream.Context()
	claims, err := s.validateStreamRequest(ctx)
	if err != nil {
		return err
	}
	userID := claims.UserID
	callUpdatesChan := make(chan *pb.CallUpdate, 100)
	s.registerCallUpdatesSubscription(userID, req.CallIds, callUpdatesChan)
	defer s.unregisterCallUpdatesSubscription(userID)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case update := <-callUpdatesChan:
			if err := stream.Send(update); err != nil {
				log.Printf("Failed to send call update to user %s: %v", userID, err)
				return err
			}
		}
	}
}

func (s *ChatStreamServer) SubscribeToCallSignaling(req *pb.CallSignalingRequest, stream pb.ChatService_SubscribeToCallSignalingServer) error {
	ctx := stream.Context()
	claims, err := s.validateStreamRequest(ctx)
	if err != nil {
		return err
	}
	userID := claims.UserID
	log.Printf("User %s subscribed to call signaling for call %s", userID, req.CallId)

	signalingChan := make(chan *pb.CallSignalingMessage, 100)
	s.registerCallSignalingSubscription(userID, req.CallId, signalingChan)
	defer s.unregisterCallSignalingSubscription(userID, req.CallId)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case signalingMsg := <-signalingChan:
			if err := stream.Send(signalingMsg); err != nil {
				log.Printf("Failed to send signaling message to user %s: %v", userID, err)
				return err
			}
		}
	}
}
