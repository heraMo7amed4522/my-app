package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	pb "quiz-services/proto"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type QuizServer struct {
	pb.UnimplementedQuizServiceServer
	db *Database
}

func NewQuizServer() *QuizServer {
	return &QuizServer{
		db: NewDatabase(),
	}
}

// GetQuizzesBySection retrieves quizzes for a specific section
func (s *QuizServer) GetQuizzesBySection(ctx context.Context, req *pb.GetQuizzesBySectionRequest) (*pb.GetQuizzesBySectionResponse, error) {
	log.Printf("GetQuizzesBySection called with section ID: %s", req.SectionId)

	if req.SectionId == "" {
		return &pb.GetQuizzesBySectionResponse{
			StatusCode: 400,
			Message:    "Section ID is required",
			Result: &pb.GetQuizzesBySectionResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Section ID cannot be empty",
					Details:   []string{"Please provide a valid section ID"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	quizzes, err := s.db.GetQuizzesBySection(req.SectionId)
	if err != nil {
		return &pb.GetQuizzesBySectionResponse{
			StatusCode: 500,
			Message:    "Failed to retrieve quizzes",
			Result: &pb.GetQuizzesBySectionResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Internal server error",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.GetQuizzesBySectionResponse{
		StatusCode: 200,
		Message:    "Quizzes retrieved successfully",
		Result: &pb.GetQuizzesBySectionResponse_Quizzes{
			Quizzes: &pb.QuizList{
				Quizzes: quizzes,
			},
		},
	}, nil
}

// GetQuizByID retrieves a quiz by its ID
func (s *QuizServer) GetQuizByID(ctx context.Context, req *pb.GetQuizByIDRequest) (*pb.GetQuizByIDResponse, error) {
	log.Printf("GetQuizByID called with ID: %s", req.Id)

	if req.Id == "" {
		return &pb.GetQuizByIDResponse{
			StatusCode: 400,
			Message:    "Quiz ID is required",
			Result: &pb.GetQuizByIDResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Quiz ID cannot be empty",
					Details:   []string{"Please provide a valid quiz ID"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	quiz, err := s.db.GetQuizByID(req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.GetQuizByIDResponse{
				StatusCode: 404,
				Message:    "Quiz not found",
				Result: &pb.GetQuizByIDResponse_Error{
					Error: &pb.ErrorDetails{
						Code:      404,
						Message:   "Quiz with the specified ID does not exist",
						Details:   []string{fmt.Sprintf("Quiz ID: %s", req.Id)},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			}, nil
		}

		return &pb.GetQuizByIDResponse{
			StatusCode: 500,
			Message:    "Failed to retrieve quiz",
			Result: &pb.GetQuizByIDResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Internal server error",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.GetQuizByIDResponse{
		StatusCode: 200,
		Message:    "Quiz retrieved successfully",
		Result: &pb.GetQuizByIDResponse_Quiz{
			Quiz: quiz,
		},
	}, nil
}

// CreateQuiz creates a new quiz
func (s *QuizServer) CreateQuiz(ctx context.Context, req *pb.CreateQuizRequest) (*pb.CreateQuizResponse, error) {
	log.Printf("CreateQuiz called for section: %s", req.Quiz.SectionId)

	if req.Quiz.SectionId == "" || req.Quiz.Question == "" {
		return &pb.CreateQuizResponse{
			StatusCode: 400,
			Message:    "Section ID and question are required",
			Result: &pb.CreateQuizResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Missing required fields",
					Details:   []string{"Section ID and question cannot be empty"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	quiz := &pb.Quiz{
		Id:            uuid.New().String(),
		SectionId:     req.Quiz.SectionId,
		Question:      req.Quiz.Question,
		Options:       req.Quiz.Options,
		CorrectAnswer: req.Quiz.CorrectAnswer,
		Explanation:   req.Quiz.Explanation,
		Difficulty:    req.Quiz.Difficulty,
		CreatedAt:     timestamppb.Now(),
	}

	createdQuiz, err := s.db.CreateQuiz(quiz)
	if err != nil {
		return &pb.CreateQuizResponse{
			StatusCode: 500,
			Message:    "Failed to create quiz",
			Result: &pb.CreateQuizResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Internal server error",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.CreateQuizResponse{
		StatusCode: 201,
		Message:    "Quiz created successfully",
		Result: &pb.CreateQuizResponse_Quiz{
			Quiz: createdQuiz,
		},
	}, nil
}

// UpdateQuiz updates an existing quiz
func (s *QuizServer) UpdateQuiz(ctx context.Context, req *pb.UpdateQuizRequest) (*pb.UpdateQuizResponse, error) {
	log.Printf("UpdateQuiz called with ID: %s", req.Id)

	if req.Id == "" {
		return &pb.UpdateQuizResponse{
			StatusCode: 400,
			Message:    "Quiz ID is required",
			Result: &pb.UpdateQuizResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Quiz ID cannot be empty",
					Details:   []string{"Please provide a valid quiz ID"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	updatedQuiz, err := s.db.UpdateQuiz(req.Id, req.Quiz)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.UpdateQuizResponse{
				StatusCode: 404,
				Message:    "Quiz not found",
				Result: &pb.UpdateQuizResponse_Error{
					Error: &pb.ErrorDetails{
						Code:      404,
						Message:   "Quiz with the specified ID does not exist",
						Details:   []string{fmt.Sprintf("Quiz ID: %s", req.Id)},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			}, nil
		}

		return &pb.UpdateQuizResponse{
			StatusCode: 500,
			Message:    "Failed to update quiz",
			Result: &pb.UpdateQuizResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Internal server error",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.UpdateQuizResponse{
		StatusCode: 200,
		Message:    "Quiz updated successfully",
		Result: &pb.UpdateQuizResponse_Quiz{
			Quiz: updatedQuiz,
		},
	}, nil
}

// DeleteQuiz deletes a quiz
func (s *QuizServer) DeleteQuiz(ctx context.Context, req *pb.DeleteQuizRequest) (*pb.DeleteQuizResponse, error) {
	log.Printf("DeleteQuiz called with ID: %s", req.Id)

	if req.Id == "" {
		return &pb.DeleteQuizResponse{
			StatusCode: 400,
			Message:    "Quiz ID is required",
			Result: &pb.DeleteQuizResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Quiz ID cannot be empty",
					Details:   []string{"Please provide a valid quiz ID"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	err := s.db.DeleteQuiz(req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.DeleteQuizResponse{
				StatusCode: 404,
				Message:    "Quiz not found",
				Result: &pb.DeleteQuizResponse_Error{
					Error: &pb.ErrorDetails{
						Code:      404,
						Message:   "Quiz with the specified ID does not exist",
						Details:   []string{fmt.Sprintf("Quiz ID: %s", req.Id)},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			}, nil
		}

		return &pb.DeleteQuizResponse{
			StatusCode: 500,
			Message:    "Failed to delete quiz",
			Result: &pb.DeleteQuizResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Internal server error",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.DeleteQuizResponse{
		StatusCode: 200,
		Message:    "Quiz deleted successfully",
		Result: &pb.DeleteQuizResponse_Success{
			Success: true,
		},
	}, nil
}

// SubmitQuizAnswer submits an answer for a quiz
func (s *QuizServer) SubmitQuizAnswer(ctx context.Context, req *pb.SubmitQuizAnswerRequest) (*pb.SubmitQuizAnswerResponse, error) {
	log.Printf("SubmitQuizAnswer called for quiz: %s by user: %s", req.QuizId, req.UserId)

	if req.UserId == "" || req.QuizId == "" || req.SelectedOption == "" {
		return &pb.SubmitQuizAnswerResponse{
			StatusCode: 400,
			Message:    "User ID, Quiz ID, and selected option are required",
			Result: &pb.SubmitQuizAnswerResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Missing required fields",
					Details:   []string{"User ID, Quiz ID, and selected option cannot be empty"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	response, err := s.db.SubmitQuizAnswer(req.UserId, req.QuizId, req.SelectedOption)
	if err != nil {
		return &pb.SubmitQuizAnswerResponse{
			StatusCode: 500,
			Message:    "Failed to submit quiz answer",
			Result: &pb.SubmitQuizAnswerResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Internal server error",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.SubmitQuizAnswerResponse{
		StatusCode: 200,
		Message:    "Quiz answer submitted successfully",
		Result: &pb.SubmitQuizAnswerResponse_Response{
			Response: response,
		},
	}, nil
}

// GetUserQuizHistory retrieves quiz history for a user
func (s *QuizServer) GetUserQuizHistory(ctx context.Context, req *pb.GetUserQuizHistoryRequest) (*pb.GetUserQuizHistoryResponse, error) {
	log.Printf("GetUserQuizHistory called for user: %s", req.UserId)

	if req.UserId == "" {
		return &pb.GetUserQuizHistoryResponse{
			StatusCode: 400,
			Message:    "User ID is required",
			Result: &pb.GetUserQuizHistoryResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "User ID cannot be empty",
					Details:   []string{"Please provide a valid user ID"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	history, err := s.db.GetUserQuizHistory(req.UserId)
	if err != nil {
		return &pb.GetUserQuizHistoryResponse{
			StatusCode: 500,
			Message:    "Failed to retrieve quiz history",
			Result: &pb.GetUserQuizHistoryResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Internal server error",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.GetUserQuizHistoryResponse{
		StatusCode: 200,
		Message:    "Quiz history retrieved successfully",
		Result: &pb.GetUserQuizHistoryResponse_Responses{
			Responses: &pb.QuizResponseList{
				Responses:  history,
				TotalCount: int32(len(history)),
				Page:       req.Page,
				Limit:      req.Limit,
			},
		},
	}, nil
}

// GetQuizStatistics retrieves statistics for a quiz
func (s *QuizServer) GetQuizStatistics(ctx context.Context, req *pb.GetQuizStatisticsRequest) (*pb.GetQuizStatisticsResponse, error) {
	log.Printf("GetQuizStatistics called for quiz: %s", req.QuizId)

	if req.QuizId == "" {
		return &pb.GetQuizStatisticsResponse{
			StatusCode: 400,
			Message:    "Quiz ID is required",
			Result: &pb.GetQuizStatisticsResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Quiz ID cannot be empty",
					Details:   []string{"Please provide a valid quiz ID"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	stats, err := s.db.GetQuizStatistics(req.QuizId)
	if err != nil {
		return &pb.GetQuizStatisticsResponse{
			StatusCode: 500,
			Message:    "Failed to retrieve quiz statistics",
			Result: &pb.GetQuizStatisticsResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Internal server error",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.GetQuizStatisticsResponse{
		StatusCode: 200,
		Message:    "Quiz statistics retrieved successfully",
		Result: &pb.GetQuizStatisticsResponse_Statistics{
			Statistics: stats,
		},
	}, nil
}
