package server

import (
	"context"
	"database/sql"
	"log"
	"time"

	pb "user-progress-services/proto"
)

type UserProgressServer struct {
	pb.UnimplementedUserProgressServiceServer
	db *Database
}

func NewUserProgressServer() *UserProgressServer {
	return &UserProgressServer{
		db: NewDatabase(),
	}
}

// GetUserProgress retrieves user progress for templates
func (s *UserProgressServer) GetUserProgress(ctx context.Context, req *pb.GetUserProgressRequest) (*pb.GetUserProgressResponse, error) {
	log.Printf("GetUserProgress called for user: %s, template: %s", req.UserId, req.TemplateId)

	if req.UserId == "" {
		return &pb.GetUserProgressResponse{
			StatusCode: 400,
			Message:    "User ID is required",
			Result: &pb.GetUserProgressResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "User ID cannot be empty",
					Details:   []string{"Please provide a valid user ID"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	progressList, err := s.db.GetUserProgress(req.UserId, req.TemplateId)
	if err != nil {
		return &pb.GetUserProgressResponse{
			StatusCode: 500,
			Message:    "Failed to retrieve user progress",
			Result: &pb.GetUserProgressResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Internal server error",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.GetUserProgressResponse{
		StatusCode: 200,
		Message:    "User progress retrieved successfully",
		Result: &pb.GetUserProgressResponse_Progress{
			Progress: &pb.UserProgressList{
				Progress:   progressList,
				TotalCount: int32(len(progressList)),
			},
		},
	}, nil
}

// UpdateProgress updates or creates user progress
func (s *UserProgressServer) UpdateProgress(ctx context.Context, req *pb.UpdateProgressRequest) (*pb.UpdateProgressResponse, error) {
	log.Printf("UpdateProgress called for user: %s, template: %s, section: %s", req.UserId, req.TemplateId, req.SectionId)

	if req.UserId == "" {
		return &pb.UpdateProgressResponse{
			StatusCode: 400,
			Message:    "User ID is required",
			Result: &pb.UpdateProgressResponse_Error{
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
		return &pb.UpdateProgressResponse{
			StatusCode: 400,
			Message:    "Template ID is required",
			Result: &pb.UpdateProgressResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Template ID cannot be empty",
					Details:   []string{"Please provide a valid template ID"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	if req.SectionId == "" {
		return &pb.UpdateProgressResponse{
			StatusCode: 400,
			Message:    "Section ID is required",
			Result: &pb.UpdateProgressResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Section ID cannot be empty",
					Details:   []string{"Please provide a valid section ID"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	if req.Progress < 0 || req.Progress > 100 {
		return &pb.UpdateProgressResponse{
			StatusCode: 400,
			Message:    "Progress must be between 0 and 100",
			Result: &pb.UpdateProgressResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Invalid progress value",
					Details:   []string{"Progress must be between 0 and 100"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	progress, err := s.db.UpdateProgress(req.UserId, req.TemplateId, req.SectionId, req.Progress, req.Completed)
	if err != nil {
		return &pb.UpdateProgressResponse{
			StatusCode: 500,
			Message:    "Failed to update progress",
			Result: &pb.UpdateProgressResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Internal server error",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.UpdateProgressResponse{
		StatusCode: 200,
		Message:    "Progress updated successfully",
		Result: &pb.UpdateProgressResponse_Progress{
			Progress: progress,
		},
	}, nil
}

// GetCompletedTemplates retrieves completed templates for a user
func (s *UserProgressServer) GetCompletedTemplates(ctx context.Context, req *pb.GetCompletedTemplatesRequest) (*pb.GetCompletedTemplatesResponse, error) {
	log.Printf("GetCompletedTemplates called for user: %s", req.UserId)

	if req.UserId == "" {
		return &pb.GetCompletedTemplatesResponse{
			StatusCode: 400,
			Message:    "User ID is required",
			Result: &pb.GetCompletedTemplatesResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "User ID cannot be empty",
					Details:   []string{"Please provide a valid user ID"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	page := req.Page
	limit := req.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	completedList, totalCount, err := s.db.GetCompletedTemplates(req.UserId, page, limit)
	if err != nil {
		return &pb.GetCompletedTemplatesResponse{
			StatusCode: 500,
			Message:    "Failed to retrieve completed templates",
			Result: &pb.GetCompletedTemplatesResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Internal server error",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.GetCompletedTemplatesResponse{
		StatusCode: 200,
		Message:    "Completed templates retrieved successfully",
		Result: &pb.GetCompletedTemplatesResponse_Completed{
			Completed: &pb.UserProgressList{
				Progress:   completedList,
				TotalCount: totalCount,
			},
		},
	}, nil
}

// GetUserAchievements retrieves user achievements
func (s *UserProgressServer) GetUserAchievements(ctx context.Context, req *pb.GetUserAchievementsRequest) (*pb.GetUserAchievementsResponse, error) {
	log.Printf("GetUserAchievements called for user: %s", req.UserId)

	if req.UserId == "" {
		return &pb.GetUserAchievementsResponse{
			StatusCode: 400,
			Message:    "User ID is required",
			Result: &pb.GetUserAchievementsResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "User ID cannot be empty",
					Details:   []string{"Please provide a valid user ID"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	achievements, err := s.db.GetUserAchievements(req.UserId)
	if err != nil {
		return &pb.GetUserAchievementsResponse{
			StatusCode: 500,
			Message:    "Failed to retrieve user achievements",
			Result: &pb.GetUserAchievementsResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Internal server error",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.GetUserAchievementsResponse{
		StatusCode: 200,
		Message:    "User achievements retrieved successfully",
		Result: &pb.GetUserAchievementsResponse_Achievements{
			Achievements: &pb.UserAchievementList{
				Achievements: achievements,
				TotalCount:   int32(len(achievements)),
			},
		},
	}, nil
}

// UnlockAchievement unlocks an achievement for a user
func (s *UserProgressServer) UnlockAchievement(ctx context.Context, req *pb.UnlockAchievementRequest) (*pb.UnlockAchievementResponse, error) {
	log.Printf("UnlockAchievement called for user: %s, achievement: %s", req.UserId, req.AchievementId)

	if req.UserId == "" {
		return &pb.UnlockAchievementResponse{
			StatusCode: 400,
			Message:    "User ID is required",
			Result: &pb.UnlockAchievementResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "User ID cannot be empty",
					Details:   []string{"Please provide a valid user ID"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	if req.AchievementId == "" {
		return &pb.UnlockAchievementResponse{
			StatusCode: 400,
			Message:    "Achievement ID is required",
			Result: &pb.UnlockAchievementResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Achievement ID cannot be empty",
					Details:   []string{"Please provide a valid achievement ID"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	achievement, err := s.db.UnlockAchievement(req.UserId, req.AchievementId)
	if err != nil {
		statusCode := 500
		message := "Failed to unlock achievement"
		if err.Error() == "achievement already unlocked" {
			statusCode = 409
			message = "Achievement already unlocked"
		} else if err == sql.ErrNoRows {
			statusCode = 404
			message = "Achievement not found"
		}

		return &pb.UnlockAchievementResponse{
			StatusCode: int32(statusCode),
			Message:    message,
			Result: &pb.UnlockAchievementResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      int32(statusCode),
					Message:   message,
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.UnlockAchievementResponse{
		StatusCode: 200,
		Message:    "Achievement unlocked successfully",
		Result: &pb.UnlockAchievementResponse_Achievement{
			Achievement: achievement,
		},
	}, nil
}

// GetLearningPath generates a learning path for a user
func (s *UserProgressServer) GetLearningPath(ctx context.Context, req *pb.GetLearningPathRequest) (*pb.GetLearningPathResponse, error) {
	log.Printf("GetLearningPath called for user: %s", req.UserId)

	if req.UserId == "" {
		return &pb.GetLearningPathResponse{
			StatusCode: 400,
			Message:    "User ID is required",
			Result: &pb.GetLearningPathResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "User ID cannot be empty",
					Details:   []string{"Please provide a valid user ID"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	learningPath, err := s.db.GetLearningPath(req.UserId)
	if err != nil {
		return &pb.GetLearningPathResponse{
			StatusCode: 500,
			Message:    "Failed to generate learning path",
			Result: &pb.GetLearningPathResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Internal server error",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.GetLearningPathResponse{
		StatusCode: 200,
		Message:    "Learning path generated successfully",
		Result: &pb.GetLearningPathResponse_Path{
			Path: learningPath,
		},
	}, nil
}