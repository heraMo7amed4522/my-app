package server

import (
	"context"
	"database/sql"
	"log"
	"time"

	pb "feedback-services/proto"
)

type FeedbackServer struct {
	pb.UnimplementedFeedbackServiceServer
	db *Database
}

func NewFeedbackServer() *FeedbackServer {
	return &FeedbackServer{
		db: NewDatabase(),
	}
}

// SubmitFeedback handles feedback submission
func (s *FeedbackServer) SubmitFeedback(ctx context.Context, req *pb.SubmitFeedbackRequest) (*pb.SubmitFeedbackResponse, error) {
	log.Printf("SubmitFeedback called for user: %s, template: %s", req.UserId, req.TemplateId)
	if req.UserId == "" {
		return &pb.SubmitFeedbackResponse{
			StatusCode: 400,
			Message:    "User ID is required",
			Result: &pb.SubmitFeedbackResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "User ID cannot be empty",
					Details:   []string{"Please provide a valid user ID"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	if req.TemplateId == "" {
		return &pb.SubmitFeedbackResponse{
			StatusCode: 400,
			Message:    "Template ID is required",
			Result: &pb.SubmitFeedbackResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Template ID cannot be empty",
					Details:   []string{"Please provide a valid template ID"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	if req.Rating < 1 || req.Rating > 5 {
		return &pb.SubmitFeedbackResponse{
			StatusCode: 400,
			Message:    "Rating must be between 1 and 5",
			Result: &pb.SubmitFeedbackResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Invalid rating value",
					Details:   []string{"Rating must be between 1.0 and 5.0"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	language := req.Language
	if language == "" {
		language = "en"
	}

	feedback, err := s.db.SubmitFeedback(req.UserId, req.TemplateId, req.Rating, req.Comment, language)
	if err != nil {
		log.Printf("Error submitting feedback: %v", err)
		return &pb.SubmitFeedbackResponse{
			StatusCode: 500,
			Message:    "Failed to submit feedback",
			Result: &pb.SubmitFeedbackResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Internal server error",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	return &pb.SubmitFeedbackResponse{
		StatusCode: 201,
		Message:    "Feedback submitted successfully",
		Result: &pb.SubmitFeedbackResponse_Feedback{
			Feedback: feedback,
		},
	}, nil
}

// GetTemplateFeedback retrieves feedback for a template
func (s *FeedbackServer) GetTemplateFeedback(ctx context.Context, req *pb.GetTemplateFeedbackRequest) (*pb.GetTemplateFeedbackResponse, error) {
	log.Printf("GetTemplateFeedback called for template: %s", req.TemplateId)

	if req.TemplateId == "" {
		return &pb.GetTemplateFeedbackResponse{
			StatusCode: 400,
			Message:    "Template ID is required",
			Result: &pb.GetTemplateFeedbackResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Template ID cannot be empty",
					Details:   []string{"Please provide a valid template ID"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	page := req.Page
	if page <= 0 {
		page = 1
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	feedbacks, totalCount, err := s.db.GetTemplateFeedback(req.TemplateId, page, limit)
	if err != nil {
		log.Printf("Error getting template feedback: %v", err)
		return &pb.GetTemplateFeedbackResponse{
			StatusCode: 500,
			Message:    "Failed to retrieve feedback",
			Result: &pb.GetTemplateFeedbackResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Internal server error",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.GetTemplateFeedbackResponse{
		StatusCode: 200,
		Message:    "Feedback retrieved successfully",
		Result: &pb.GetTemplateFeedbackResponse_Feedback{
			Feedback: &pb.FeedbackList{
				Feedback:   feedbacks,
				TotalCount: totalCount,
				Page:       page,
				Limit:      limit,
			},
		},
	}, nil
}

// GetFeedbackStatistics calculates feedback statistics
func (s *FeedbackServer) GetFeedbackStatistics(ctx context.Context, req *pb.GetFeedbackStatisticsRequest) (*pb.GetFeedbackStatisticsResponse, error) {
	log.Printf("GetFeedbackStatistics called for template: %s", req.TemplateId)
	if req.TemplateId == "" {
		return &pb.GetFeedbackStatisticsResponse{
			StatusCode: 400,
			Message:    "Template ID is required",
			Result: &pb.GetFeedbackStatisticsResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Template ID cannot be empty",
					Details:   []string{"Please provide a valid template ID"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	statistics, err := s.db.GetFeedbackStatistics(req.TemplateId)
	if err != nil {
		log.Printf("Error getting feedback statistics: %v", err)
		return &pb.GetFeedbackStatisticsResponse{
			StatusCode: 500,
			Message:    "Failed to retrieve statistics",
			Result: &pb.GetFeedbackStatisticsResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Internal server error",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.GetFeedbackStatisticsResponse{
		StatusCode: 200,
		Message:    "Statistics retrieved successfully",
		Result: &pb.GetFeedbackStatisticsResponse_Statistics{
			Statistics: statistics,
		},
	}, nil
}

// UpdateFeedback updates existing feedback
func (s *FeedbackServer) UpdateFeedback(ctx context.Context, req *pb.UpdateFeedbackRequest) (*pb.UpdateFeedbackResponse, error) {
	log.Printf("UpdateFeedback called for feedback ID: %s", req.Id)

	if req.Id == "" {
		return &pb.UpdateFeedbackResponse{
			StatusCode: 400,
			Message:    "Feedback ID is required",
			Result: &pb.UpdateFeedbackResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Feedback ID cannot be empty",
					Details:   []string{"Please provide a valid feedback ID"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	if req.Rating < 1 || req.Rating > 5 {
		return &pb.UpdateFeedbackResponse{
			StatusCode: 400,
			Message:    "Rating must be between 1 and 5",
			Result: &pb.UpdateFeedbackResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Invalid rating value",
					Details:   []string{"Rating must be between 1.0 and 5.0"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	feedback, err := s.db.UpdateFeedback(req.Id, req.Rating, req.Comment)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.UpdateFeedbackResponse{
				StatusCode: 404,
				Message:    "Feedback not found",
				Result: &pb.UpdateFeedbackResponse_Error{
					Error: &pb.ErrorDetails{
						Code:      404,
						Message:   "Feedback with the specified ID does not exist",
						Details:   []string{"Please check the feedback ID"},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			}, nil
		}

		log.Printf("Error updating feedback: %v", err)
		return &pb.UpdateFeedbackResponse{
			StatusCode: 500,
			Message:    "Failed to update feedback",
			Result: &pb.UpdateFeedbackResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Internal server error",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.UpdateFeedbackResponse{
		StatusCode: 200,
		Message:    "Feedback updated successfully",
		Result: &pb.UpdateFeedbackResponse_Feedback{
			Feedback: feedback,
		},
	}, nil
}

// DeleteFeedback removes feedback
func (s *FeedbackServer) DeleteFeedback(ctx context.Context, req *pb.DeleteFeedbackRequest) (*pb.DeleteFeedbackResponse, error) {
	log.Printf("DeleteFeedback called for feedback ID: %s", req.Id)

	if req.Id == "" {
		return &pb.DeleteFeedbackResponse{
			StatusCode: 400,
			Message:    "Feedback ID is required",
			Result: &pb.DeleteFeedbackResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Feedback ID cannot be empty",
					Details:   []string{"Please provide a valid feedback ID"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	err := s.db.DeleteFeedback(req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.DeleteFeedbackResponse{
				StatusCode: 404,
				Message:    "Feedback not found",
				Result: &pb.DeleteFeedbackResponse_Error{
					Error: &pb.ErrorDetails{
						Code:      404,
						Message:   "Feedback with the specified ID does not exist",
						Details:   []string{"Please check the feedback ID"},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			}, nil
		}

		log.Printf("Error deleting feedback: %v", err)
		return &pb.DeleteFeedbackResponse{
			StatusCode: 500,
			Message:    "Failed to delete feedback",
			Result: &pb.DeleteFeedbackResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Internal server error",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.DeleteFeedbackResponse{
		StatusCode: 200,
		Message:    "Feedback deleted successfully",
		Result: &pb.DeleteFeedbackResponse_Success{
			Success: true,
		},
	}, nil
}
