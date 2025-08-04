package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	pb "quiz-services/proto"

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
		dbPassword = "password"
	}
	if dbName == "" {
		dbName = "myapp"
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

func (db *Database) GetQuizzesBySection(sectionId string) ([]*pb.Quiz, error) {
	query := `
		SELECT id, section_id, question, options, correct_answer, explanation, difficulty, created_at
		FROM quizzes 
		WHERE section_id = $1
		ORDER BY created_at ASC
	`

	rows, err := db.conn.Query(query, sectionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quizzes []*pb.Quiz
	for rows.Next() {
		var quiz pb.Quiz
		var createdAt sql.NullTime
		var options []byte

		err := rows.Scan(
			&quiz.Id,
			&quiz.SectionId,
			&quiz.Question,
			&options,
			&quiz.CorrectAnswer,
			&quiz.Explanation,
			&quiz.Difficulty,
			&createdAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse options JSON
		if quiz.Options == nil {
			quiz.Options = make(map[string]string)
		}
		// You would need to implement JSON parsing here

		if createdAt.Valid {
			quiz.CreatedAt = timestamppb.New(createdAt.Time)
		}

		quizzes = append(quizzes, &quiz)
	}

	return quizzes, nil
}

func (db *Database) GetQuizByID(id string) (*pb.Quiz, error) {
	query := `
		SELECT id, section_id, question, options, correct_answer, explanation, difficulty, created_at
		FROM quizzes 
		WHERE id = $1
	`

	var quiz pb.Quiz
	var createdAt sql.NullTime
	var options []byte

	err := db.conn.QueryRow(query, id).Scan(
		&quiz.Id,
		&quiz.SectionId,
		&quiz.Question,
		&options,
		&quiz.CorrectAnswer,
		&quiz.Explanation,
		&quiz.Difficulty,
		&createdAt,
	)

	if err != nil {
		return nil, err
	}

	// Parse options JSON
	if quiz.Options == nil {
		quiz.Options = make(map[string]string)
	}
	// You would need to implement JSON parsing here

	if createdAt.Valid {
		quiz.CreatedAt = timestamppb.New(createdAt.Time)
	}

	return &quiz, nil
}

func (db *Database) CreateQuiz(quiz *pb.Quiz) (*pb.Quiz, error) {
	query := `
		INSERT INTO quizzes (id, section_id, question, options, correct_answer, explanation, difficulty, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, section_id, question, options, correct_answer, explanation, difficulty, created_at
	`

	// Convert options map to JSON
	optionsJSON := "{}" // You would need to implement JSON marshaling here

	var createdQuiz pb.Quiz
	var createdAt sql.NullTime
	var options []byte

	err := db.conn.QueryRow(query,
		quiz.Id,
		quiz.SectionId,
		quiz.Question,
		optionsJSON,
		quiz.CorrectAnswer,
		quiz.Explanation,
		quiz.Difficulty,
		time.Now(),
	).Scan(
		&createdQuiz.Id,
		&createdQuiz.SectionId,
		&createdQuiz.Question,
		&options,
		&createdQuiz.CorrectAnswer,
		&createdQuiz.Explanation,
		&createdQuiz.Difficulty,
		&createdAt,
	)

	if err != nil {
		return nil, err
	}

	// Parse options JSON
	if createdQuiz.Options == nil {
		createdQuiz.Options = make(map[string]string)
	}

	if createdAt.Valid {
		createdQuiz.CreatedAt = timestamppb.New(createdAt.Time)
	}

	return &createdQuiz, nil
}

func (db *Database) UpdateQuiz(id string, quiz *pb.Quiz) (*pb.Quiz, error) {
	query := `
		UPDATE quizzes 
		SET question = $2, options = $3, correct_answer = $4, explanation = $5, difficulty = $6
		WHERE id = $1
		RETURNING id, section_id, question, options, correct_answer, explanation, difficulty, created_at
	`

	// Convert options map to JSON
	optionsJSON := "{}" // You would need to implement JSON marshaling here

	var updatedQuiz pb.Quiz
	var createdAt sql.NullTime
	var options []byte

	err := db.conn.QueryRow(query,
		id,
		quiz.Question,
		optionsJSON,
		quiz.CorrectAnswer,
		quiz.Explanation,
		quiz.Difficulty,
	).Scan(
		&updatedQuiz.Id,
		&updatedQuiz.SectionId,
		&updatedQuiz.Question,
		&options,
		&updatedQuiz.CorrectAnswer,
		&updatedQuiz.Explanation,
		&updatedQuiz.Difficulty,
		&createdAt,
	)

	if err != nil {
		return nil, err
	}

	// Parse options JSON
	if updatedQuiz.Options == nil {
		updatedQuiz.Options = make(map[string]string)
	}

	if createdAt.Valid {
		updatedQuiz.CreatedAt = timestamppb.New(createdAt.Time)
	}

	return &updatedQuiz, nil
}

func (db *Database) DeleteQuiz(id string) error {
	query := `DELETE FROM quizzes WHERE id = $1`

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

func (db *Database) SubmitQuizAnswer(userId, quizId, selectedOption string) (*pb.QuizResponse, error) {
	// First, get the correct answer
	var correctAnswer string
	err := db.conn.QueryRow("SELECT correct_answer FROM quizzes WHERE id = $1", quizId).Scan(&correctAnswer)
	if err != nil {
		return nil, err
	}

	isCorrect := selectedOption == correctAnswer

	// Insert the response
	response := &pb.QuizResponse{
		Id:             uuid.New().String(),
		UserId:         userId,
		QuizId:         quizId,
		SelectedOption: selectedOption,
		IsCorrect:      isCorrect,
		AnsweredAt:     timestamppb.Now(),
	}

	query := `
		INSERT INTO quiz_responses (id, user_id, quiz_id, selected_option, is_correct, answered_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = db.conn.Exec(query,
		response.Id,
		response.UserId,
		response.QuizId,
		response.SelectedOption,
		response.IsCorrect,
		time.Now(),
	)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (db *Database) GetUserQuizHistory(userId string) ([]*pb.QuizResponse, error) {
	query := `
		SELECT id, user_id, quiz_id, selected_option, is_correct, answered_at
		FROM quiz_responses 
		WHERE user_id = $1
		ORDER BY answered_at DESC
	`

	rows, err := db.conn.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var responses []*pb.QuizResponse
	for rows.Next() {
		var response pb.QuizResponse
		var answeredAt sql.NullTime

		err := rows.Scan(
			&response.Id,
			&response.UserId,
			&response.QuizId,
			&response.SelectedOption,
			&response.IsCorrect,
			&answeredAt,
		)
		if err != nil {
			return nil, err
		}

		if answeredAt.Valid {
			response.AnsweredAt = timestamppb.New(answeredAt.Time)
		}

		responses = append(responses, &response)
	}

	return responses, nil
}

func (db *Database) GetQuizStatistics(quizId string) (*pb.QuizStatistics, error) {
	query := `
		SELECT 
			COUNT(*) as total_attempts,
			COUNT(CASE WHEN is_correct = true THEN 1 END) as correct_attempts
		FROM quiz_responses 
		WHERE quiz_id = $1
	`

	var stats pb.QuizStatistics
	var totalAttempts, correctAttempts int32

	err := db.conn.QueryRow(query, quizId).Scan(&totalAttempts, &correctAttempts)
	if err != nil {
		return nil, err
	}

	stats.QuizId = quizId
	stats.TotalAttempts = totalAttempts
	stats.CorrectAttempts = correctAttempts

	if totalAttempts > 0 {
		stats.SuccessRate = float32(correctAttempts) / float32(totalAttempts) * 100
	}

	return &stats, nil
}

// Helper function to convert options map to JSON
func optionsToJSON(options map[string]string) (string, error) {
	if options == nil {
		return "{}", nil
	}
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return "{}", err
	}
	return string(jsonBytes), nil
}

// Helper function to convert JSON to options map
func jsonToOptions(jsonStr string) (map[string]string, error) {
	var options map[string]string
	err := json.Unmarshal([]byte(jsonStr), &options)
	if err != nil {
		return make(map[string]string), err
	}
	return options, nil
}
