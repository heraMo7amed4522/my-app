package server

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	pb "feedback-services/proto"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Database struct {
	conn *sql.DB
}

func NewDatabase() *Database {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "5432"
	}
	if dbUser == "" {
		dbUser = "postgres"
	}
	if dbPassword == "" {
		dbPassword = "2521"
	}
	if dbName == "" {
		dbName = "userdb"
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to the database")
	return &Database{conn: db}
}

// SubmitFeedback creates a new feedback entry
func (db *Database) SubmitFeedback(userID, templateID string, rating float32, comment, language string) (*pb.Feedback, error) {
	id := uuid.New().String()
	createdAt := time.Now()

	query := `
		INSERT INTO template_feedback (id, user_id, template_id, rating, comment, language, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, template_id, rating, comment, language, created_at
	`

	var feedback pb.Feedback
	var dbCreatedAt time.Time

	err := db.conn.QueryRow(query, id, userID, templateID, rating, comment, language, createdAt).Scan(
		&feedback.Id,
		&feedback.UserId,
		&feedback.TemplateId,
		&feedback.Rating,
		&feedback.Comment,
		&feedback.Language,
		&dbCreatedAt,
	)

	if err != nil {
		return nil, err
	}

	feedback.CreatedAt = timestamppb.New(dbCreatedAt)
	return &feedback, nil
}

// GetTemplateFeedback retrieves feedback for a specific template with pagination
func (db *Database) GetTemplateFeedback(templateID string, page, limit int32) ([]*pb.Feedback, int32, error) {
	offset := (page - 1) * limit

	// Get total count
	countQuery := "SELECT COUNT(*) FROM template_feedback WHERE template_id = $1"
	var totalCount int32
	err := db.conn.QueryRow(countQuery, templateID).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Get feedback with pagination
	query := `
		SELECT id, user_id, template_id, rating, comment, language, created_at
		FROM template_feedback 
		WHERE template_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := db.conn.Query(query, templateID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var feedbacks []*pb.Feedback
	for rows.Next() {
		var feedback pb.Feedback
		var createdAt time.Time

		err := rows.Scan(
			&feedback.Id,
			&feedback.UserId,
			&feedback.TemplateId,
			&feedback.Rating,
			&feedback.Comment,
			&feedback.Language,
			&createdAt,
		)
		if err != nil {
			return nil, 0, err
		}

		feedback.CreatedAt = timestamppb.New(createdAt)
		feedbacks = append(feedbacks, &feedback)
	}

	return feedbacks, totalCount, nil
}

// GetFeedbackStatistics calculates statistics for a template
func (db *Database) GetFeedbackStatistics(templateID string) (*pb.FeedbackStatistics, error) {
	query := `
		SELECT 
			COUNT(*) as total_ratings,
			AVG(rating) as average_rating,
			COUNT(CASE WHEN rating >= 1 AND rating < 2 THEN 1 END) as rating_1,
			COUNT(CASE WHEN rating >= 2 AND rating < 3 THEN 1 END) as rating_2,
			COUNT(CASE WHEN rating >= 3 AND rating < 4 THEN 1 END) as rating_3,
			COUNT(CASE WHEN rating >= 4 AND rating < 5 THEN 1 END) as rating_4,
			COUNT(CASE WHEN rating = 5 THEN 1 END) as rating_5
		FROM template_feedback 
		WHERE template_id = $1
	`

	var totalRatings int32
	var averageRating float64
	var rating1, rating2, rating3, rating4, rating5 int32

	err := db.conn.QueryRow(query, templateID).Scan(
		&totalRatings,
		&averageRating,
		&rating1,
		&rating2,
		&rating3,
		&rating4,
		&rating5,
	)
	if err != nil {
		return nil, err
	}

	ratingDistribution := map[string]int32{
		"1": rating1,
		"2": rating2,
		"3": rating3,
		"4": rating4,
		"5": rating5,
	}

	return &pb.FeedbackStatistics{
		TemplateId:         templateID,
		AverageRating:      float32(averageRating),
		TotalRatings:       totalRatings,
		RatingDistribution: ratingDistribution,
	}, nil
}

// UpdateFeedback updates an existing feedback entry
func (db *Database) UpdateFeedback(id string, rating float32, comment string) (*pb.Feedback, error) {
	query := `
		UPDATE template_feedback 
		SET rating = $2, comment = $3
		WHERE id = $1
		RETURNING id, user_id, template_id, rating, comment, language, created_at
	`

	var feedback pb.Feedback
	var createdAt time.Time

	err := db.conn.QueryRow(query, id, rating, comment).Scan(
		&feedback.Id,
		&feedback.UserId,
		&feedback.TemplateId,
		&feedback.Rating,
		&feedback.Comment,
		&feedback.Language,
		&createdAt,
	)

	if err != nil {
		return nil, err
	}

	feedback.CreatedAt = timestamppb.New(createdAt)
	return &feedback, nil
}

// DeleteFeedback removes a feedback entry
func (db *Database) DeleteFeedback(id string) error {
	query := "DELETE FROM template_feedback WHERE id = $1"
	result, err := db.conn.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
