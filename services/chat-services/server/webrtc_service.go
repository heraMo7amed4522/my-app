package server

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	pb "chat-services/proto"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// WebRTCService handles WebRTC-specific operations
type WebRTCService struct {
	peerConnections   map[string]*PeerConnection
	signalingChannels map[string]chan *pb.CallSignalingMessage
	mutex             sync.RWMutex
}

// PeerConnection represents a WebRTC peer connection
type PeerConnection struct {
	CallID          string
	LocalUserID     string
	RemoteUserID    string
	ConnectionState string
	CreatedAt       time.Time
	LastActivity    time.Time
}

// WebRTCOffer represents a WebRTC offer
type WebRTCOffer struct {
	Type string `json:"type"`
	SDP  string `json:"sdp"`
}

// WebRTCAnswer represents a WebRTC answer
type WebRTCAnswer struct {
	Type string `json:"type"`
	SDP  string `json:"sdp"`
}

// ICECandidate represents an ICE candidate
type ICECandidate struct {
	Candidate     string `json:"candidate"`
	SDPMid        string `json:"sdpMid"`
	SDPMLineIndex int    `json:"sdpMLineIndex"`
}

// NewWebRTCService creates a new WebRTC service
func NewWebRTCService() *WebRTCService {
	return &WebRTCService{
		peerConnections:   make(map[string]*PeerConnection),
		signalingChannels: make(map[string]chan *pb.CallSignalingMessage),
	}
}

// CreatePeerConnection creates a new peer connection
func (w *WebRTCService) CreatePeerConnection(callID, localUserID, remoteUserID string) *PeerConnection {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	pc := &PeerConnection{
		CallID:          callID,
		LocalUserID:     localUserID,
		RemoteUserID:    remoteUserID,
		ConnectionState: "new",
		CreatedAt:       time.Now(),
		LastActivity:    time.Now(),
	}

	w.peerConnections[callID] = pc
	return pc
}

// ProcessSignalingMessage processes WebRTC signaling messages
func (w *WebRTCService) ProcessSignalingMessage(msg *pb.CallSignalingMessage) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	log.Printf("Processing signaling message for call %s: %s", msg.CallId, msg.SignalingType)

	switch msg.SignalingType {
	case pb.SignalingType_SIGNALING_OFFER:
		return w.handleOffer(msg)
	case pb.SignalingType_SIGNALING_ANSWER:
		return w.handleAnswer(msg)
	case pb.SignalingType_SIGNALING_ICE_CANDIDATE:
		return w.handleICECandidate(msg)
	case pb.SignalingType_SIGNALING_ICE_GATHERING_COMPLETE:
		return w.handleICEGatheringComplete(msg)
	case pb.SignalingType_SIGNALING_RENEGOTIATION:
		return w.handleRenegotiation(msg)
	default:
		log.Printf("Unknown signaling type: %s", msg.SignalingType)
		return nil
	}
}

// handleOffer processes WebRTC offer
func (w *WebRTCService) handleOffer(msg *pb.CallSignalingMessage) error {
	var offer WebRTCOffer
	if err := json.Unmarshal([]byte(msg.Payload), &offer); err != nil {
		return err
	}

	log.Printf("Received WebRTC offer for call %s", msg.CallId)

	// Update peer connection state
	if pc, exists := w.peerConnections[msg.CallId]; exists {
		pc.ConnectionState = "have-remote-offer"
		pc.LastActivity = time.Now()
	}

	return nil
}

// handleAnswer processes WebRTC answer
func (w *WebRTCService) handleAnswer(msg *pb.CallSignalingMessage) error {
	var answer WebRTCAnswer
	if err := json.Unmarshal([]byte(msg.Payload), &answer); err != nil {
		return err
	}

	log.Printf("Received WebRTC answer for call %s", msg.CallId)

	// Update peer connection state
	if pc, exists := w.peerConnections[msg.CallId]; exists {
		pc.ConnectionState = "stable"
		pc.LastActivity = time.Now()
	}

	return nil
}

// handleICECandidate processes ICE candidates
func (w *WebRTCService) handleICECandidate(msg *pb.CallSignalingMessage) error {
	var candidate ICECandidate
	if err := json.Unmarshal([]byte(msg.Payload), &candidate); err != nil {
		return err
	}

	log.Printf("Received ICE candidate for call %s", msg.CallId)

	// Update peer connection activity
	if pc, exists := w.peerConnections[msg.CallId]; exists {
		pc.LastActivity = time.Now()
	}

	return nil
}

// handleICEGatheringComplete processes ICE gathering completion
func (w *WebRTCService) handleICEGatheringComplete(msg *pb.CallSignalingMessage) error {
	log.Printf("ICE gathering complete for call %s", msg.CallId)

	// Update peer connection state
	if pc, exists := w.peerConnections[msg.CallId]; exists {
		pc.ConnectionState = "connected"
		pc.LastActivity = time.Now()
	}

	return nil
}

// handleRenegotiation processes renegotiation requests
func (w *WebRTCService) handleRenegotiation(msg *pb.CallSignalingMessage) error {
	log.Printf("Renegotiation requested for call %s", msg.CallId)

	// Update peer connection activity
	if pc, exists := w.peerConnections[msg.CallId]; exists {
		pc.LastActivity = time.Now()
	}

	return nil
}

// CreateOffer creates a WebRTC offer
func (w *WebRTCService) CreateOffer(callID string) (*pb.CallSignalingMessage, error) {
	offer := WebRTCOffer{
		Type: "offer",
		SDP:  "v=0\r\no=- 123456789 2 IN IP4 127.0.0.1\r\n...", // Simplified SDP
	}

	payload, err := json.Marshal(offer)
	if err != nil {
		return nil, err
	}

	return &pb.CallSignalingMessage{
		CallId:        callID,
		SignalingType: pb.SignalingType_SIGNALING_OFFER,
		Payload:       string(payload),
		Timestamp:     timestamppb.New(time.Now()),
	}, nil
}

// CreateAnswer creates a WebRTC answer
func (w *WebRTCService) CreateAnswer(callID string) (*pb.CallSignalingMessage, error) {
	answer := WebRTCAnswer{
		Type: "answer",
		SDP:  "v=0\r\no=- 987654321 2 IN IP4 127.0.0.1\r\n...", // Simplified SDP
	}

	payload, err := json.Marshal(answer)
	if err != nil {
		return nil, err
	}

	return &pb.CallSignalingMessage{
		CallId:        callID,
		SignalingType: pb.SignalingType_SIGNALING_ANSWER,
		Payload:       string(payload),
		Timestamp:     timestamppb.New(time.Now()),
	}, nil
}

// GetPeerConnection retrieves a peer connection
func (w *WebRTCService) GetPeerConnection(callID string) (*PeerConnection, bool) {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	pc, exists := w.peerConnections[callID]
	return pc, exists
}

// RemovePeerConnection removes a peer connection
func (w *WebRTCService) RemovePeerConnection(callID string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	delete(w.peerConnections, callID)
}

// CleanupInactivePeerConnections removes inactive peer connections
func (w *WebRTCService) CleanupInactivePeerConnections(timeout time.Duration) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	now := time.Now()
	for callID, pc := range w.peerConnections {
		if now.Sub(pc.LastActivity) > timeout {
			log.Printf("Cleaning up inactive peer connection for call %s", callID)
			delete(w.peerConnections, callID)
		}
	}
}
