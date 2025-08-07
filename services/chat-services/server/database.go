package server

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

func (d *Database) QueryRow(checkQuery string, id string) {
	panic("unimplemented")
}

func (d *Database) Exec(query string, notificationID string, param3 string, s string, title string, string sql.NullString, id string, param8 sql.NullString, param9 sql.NullString, param10 sql.NullString, read bool, time time.Time) (any, error) {
	panic("unimplemented")
}

// NewDatabase creates a new database connection
func NewDatabase() *Database {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "2521")
	dbname := getEnv("DB_NAME", "userdb")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to database")

	return &Database{DB: db}
}

func (d *Database) Close() {
	if d.DB != nil {
		d.DB.Close()
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
