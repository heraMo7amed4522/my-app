package server

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	pb "user-progress-services/proto"

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

	log.Println("Successfully connected to database")

	return &Database{conn: db}
}

// GetUserProgress retrieves user progress for a specific template or all templates
func (db *Database) GetUserProgress(userID, templateID string) ([]*pb.UserProgress, error) {
	var query string
	var args []interface{}

	if templateID != "" {
		query = `SELECT id, user_id, template_id, section_id, progress, completed, last_viewed, created_at 
				 FROM user_progress WHERE user_id = $1 AND template_id = $2 ORDER BY created_at DESC`
		args = []interface{}{userID, templateID}
	} else {
		query = `SELECT id, user_id, template_id, section_id, progress, completed, last_viewed, created_at 
				 FROM user_progress WHERE user_id = $1 ORDER BY created_at DESC`
		args = []interface{}{userID}
	}

	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var progressList []*pb.UserProgress
	for rows.Next() {
		var progress pb.UserProgress
		var lastViewed, createdAt time.Time

		err := rows.Scan(
			&progress.Id,
			&progress.UserId,
			&progress.TemplateId,
			&progress.SectionId,
			&progress.Progress,
			&progress.Completed,
			&lastViewed,
			&createdAt,
		)
		if err != nil {
			return nil, err
		}

		progress.LastViewed = timestamppb.New(lastViewed)
		progress.CreatedAt = timestamppb.New(createdAt)
		progressList = append(progressList, &progress)
	}

	return progressList, nil
}

