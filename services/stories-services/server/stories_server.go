package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	pb "stories-services/proto"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type StoriesServer struct {
	pb.UnimplementedStoriesServiceServer
	db    *Database
	cache *RedisCache
}

func NewStoriesServer() *StoriesServer {
	return &StoriesServer{
		db:    NewDatabase(),
		cache: NewRedisCache(),
	}
}

// MARK: CreateStory
func (s *StoriesServer) CreateStory(ctx context.Context, req *pb.CreateStoryRequest) (*pb.CreateStoryResponse, error) {
	log.Printf("CreateStory called for user: %s", req.UserId)

	// Validate required fields
	if req.UserId == "" || req.Title == "" {
		return &pb.CreateStoryResponse{
			StatusCode: 400,
			Message:    "User ID and title are required",
			Result: &pb.CreateStoryResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Missing required fields",
					Details:   []string{"User ID and title are required"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	storyID := uuid.New().String()
	language := req.Language
	if language == "" {
		language = "en"
	}

	query := `
		INSERT INTO stories (id, user_id, title, content, author, category, tags, image_url, language, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`

	var createdAt, updatedAt time.Time
	err := s.db.DB.QueryRow(query, storyID, req.UserId, req.Title, req.Content, req.Author,
		req.Category, pq.Array(req.Tags), req.ImageUrl, language, time.Now(), time.Now()).Scan(
		&storyID, &createdAt, &updatedAt)

	if err != nil {
		log.Printf("Database error: %v", err)
		return &pb.CreateStoryResponse{
			StatusCode: 500,
			Message:    "Failed to create story",
			Result: &pb.CreateStoryResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Database error occurred",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Clear cache for user stories and all stories
	if s.cache != nil {
		s.cache.DeletePattern(fmt.Sprintf("stories:user:%s:*", req.UserId))
		s.cache.DeletePattern("stories:all:*")
	}

	story := &pb.Story{
		Id:         storyID,
		UserId:     req.UserId,
		Title:      req.Title,
		Content:    req.Content,
		Author:     req.Author,
		Status:     "DRAFT",
		Category:   req.Category,
		Tags:       req.Tags,
		ImageUrl:   req.ImageUrl,
		Language:   language,
		IsFeatured: false,
		ViewCount:  0,
		LikeCount:  0,
		CreatedAt:  createdAt.Format(time.RFC3339),
		UpdatedAt:  updatedAt.Format(time.RFC3339),
	}

	return &pb.CreateStoryResponse{
		StatusCode: 201,
		Message:    "Story created successfully",
		Result: &pb.CreateStoryResponse_Story{
			Story: story,
		},
	}, nil
}

// MARK: GetStoryById (with Redis caching)
func (s *StoriesServer) GetStoryById(ctx context.Context, req *pb.GetStoryByIdRequest) (*pb.GetStoryByIdResponse, error) {
	log.Printf("GetStoryById called with ID: %s", req.StoryId)

	if req.StoryId == "" {
		return &pb.GetStoryByIdResponse{
			StatusCode: 400,
			Message:    "Story ID is required",
			Result: &pb.GetStoryByIdResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Missing story ID",
					Details:   []string{"Story ID is required"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Try to get from cache first
	cacheKey := fmt.Sprintf("story:%s", req.StoryId)
	var story pb.Story
	if s.cache != nil {
		if err := s.cache.Get(cacheKey, &story); err == nil {
			log.Printf("Story found in cache: %s", req.StoryId)
			// Increment view count
			go s.incrementViewCount(req.StoryId)
			return &pb.GetStoryByIdResponse{
				StatusCode: 200,
				Message:    "Story retrieved from cache",
				Result: &pb.GetStoryByIdResponse_Story{
					Story: &story,
				},
			}, nil
		}
	}

	// Get from database
	query := `
		SELECT id, user_id, title, content, author, status, category, tags, image_url, 
		       language, is_featured, view_count, like_count, published_at, created_at, updated_at
		FROM stories WHERE id = $1
	`

	row := s.db.DB.QueryRow(query, req.StoryId)

	var tags pq.StringArray
	var publishedAt sql.NullTime
	var createdAt, updatedAt time.Time

	err := row.Scan(&story.Id, &story.UserId, &story.Title, &story.Content, &story.Author,
		&story.Status, &story.Category, &tags, &story.ImageUrl, &story.Language,
		&story.IsFeatured, &story.ViewCount, &story.LikeCount, &publishedAt, &createdAt, &updatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.GetStoryByIdResponse{
				StatusCode: 404,
				Message:    "Story not found",
				Result: &pb.GetStoryByIdResponse_Error{
					Error: &pb.ErrorDetails{
						Code:      404,
						Message:   "Story not found",
						Details:   []string{"No story found with this ID"},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			}, nil
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	story.Tags = []string(tags)
	story.CreatedAt = createdAt.Format(time.RFC3339)
	story.UpdatedAt = updatedAt.Format(time.RFC3339)
	if publishedAt.Valid {
		story.PublishedAt = publishedAt.Time.Format(time.RFC3339)
	}

	// Cache the story for 1 hour
	if s.cache != nil {
		s.cache.Set(cacheKey, story, time.Hour)
	}

	// Increment view count
	go s.incrementViewCount(req.StoryId)

	return &pb.GetStoryByIdResponse{
		StatusCode: 200,
		Message:    "Story retrieved successfully",
		Result: &pb.GetStoryByIdResponse_Story{
			Story: &story,
		},
	}, nil
}

// MARK: GetStoriesByUser (with Redis caching)
func (s *StoriesServer) GetStoriesByUser(ctx context.Context, req *pb.GetStoriesByUserRequest) (*pb.GetStoriesByUserResponse, error) {
	log.Printf("GetStoriesByUser called for user: %s", req.UserId)

	if req.UserId == "" {
		return &pb.GetStoriesByUserResponse{
			StatusCode: 400,
			Message:    "User ID is required",
			Result: &pb.GetStoriesByUserResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Missing user ID",
					Details:   []string{"User ID is required"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Set defaults
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	// Try to get from cache first
	cacheKey := fmt.Sprintf("stories:user:%s:page:%d:limit:%d:status:%s", req.UserId, page, limit, req.Status)
	var storiesList pb.StoriesList
	if s.cache != nil {
		if err := s.cache.Get(cacheKey, &storiesList); err == nil {
			log.Printf("User stories found in cache: %s", req.UserId)
			return &pb.GetStoriesByUserResponse{
				StatusCode: 200,
				Message:    "Stories retrieved from cache",
				Result: &pb.GetStoriesByUserResponse_Stories{
					Stories: &storiesList,
				},
			}, nil
		}
	}

	// Build query
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(`
		SELECT id, user_id, title, content, author, status, category, tags, image_url, 
		       language, is_featured, view_count, like_count, published_at, created_at, updated_at
		FROM stories WHERE user_id = $1
	`)

	args := []interface{}{req.UserId}
	argIndex := 2

	if req.Status != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND status = $%d", argIndex))
		args = append(args, req.Status)
		argIndex++
	}

	queryBuilder.WriteString(" ORDER BY created_at DESC")
	queryBuilder.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1))
	args = append(args, limit, (page-1)*limit)

	// Get total count
	countQuery := "SELECT COUNT(*) FROM stories WHERE user_id = $1"
	countArgs := []interface{}{req.UserId}
	if req.Status != "" {
		countQuery += " AND status = $2"
		countArgs = append(countArgs, req.Status)
	}

	var totalCount int32
	err := s.db.DB.QueryRow(countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %v", err)
	}

	// Execute main query
	rows, err := s.db.DB.Query(queryBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("database error: %v", err)
	}
	defer rows.Close()

	var stories []*pb.Story
	for rows.Next() {
		var story pb.Story
		var tags pq.StringArray
		var publishedAt sql.NullTime
		var createdAt, updatedAt time.Time

		err := rows.Scan(&story.Id, &story.UserId, &story.Title, &story.Content, &story.Author,
			&story.Status, &story.Category, &tags, &story.ImageUrl, &story.Language,
			&story.IsFeatured, &story.ViewCount, &story.LikeCount, &publishedAt, &createdAt, &updatedAt)

		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}

		story.Tags = []string(tags)
		story.CreatedAt = createdAt.Format(time.RFC3339)
		story.UpdatedAt = updatedAt.Format(time.RFC3339)
		if publishedAt.Valid {
			story.PublishedAt = publishedAt.Time.Format(time.RFC3339)
		}

		stories = append(stories, &story)
	}

	storiesList = pb.StoriesList{
		Stories:    stories,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
	}

	// Cache for 30 minutes
	if s.cache != nil {
		s.cache.Set(cacheKey, storiesList, 30*time.Minute)
	}

	return &pb.GetStoriesByUserResponse{
		StatusCode: 200,
		Message:    "Stories retrieved successfully",
		Result: &pb.GetStoriesByUserResponse_Stories{
			Stories: &storiesList,
		},
	}, nil
}

// MARK: GetAllStories (with Redis caching)
func (s *StoriesServer) GetAllStories(ctx context.Context, req *pb.GetAllStoriesRequest) (*pb.GetAllStoriesResponse, error) {
	log.Printf("GetAllStories called")

	// Set defaults
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	// Try to get from cache first
	cacheKey := fmt.Sprintf("stories:all:page:%d:limit:%d:category:%s:status:%s:featured:%t",
		page, limit, req.Category, req.Status, req.FeaturedOnly)
	var storiesList pb.StoriesList
	if s.cache != nil {
		if err := s.cache.Get(cacheKey, &storiesList); err == nil {
			log.Printf("All stories found in cache")
			return &pb.GetAllStoriesResponse{
				StatusCode: 200,
				Message:    "Stories retrieved from cache",
				Result: &pb.GetAllStoriesResponse_Stories{
					Stories: &storiesList,
				},
			}, nil
		}
	}

	// Build query
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(`
		SELECT id, user_id, title, content, author, status, category, tags, image_url, 
		       language, is_featured, view_count, like_count, published_at, created_at, updated_at
		FROM stories WHERE 1=1
	`)

	args := []interface{}{}
	argIndex := 1

	if req.Category != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND category = $%d", argIndex))
		args = append(args, req.Category)
		argIndex++
	}

	if req.Status != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND status = $%d", argIndex))
		args = append(args, req.Status)
		argIndex++
	}

	if req.FeaturedOnly {
		queryBuilder.WriteString(" AND is_featured = true")
	}

	queryBuilder.WriteString(" ORDER BY created_at DESC")
	queryBuilder.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1))
	args = append(args, limit, (page-1)*limit)

	// Get total count
	countQueryBuilder := strings.Builder{}
	countQueryBuilder.WriteString("SELECT COUNT(*) FROM stories WHERE 1=1")
	countArgs := []interface{}{}
	countArgIndex := 1

	if req.Category != "" {
		countQueryBuilder.WriteString(fmt.Sprintf(" AND category = $%d", countArgIndex))
		countArgs = append(countArgs, req.Category)
		countArgIndex++
	}

	if req.Status != "" {
		countQueryBuilder.WriteString(fmt.Sprintf(" AND status = $%d", countArgIndex))
		countArgs = append(countArgs, req.Status)
		countArgIndex++
	}

	if req.FeaturedOnly {
		countQueryBuilder.WriteString(" AND is_featured = true")
	}

	var totalCount int32
	err := s.db.DB.QueryRow(countQueryBuilder.String(), countArgs...).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %v", err)
	}

	// Execute main query
	rows, err := s.db.DB.Query(queryBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("database error: %v", err)
	}
	defer rows.Close()

	var stories []*pb.Story
	for rows.Next() {
		var story pb.Story
		var tags pq.StringArray
		var publishedAt sql.NullTime
		var createdAt, updatedAt time.Time

		err := rows.Scan(&story.Id, &story.UserId, &story.Title, &story.Content, &story.Author,
			&story.Status, &story.Category, &tags, &story.ImageUrl, &story.Language,
			&story.IsFeatured, &story.ViewCount, &story.LikeCount, &publishedAt, &createdAt, &updatedAt)

		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}

		story.Tags = []string(tags)
		story.CreatedAt = createdAt.Format(time.RFC3339)
		story.UpdatedAt = updatedAt.Format(time.RFC3339)
		if publishedAt.Valid {
			story.PublishedAt = publishedAt.Time.Format(time.RFC3339)
		}

		stories = append(stories, &story)
	}

	storiesList = pb.StoriesList{
		Stories:    stories,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
	}

	// Cache for 15 minutes
	if s.cache != nil {
		s.cache.Set(cacheKey, storiesList, 15*time.Minute)
	}

	return &pb.GetAllStoriesResponse{
		StatusCode: 200,
		Message:    "Stories retrieved successfully",
		Result: &pb.GetAllStoriesResponse_Stories{
			Stories: &storiesList,
		},
	}, nil
}

// MARK: UpdateStory
func (s *StoriesServer) UpdateStory(ctx context.Context, req *pb.UpdateStoryRequest) (*pb.UpdateStoryResponse, error) {
	log.Printf("UpdateStory called for story: %s", req.StoryId)

	if req.StoryId == "" {
		return &pb.UpdateStoryResponse{
			StatusCode: 400,
			Message:    "Story ID is required",
			Result: &pb.UpdateStoryResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Missing story ID",
					Details:   []string{"Story ID is required"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Build update query dynamically
	updateFields := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Title != "" {
		updateFields = append(updateFields, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, req.Title)
		argIndex++
	}

	if req.Content != "" {
		updateFields = append(updateFields, fmt.Sprintf("content = $%d", argIndex))
		args = append(args, req.Content)
		argIndex++
	}

	if req.Category != "" {
		updateFields = append(updateFields, fmt.Sprintf("category = $%d", argIndex))
		args = append(args, req.Category)
		argIndex++
	}

	if len(req.Tags) > 0 {
		updateFields = append(updateFields, fmt.Sprintf("tags = $%d", argIndex))
		args = append(args, pq.Array(req.Tags))
		argIndex++
	}

	if req.ImageUrl != "" {
		updateFields = append(updateFields, fmt.Sprintf("image_url = $%d", argIndex))
		args = append(args, req.ImageUrl)
		argIndex++
	}

	if len(updateFields) == 0 {
		return &pb.UpdateStoryResponse{
			StatusCode: 400,
			Message:    "No fields to update",
			Result: &pb.UpdateStoryResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "No update fields provided",
					Details:   []string{"At least one field must be provided for update"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Add updated_at
	updateFields = append(updateFields, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	// Add WHERE clause
	args = append(args, req.StoryId)

	query := fmt.Sprintf(`
		UPDATE stories SET %s WHERE id = $%d
		RETURNING id, user_id, title, content, author, status, category, tags, image_url, 
		          language, is_featured, view_count, like_count, published_at, created_at, updated_at
	`, strings.Join(updateFields, ", "), argIndex)

	row := s.db.DB.QueryRow(query, args...)

	var story pb.Story
	var tags pq.StringArray
	var publishedAt sql.NullTime
	var createdAt, updatedAt time.Time

	err := row.Scan(&story.Id, &story.UserId, &story.Title, &story.Content, &story.Author,
		&story.Status, &story.Category, &tags, &story.ImageUrl, &story.Language,
		&story.IsFeatured, &story.ViewCount, &story.LikeCount, &publishedAt, &createdAt, &updatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.UpdateStoryResponse{
				StatusCode: 404,
				Message:    "Story not found",
				Result: &pb.UpdateStoryResponse_Error{
					Error: &pb.ErrorDetails{
						Code:      404,
						Message:   "Story not found",
						Details:   []string{"No story found with this ID"},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			}, nil
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	story.Tags = []string(tags)
	story.CreatedAt = createdAt.Format(time.RFC3339)
	story.UpdatedAt = updatedAt.Format(time.RFC3339)
	if publishedAt.Valid {
		story.PublishedAt = publishedAt.Time.Format(time.RFC3339)
	}

	// Clear cache
	if s.cache != nil {
		s.cache.Delete(fmt.Sprintf("story:%s", req.StoryId))
		s.cache.DeletePattern(fmt.Sprintf("stories:user:%s:*", story.UserId))
		s.cache.DeletePattern("stories:all:*")
	}

	return &pb.UpdateStoryResponse{
		StatusCode: 200,
		Message:    "Story updated successfully",
		Result: &pb.UpdateStoryResponse_Story{
			Story: &story,
		},
	}, nil
}

// MARK: DeleteStory
func (s *StoriesServer) DeleteStory(ctx context.Context, req *pb.DeleteStoryRequest) (*pb.DeleteStoryResponse, error) {
	log.Printf("DeleteStory called for story: %s", req.StoryId)

	if req.StoryId == "" || req.UserId == "" {
		return &pb.DeleteStoryResponse{
			StatusCode: 400,
			Message:    "Story ID and User ID are required",
			Error: &pb.ErrorDetails{
				Code:      400,
				Message:   "Missing required fields",
				Details:   []string{"Story ID and User ID are required"},
				Timestamp: time.Now().Format(time.RFC3339),
			},
		}, nil
	}

	// Delete story (with user verification)
	query := `DELETE FROM stories WHERE id = $1 AND user_id = $2`
	result, err := s.db.DB.Exec(query, req.StoryId, req.UserId)
	if err != nil {
		return &pb.DeleteStoryResponse{
			StatusCode: 500,
			Message:    "Failed to delete story",
			Error: &pb.ErrorDetails{
				Code:      500,
				Message:   "Database error occurred",
				Details:   []string{err.Error()},
				Timestamp: time.Now().Format(time.RFC3339),
			},
		}, nil
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return &pb.DeleteStoryResponse{
			StatusCode: 404,
			Message:    "Story not found or unauthorized",
			Error: &pb.ErrorDetails{
				Code:      404,
				Message:   "Story not found or you don't have permission",
				Details:   []string{"No story found with this ID for this user"},
				Timestamp: time.Now().Format(time.RFC3339),
			},
		}, nil
	}

	// Clear cache
	if s.cache != nil {
		s.cache.Delete(fmt.Sprintf("story:%s", req.StoryId))
		s.cache.DeletePattern(fmt.Sprintf("stories:user:%s:*", req.UserId))
		s.cache.DeletePattern("stories:all:*")
	}

	return &pb.DeleteStoryResponse{
		StatusCode: 200,
		Message:    "Story deleted successfully",
	}, nil
}

// MARK: PublishStory
func (s *StoriesServer) PublishStory(ctx context.Context, req *pb.PublishStoryRequest) (*pb.PublishStoryResponse, error) {
	log.Printf("PublishStory called for story: %s", req.StoryId)

	if req.StoryId == "" || req.UserId == "" {
		return &pb.PublishStoryResponse{
			StatusCode: 400,
			Message:    "Story ID and User ID are required",
			Result: &pb.PublishStoryResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Missing required fields",
					Details:   []string{"Story ID and User ID are required"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Update story status to PUBLISHED
	query := `
		UPDATE stories SET status = 'PUBLISHED', published_at = $1, updated_at = $2 
		WHERE id = $3 AND user_id = $4
		RETURNING id, user_id, title, content, author, status, category, tags, image_url, 
		          language, is_featured, view_count, like_count, published_at, created_at, updated_at
	`

	now := time.Now()
	row := s.db.DB.QueryRow(query, now, now, req.StoryId, req.UserId)

	var story pb.Story
	var tags pq.StringArray
	var publishedAt sql.NullTime
	var createdAt, updatedAt time.Time

	err := row.Scan(&story.Id, &story.UserId, &story.Title, &story.Content, &story.Author,
		&story.Status, &story.Category, &tags, &story.ImageUrl, &story.Language,
		&story.IsFeatured, &story.ViewCount, &story.LikeCount, &publishedAt, &createdAt, &updatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.PublishStoryResponse{
				StatusCode: 404,
				Message:    "Story not found or unauthorized",
				Result: &pb.PublishStoryResponse_Error{
					Error: &pb.ErrorDetails{
						Code:      404,
						Message:   "Story not found or you don't have permission",
						Details:   []string{"No story found with this ID for this user"},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			}, nil
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	story.Tags = []string(tags)
	story.CreatedAt = createdAt.Format(time.RFC3339)
	story.UpdatedAt = updatedAt.Format(time.RFC3339)
	if publishedAt.Valid {
		story.PublishedAt = publishedAt.Time.Format(time.RFC3339)
	}

	// Clear cache
	if s.cache != nil {
		s.cache.Delete(fmt.Sprintf("story:%s", req.StoryId))
		s.cache.DeletePattern(fmt.Sprintf("stories:user:%s:*", req.UserId))
		s.cache.DeletePattern("stories:all:*")
	}

	return &pb.PublishStoryResponse{
		StatusCode: 200,
		Message:    "Story published successfully",
		Result: &pb.PublishStoryResponse_Story{
			Story: &story,
		},
	}, nil
}

// MARK: ReactToStory
func (s *StoriesServer) ReactToStory(ctx context.Context, req *pb.ReactToStoryRequest) (*pb.ReactToStoryResponse, error) {
	log.Printf("ReactToStory called for story: %s", req.StoryId)

	if req.StoryId == "" || req.UserId == "" || req.ReactionType == "" {
		return &pb.ReactToStoryResponse{
			StatusCode: 400,
			Message:    "Story ID, User ID, and Reaction Type are required",
			Error: &pb.ErrorDetails{
				Code:      400,
				Message:   "Missing required fields",
				Details:   []string{"Story ID, User ID, and Reaction Type are required"},
				Timestamp: time.Now().Format(time.RFC3339),
			},
		}, nil
	}

	// Insert or update reaction
	reactionQuery := `
		INSERT INTO story_reactions (id, story_id, user_id, reaction_type, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (story_id, user_id) 
		DO UPDATE SET reaction_type = $4, created_at = $5
	`

	reactionID := uuid.New().String()
	_, err := s.db.DB.Exec(reactionQuery, reactionID, req.StoryId, req.UserId, req.ReactionType, time.Now())
	if err != nil {
		return &pb.ReactToStoryResponse{
			StatusCode: 500,
			Message:    "Failed to add reaction",
			Error: &pb.ErrorDetails{
				Code:      500,
				Message:   "Database error occurred",
				Details:   []string{err.Error()},
				Timestamp: time.Now().Format(time.RFC3339),
			},
		}, nil
	}

	// Update like count in stories table
	go s.updateLikeCount(req.StoryId)

	// Clear cache
	if s.cache != nil {
		s.cache.Delete(fmt.Sprintf("story:%s", req.StoryId))
	}

	return &pb.ReactToStoryResponse{
		StatusCode: 200,
		Message:    "Reaction added successfully",
	}, nil
}

// Helper functions
func (s *StoriesServer) incrementViewCount(storyID string) {
	query := `UPDATE stories SET view_count = view_count + 1 WHERE id = $1`
	_, err := s.db.DB.Exec(query, storyID)
	if err != nil {
		log.Printf("Failed to increment view count: %v", err)
	}

	// Clear cache to ensure fresh data
	if s.cache != nil {
		s.cache.Delete(fmt.Sprintf("story:%s", storyID))
	}
}

func (s *StoriesServer) updateLikeCount(storyID string) {
	// Count likes for this story
	countQuery := `SELECT COUNT(*) FROM story_reactions WHERE story_id = $1 AND reaction_type = 'LIKE'`
	var likeCount int32
	err := s.db.DB.QueryRow(countQuery, storyID).Scan(&likeCount)
	if err != nil {
		log.Printf("Failed to count likes: %v", err)
		return
	}

	// Update story like count
	updateQuery := `UPDATE stories SET like_count = $1 WHERE id = $2`
	_, err = s.db.DB.Exec(updateQuery, likeCount, storyID)
	if err != nil {
		log.Printf("Failed to update like count: %v", err)
	}

	// Clear cache
	if s.cache != nil {
		s.cache.Delete(fmt.Sprintf("story:%s", storyID))
	}
}
