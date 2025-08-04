package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	pb "history-template-services/proto"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type HistoryTemplateServer struct {
	pb.UnimplementedHistoryTemplateServiceServer
	db *Database
}

func NewHistoryTemplateServer() *HistoryTemplateServer {
	return &HistoryTemplateServer{
		db: NewDatabase(),
	}
}

// GetTemplateByID retrieves a history template by its ID
func (s *HistoryTemplateServer) GetTemplateByID(ctx context.Context, req *pb.GetTemplateByIDRequest) (*pb.GetTemplateByIDResponse, error) {
	log.Printf("GetTemplateByID called with ID: %s", req.Id)

	if req.Id == "" {
		return &pb.GetTemplateByIDResponse{
			StatusCode: 400,
			Message:    "Template ID is required",
			Result: &pb.GetTemplateByIDResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Template ID cannot be empty",
					Details:   []string{"Please provide a valid template ID"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	template, err := s.db.GetTemplateByID(req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.GetTemplateByIDResponse{
				StatusCode: 404,
				Message:    "Template not found",
				Result: &pb.GetTemplateByIDResponse_Error{
					Error: &pb.ErrorDetails{
						Code:      404,
						Message:   "Template with the specified ID does not exist",
						Details:   []string{fmt.Sprintf("No template found with ID: %s", req.Id)},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			}, nil
		}

		log.Printf("Database error: %v", err)
		return &pb.GetTemplateByIDResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result: &pb.GetTemplateByIDResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Failed to retrieve template",
					Details:   []string{"Database operation failed"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.GetTemplateByIDResponse{
		StatusCode: 200,
		Message:    "Template retrieved successfully",
		Result: &pb.GetTemplateByIDResponse_Template{
			Template: template,
		},
	}, nil
}

// GetAllTemplates retrieves all history templates with pagination
func (s *HistoryTemplateServer) GetAllTemplates(ctx context.Context, req *pb.GetAllTemplatesRequest) (*pb.GetAllTemplatesResponse, error) {
	log.Printf("GetAllTemplates called with page: %d, limit: %d", req.Page, req.Limit)

	// Set default pagination values
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // Maximum limit
	}

	templates, totalCount, err := s.db.GetAllTemplates(page, limit, req.SortBy, req.Order)
	if err != nil {
		log.Printf("Database error: %v", err)
		return &pb.GetAllTemplatesResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result: &pb.GetAllTemplatesResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Failed to retrieve templates",
					Details:   []string{"Database operation failed"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.GetAllTemplatesResponse{
		StatusCode: 200,
		Message:    "Templates retrieved successfully",
		Result: &pb.GetAllTemplatesResponse_Templates{
			Templates: &pb.TemplateList{
				Templates:  templates,
				TotalCount: totalCount,
				Page:       page,
				Limit:      limit,
			},
		},
	}, nil
}

// CreateTemplate creates a new history template
func (s *HistoryTemplateServer) CreateTemplate(ctx context.Context, req *pb.CreateTemplateRequest) (*pb.CreateTemplateResponse, error) {
	log.Printf("CreateTemplate called with title: %s", req.Template.Title)

	if req.Template == nil {
		return &pb.CreateTemplateResponse{
			StatusCode: 400,
			Message:    "Template data is required",
			Result: &pb.CreateTemplateResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Template object cannot be nil",
					Details:   []string{"Please provide valid template data"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Validate required fields
	if req.Template.Title == "" {
		return &pb.CreateTemplateResponse{
			StatusCode: 400,
			Message:    "Template title is required",
			Result: &pb.CreateTemplateResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Title cannot be empty",
					Details:   []string{"Please provide a valid title"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Generate ID and set timestamps
	req.Template.Id = uuid.New().String()
	now := timestamppb.Now()
	req.Template.CreatedAt = now
	req.Template.UpdatedAt = now

	createdTemplate, err := s.db.CreateTemplate(req.Template)
	if err != nil {
		log.Printf("Database error: %v", err)
		return &pb.CreateTemplateResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result: &pb.CreateTemplateResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Failed to create template",
					Details:   []string{"Database operation failed"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.CreateTemplateResponse{
		StatusCode: 201,
		Message:    "Template created successfully",
		Result: &pb.CreateTemplateResponse_Template{
			Template: createdTemplate,
		},
	}, nil
}

// UpdateTemplate updates an existing history template
func (s *HistoryTemplateServer) UpdateTemplate(ctx context.Context, req *pb.UpdateTemplateRequest) (*pb.UpdateTemplateResponse, error) {
	log.Printf("UpdateTemplate called with ID: %s", req.Id)

	if req.Id == "" {
		return &pb.UpdateTemplateResponse{
			StatusCode: 400,
			Message:    "Template ID is required",
			Result: &pb.UpdateTemplateResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Template ID cannot be empty",
					Details:   []string{"Please provide a valid template ID"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	if req.Template == nil {
		return &pb.UpdateTemplateResponse{
			StatusCode: 400,
			Message:    "Template data is required",
			Result: &pb.UpdateTemplateResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Template object cannot be nil",
					Details:   []string{"Please provide valid template data"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Set updated timestamp
	req.Template.UpdatedAt = timestamppb.Now()

	updatedTemplate, err := s.db.UpdateTemplate(req.Id, req.Template)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.UpdateTemplateResponse{
				StatusCode: 404,
				Message:    "Template not found",
				Result: &pb.UpdateTemplateResponse_Error{
					Error: &pb.ErrorDetails{
						Code:      404,
						Message:   "Template with the specified ID does not exist",
						Details:   []string{fmt.Sprintf("No template found with ID: %s", req.Id)},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			}, nil
		}

		log.Printf("Database error: %v", err)
		return &pb.UpdateTemplateResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result: &pb.UpdateTemplateResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Failed to update template",
					Details:   []string{"Database operation failed"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.UpdateTemplateResponse{
		StatusCode: 200,
		Message:    "Template updated successfully",
		Result: &pb.UpdateTemplateResponse_Template{
			Template: updatedTemplate,
		},
	}, nil
}

// DeleteTemplate deletes a history template
func (s *HistoryTemplateServer) DeleteTemplate(ctx context.Context, req *pb.DeleteTemplateRequest) (*pb.DeleteTemplateResponse, error) {
	log.Printf("DeleteTemplate called with ID: %s", req.Id)

	if req.Id == "" {
		return &pb.DeleteTemplateResponse{
			StatusCode: 400,
			Message:    "Template ID is required",
			Result: &pb.DeleteTemplateResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Template ID cannot be empty",
					Details:   []string{"Please provide a valid template ID"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	err := s.db.DeleteTemplate(req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.DeleteTemplateResponse{
				StatusCode: 404,
				Message:    "Template not found",
				Result: &pb.DeleteTemplateResponse_Error{
					Error: &pb.ErrorDetails{
						Code:      404,
						Message:   "Template with the specified ID does not exist",
						Details:   []string{fmt.Sprintf("No template found with ID: %s", req.Id)},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			}, nil
		}

		log.Printf("Database error: %v", err)
		return &pb.DeleteTemplateResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result: &pb.DeleteTemplateResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Failed to delete template",
					Details:   []string{"Database operation failed"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.DeleteTemplateResponse{
		StatusCode: 200,
		Message:    "Template deleted successfully",
		Result: &pb.DeleteTemplateResponse_Success{
			Success: true,
		},
	}, nil
}

// GetTemplatesByEra retrieves templates by era
func (s *HistoryTemplateServer) GetTemplatesByEra(ctx context.Context, req *pb.GetTemplatesByEraRequest) (*pb.GetTemplatesByEraResponse, error) {
	log.Printf("GetTemplatesByEra called with era: %s", req.Era)

	if req.Era == "" {
		return &pb.GetTemplatesByEraResponse{
			StatusCode: 400,
			Message:    "Era is required",
			Result: &pb.GetTemplatesByEraResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Era cannot be empty",
					Details:   []string{"Please provide a valid era"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Set default pagination values
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

	templates, totalCount, err := s.db.GetTemplatesByEra(req.Era, page, limit)
	if err != nil {
		log.Printf("Database error: %v", err)
		return &pb.GetTemplatesByEraResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result: &pb.GetTemplatesByEraResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Failed to retrieve templates",
					Details:   []string{"Database operation failed"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.GetTemplatesByEraResponse{
		StatusCode: 200,
		Message:    "Templates retrieved successfully",
		Result: &pb.GetTemplatesByEraResponse_Templates{
			Templates: &pb.TemplateList{
				Templates:  templates,
				TotalCount: totalCount,
				Page:       page,
				Limit:      limit,
			},
		},
	}, nil
}

// GetTemplatesByDynasty retrieves templates by dynasty
func (s *HistoryTemplateServer) GetTemplatesByDynasty(ctx context.Context, req *pb.GetTemplatesByDynastyRequest) (*pb.GetTemplatesByDynastyResponse, error) {
	log.Printf("GetTemplatesByDynasty called with dynasty: %d", req.Dynasty)

	if req.Dynasty <= 0 {
		return &pb.GetTemplatesByDynastyResponse{
			StatusCode: 400,
			Message:    "Dynasty is required",
			Result: &pb.GetTemplatesByDynastyResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Dynasty must be a positive number",
					Details:   []string{"Please provide a valid dynasty number"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Set default pagination values
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

	templates, totalCount, err := s.db.GetTemplatesByDynasty(req.Dynasty, page, limit)
	if err != nil {
		log.Printf("Database error: %v", err)
		return &pb.GetTemplatesByDynastyResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result: &pb.GetTemplatesByDynastyResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Failed to retrieve templates",
					Details:   []string{"Database operation failed"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.GetTemplatesByDynastyResponse{
		StatusCode: 200,
		Message:    "Templates retrieved successfully",
		Result: &pb.GetTemplatesByDynastyResponse_Templates{
			Templates: &pb.TemplateList{
				Templates:  templates,
				TotalCount: totalCount,
				Page:       page,
				Limit:      limit,
			},
		},
	}, nil
}

// GetTemplatesByPharaoh retrieves templates by pharaoh ID
func (s *HistoryTemplateServer) GetTemplatesByPharaoh(ctx context.Context, req *pb.GetTemplatesByPharaohRequest) (*pb.GetTemplatesByPharaohResponse, error) {
	log.Printf("GetTemplatesByPharaoh called with pharaoh_id: %s", req.PharaohId)

	if req.PharaohId == "" {
		return &pb.GetTemplatesByPharaohResponse{
			StatusCode: 400,
			Message:    "Pharaoh ID is required",
			Result: &pb.GetTemplatesByPharaohResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Pharaoh ID cannot be empty",
					Details:   []string{"Please provide a valid pharaoh ID"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Set default pagination values
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

	templates, totalCount, err := s.db.GetTemplatesByPharaoh(req.PharaohId, page, limit)
	if err != nil {
		log.Printf("Database error: %v", err)
		return &pb.GetTemplatesByPharaohResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result: &pb.GetTemplatesByPharaohResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Failed to retrieve templates",
					Details:   []string{"Database operation failed"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.GetTemplatesByPharaohResponse{
		StatusCode: 200,
		Message:    "Templates retrieved successfully",
		Result: &pb.GetTemplatesByPharaohResponse_Templates{
			Templates: &pb.TemplateList{
				Templates:  templates,
				TotalCount: totalCount,
				Page:       page,
				Limit:      limit,
			},
		},
	}, nil
}

// GetTemplatesByDifficulty retrieves templates by difficulty
func (s *HistoryTemplateServer) GetTemplatesByDifficulty(ctx context.Context, req *pb.GetTemplatesByDifficultyRequest) (*pb.GetTemplatesByDifficultyResponse, error) {
	log.Printf("GetTemplatesByDifficulty called with difficulty: %s", req.Difficulty)

	if req.Difficulty == "" {
		return &pb.GetTemplatesByDifficultyResponse{
			StatusCode: 400,
			Message:    "Difficulty is required",
			Result: &pb.GetTemplatesByDifficultyResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Difficulty cannot be empty",
					Details:   []string{"Please provide a valid difficulty level"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Validate difficulty
	if !ValidateDifficulty(req.Difficulty) {
		return &pb.GetTemplatesByDifficultyResponse{
			StatusCode: 400,
			Message:    "Invalid difficulty level",
			Result: &pb.GetTemplatesByDifficultyResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Difficulty must be one of: beginner, intermediate, advanced",
					Details:   []string{"Please provide a valid difficulty level"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Set default pagination values
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

	templates, totalCount, err := s.db.GetTemplatesByDifficulty(req.Difficulty, page, limit)
	if err != nil {
		log.Printf("Database error: %v", err)
		return &pb.GetTemplatesByDifficultyResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result: &pb.GetTemplatesByDifficultyResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Failed to retrieve templates",
					Details:   []string{"Database operation failed"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.GetTemplatesByDifficultyResponse{
		StatusCode: 200,
		Message:    "Templates retrieved successfully",
		Result: &pb.GetTemplatesByDifficultyResponse_Templates{
			Templates: &pb.TemplateList{
				Templates:  templates,
				TotalCount: totalCount,
				Page:       page,
				Limit:      limit,
			},
		},
	}, nil
}

// SearchTemplates searches templates by query
func (s *HistoryTemplateServer) SearchTemplates(ctx context.Context, req *pb.SearchTemplatesRequest) (*pb.SearchTemplatesResponse, error) {
	log.Printf("SearchTemplates called with query: %s", req.Query)

	if req.Query == "" {
		return &pb.SearchTemplatesResponse{
			StatusCode: 400,
			Message:    "Search query is required",
			Result: &pb.SearchTemplatesResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Query cannot be empty",
					Details:   []string{"Please provide a search query"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Set default pagination values
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

	templates, totalCount, err := s.db.SearchTemplates(req.Query, req.Fields, page, limit)
	if err != nil {
		log.Printf("Database error: %v", err)
		return &pb.SearchTemplatesResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result: &pb.SearchTemplatesResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Failed to search templates",
					Details:   []string{"Database operation failed"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.SearchTemplatesResponse{
		StatusCode: 200,
		Message:    "Templates found successfully",
		Result: &pb.SearchTemplatesResponse_Templates{
			Templates: &pb.TemplateList{
				Templates:  templates,
				TotalCount: totalCount,
				Page:       page,
				Limit:      limit,
			},
		},
	}, nil
}

// GetTemplatesByTag retrieves templates by tag
func (s *HistoryTemplateServer) GetTemplatesByTag(ctx context.Context, req *pb.GetTemplatesByTagRequest) (*pb.GetTemplatesByTagResponse, error) {
	log.Printf("GetTemplatesByTag called with tag: %s", req.Tag)

	if req.Tag == "" {
		return &pb.GetTemplatesByTagResponse{
			StatusCode: 400,
			Message:    "Tag is required",
			Result: &pb.GetTemplatesByTagResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Tag cannot be empty",
					Details:   []string{"Please provide a valid tag"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Set default pagination values
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

	templates, totalCount, err := s.db.GetTemplatesByTag(req.Tag, page, limit)
	if err != nil {
		log.Printf("Database error: %v", err)
		return &pb.GetTemplatesByTagResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result: &pb.GetTemplatesByTagResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Failed to retrieve templates",
					Details:   []string{"Database operation failed"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.GetTemplatesByTagResponse{
		StatusCode: 200,
		Message:    "Templates retrieved successfully",
		Result: &pb.GetTemplatesByTagResponse_Templates{
			Templates: &pb.TemplateList{
				Templates:  templates,
				TotalCount: totalCount,
				Page:       page,
				Limit:      limit,
			},
		},
	}, nil
}

// GetRelatedTemplates retrieves related templates
func (s *HistoryTemplateServer) GetRelatedTemplates(ctx context.Context, req *pb.GetRelatedTemplatesRequest) (*pb.GetRelatedTemplatesResponse, error) {
	log.Printf("GetRelatedTemplates called with template_id: %s", req.TemplateId)

	if req.TemplateId == "" {
		return &pb.GetRelatedTemplatesResponse{
			StatusCode: 400,
			Message:    "Template ID is required",
			Result: &pb.GetRelatedTemplatesResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Template ID cannot be empty",
					Details:   []string{"Please provide a valid template ID"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 5
	}
	if limit > 20 {
		limit = 20
	}

	templates, totalCount, err := s.db.GetRelatedTemplates(req.TemplateId, limit)
	if err != nil {
		log.Printf("Database error: %v", err)
		return &pb.GetRelatedTemplatesResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result: &pb.GetRelatedTemplatesResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Failed to retrieve related templates",
					Details:   []string{"Database operation failed"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.GetRelatedTemplatesResponse{
		StatusCode: 200,
		Message:    "Related templates retrieved successfully",
		Result: &pb.GetRelatedTemplatesResponse_Templates{
			Templates: &pb.TemplateList{
				Templates:  templates,
				TotalCount: totalCount,
				Page:       1,
				Limit:      limit,
			},
		},
	}, nil
}