// UpdateProgress updates or creates user progress
func (db *Database) UpdateProgress(userID, templateID, sectionID string, progress float32, completed bool) (*pb.UserProgress, error) {
	// Check if progress already exists
	existingProgress, err := db.getProgressByUserAndTemplate(userID, templateID, sectionID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if existingProgress != nil {
		// Update existing progress
		query := `UPDATE user_progress SET progress = $1, completed = $2, last_viewed = $3 
				  WHERE user_id = $4 AND template_id = $5 AND section_id = $6 RETURNING id, created_at`
		
		var id string
		var createdAt time.Time
		now := time.Now()
		
		err = db.conn.QueryRow(query, progress, completed, now, userID, templateID, sectionID).Scan(&id, &createdAt)
		if err != nil {
			return nil, err
		}

		return &pb.UserProgress{
			Id:         id,
			UserId:     userID,
			TemplateId: templateID,
			SectionId:  sectionID,
			Progress:   progress,
			Completed:  completed,
			LastViewed: timestamppb.New(now),
			CreatedAt:  timestamppb.New(createdAt),
		}, nil
	} else {
		// Create new progress
		id := uuid.New().String()
		now := time.Now()
		
		query := `INSERT INTO user_progress (id, user_id, template_id, section_id, progress, completed, last_viewed, created_at) 
				  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
		
		_, err = db.conn.Exec(query, id, userID, templateID, sectionID, progress, completed, now, now)
		if err != nil {
			return nil, err
		}

		return &pb.UserProgress{
			Id:         id,
			UserId:     userID,
			TemplateId: templateID,
			SectionId:  sectionID,
			Progress:   progress,
			Completed:  completed,
			LastViewed: timestamppb.New(now),
			CreatedAt:  timestamppb.New(now),
		}, nil
	}
}

// GetCompletedTemplates retrieves completed templates for a user
func (db *Database) GetCompletedTemplates(userID string, page, limit int32) ([]*pb.UserProgress, int32, error) {
	offset := (page - 1) * limit
	
	query := `SELECT id, user_id, template_id, section_id, progress, completed, last_viewed, created_at 
			  FROM user_progress WHERE user_id = $1 AND completed = true 
			  ORDER BY last_viewed DESC LIMIT $2 OFFSET $3`
	
	rows, err := db.conn.Query(query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var progressList []*pb.UserProgress
	for rows.Next() {
		var progress pb.UserProgress
		var lastViewed, createdAt time.Time

		err := rows.Scan(
			&progress.Id,
			&progress.UserId,
			&progress.TemplateId,
			&progress.SectionId,
			&progress.Progress,
			&progress.Completed,
			&lastViewed,
			&createdAt,
		)
		if err != nil {
			return nil, 0, err
		}

		progress.LastViewed = timestamppb.New(lastViewed)
		progress.CreatedAt = timestamppb.New(createdAt)
		progressList = append(progressList, &progress)
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM user_progress WHERE user_id = $1 AND completed = true`
	var totalCount int32
	err = db.conn.QueryRow(countQuery, userID).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	return progressList, totalCount, nil
}

// GetUserAchievements retrieves user achievements
func (db *Database) GetUserAchievements(userID string) ([]*pb.UserAchievement, error) {
	query := `SELECT ua.id, ua.user_id, ua.achievement_id, ua.unlocked_at,
					 a.id, a.code, a.name, a.description, a.badge_url, a.created_at
			  FROM user_achievements ua
			  JOIN achievements a ON ua.achievement_id = a.id
			  WHERE ua.user_id = $1 ORDER BY ua.unlocked_at DESC`
	
	rows, err := db.conn.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var achievements []*pb.UserAchievement
	for rows.Next() {
		var userAchievement pb.UserAchievement
		var achievement pb.Achievement
		var unlockedAt, achievementCreatedAt time.Time

		err := rows.Scan(
			&userAchievement.Id,
			&userAchievement.UserId,
			&userAchievement.AchievementId,
			&unlockedAt,
			&achievement.Id,
			&achievement.Code,
			&achievement.Name,
			&achievement.Description,
			&achievement.BadgeUrl,
			&achievementCreatedAt,
		)
		if err != nil {
			return nil, err
		}

		achievement.CreatedAt = timestamppb.New(achievementCreatedAt)
		userAchievement.Achievement = &achievement
		userAchievement.UnlockedAt = timestamppb.New(unlockedAt)
		achievements = append(achievements, &userAchievement)
	}

	return achievements, nil
}

// UnlockAchievement unlocks an achievement for a user
func (db *Database) UnlockAchievement(userID, achievementID string) (*pb.UserAchievement, error) {
	// Check if achievement already unlocked
	existingQuery := `SELECT id FROM user_achievements WHERE user_id = $1 AND achievement_id = $2`
	var existingID string
	err := db.conn.QueryRow(existingQuery, userID, achievementID).Scan(&existingID)
	if err == nil {
		return nil, fmt.Errorf("achievement already unlocked")
	}
	if err != sql.ErrNoRows {
		return nil, err
	}

	// Get achievement details
	achievementQuery := `SELECT id, code, name, description, badge_url, created_at FROM achievements WHERE id = $1`
	var achievement pb.Achievement
	var achievementCreatedAt time.Time
	
	err = db.conn.QueryRow(achievementQuery, achievementID).Scan(
		&achievement.Id,
		&achievement.Code,
		&achievement.Name,
		&achievement.Description,
		&achievement.BadgeUrl,
		&achievementCreatedAt,
	)
	if err != nil {
		return nil, err
	}
	achievement.CreatedAt = timestamppb.New(achievementCreatedAt)

	// Unlock achievement
	id := uuid.New().String()
	now := time.Now()
	
	insertQuery := `INSERT INTO user_achievements (id, user_id, achievement_id, unlocked_at) VALUES ($1, $2, $3, $4)`
	_, err = db.conn.Exec(insertQuery, id, userID, achievementID, now)
	if err != nil {
		return nil, err
	}

	return &pb.UserAchievement{
		Id:            id,
		UserId:        userID,
		AchievementId: achievementID,
		Achievement:   &achievement,
		UnlockedAt:    timestamppb.New(now),
	}, nil
}

// GetLearningPath generates a learning path for a user
func (db *Database) GetLearningPath(userID string) (*pb.LearningPath, error) {
	// Get completed templates count
	completedQuery := `SELECT COUNT(DISTINCT template_id) FROM user_progress WHERE user_id = $1 AND completed = true`
	var completedCount int32
	err := db.conn.QueryRow(completedQuery, userID).Scan(&completedCount)
	if err != nil {
		return nil, err
	}

	// Get total templates count (this would typically come from template service)
	totalQuery := `SELECT COUNT(*) FROM history_templates`
	var totalCount int32
	err = db.conn.QueryRow(totalQuery).Scan(&totalCount)
	if err != nil {
		// If history_templates table doesn't exist, use a default
		totalCount = 100
	}

	// Calculate overall progress
	var overallProgress float32
	if totalCount > 0 {
		overallProgress = float32(completedCount) / float32(totalCount) * 100
	}

	// Get recommended templates (templates not yet completed)
	recommendedQuery := `SELECT DISTINCT ht.id FROM history_templates ht 
						  LEFT JOIN user_progress up ON ht.id = up.template_id AND up.user_id = $1 AND up.completed = true
						  WHERE up.template_id IS NULL 
						  ORDER BY ht.difficulty, ht.created_at LIMIT 10`
	
	rows, err := db.conn.Query(recommendedQuery, userID)
	if err != nil {
		// If query fails, return empty recommendations
		return &pb.LearningPath{
			UserId:                   userID,
			RecommendedTemplateIds:   []string{},
			OverallProgress:          overallProgress,
			CompletedTemplates:       completedCount,
			TotalTemplates:           totalCount,
		}, nil
	}
	defer rows.Close()

	var recommendedIDs []string
	for rows.Next() {
		var templateID string
		err := rows.Scan(&templateID)
		if err != nil {
			continue
		}
		recommendedIDs = append(recommendedIDs, templateID)
	}

	return &pb.LearningPath{
		UserId:                 userID,
		RecommendedTemplateIds: recommendedIDs,
		OverallProgress:        overallProgress,
		CompletedTemplates:     completedCount,
		TotalTemplates:         totalCount,
	}, nil
}

// Helper function to get existing progress
func (db *Database) getProgressByUserAndTemplate(userID, templateID, sectionID string) (*pb.UserProgress, error) {
	query := `SELECT id, user_id, template_id, section_id, progress, completed, last_viewed, created_at 
			  FROM user_progress WHERE user_id = $1 AND template_id = $2 AND section_id = $3`
	
	var progress pb.UserProgress
	var lastViewed, createdAt time.Time
	
	err := db.conn.QueryRow(query, userID, templateID, sectionID).Scan(
		&progress.Id,
		&progress.UserId,
		&progress.TemplateId,
		&progress.SectionId,
		&progress.Progress,
		&progress.Completed,
		&lastViewed,
		&createdAt,
	)
	if err != nil {
		return nil, err
	}

	progress.LastViewed = timestamppb.New(lastViewed)
	progress.CreatedAt = timestamppb.New(createdAt)
	return &progress, nil
}