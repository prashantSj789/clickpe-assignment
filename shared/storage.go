package shared

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func init() {
	// Try to load .env file if running locally
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, assuming environment variables are set externally")
	}
}

func ConnectDB() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	return db, nil
}


func InsertUser(db *sql.DB, user User) error {
	query := `
	INSERT INTO users (id, name, email, monthly_income, credit_score, employment_status, age)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`
	
	_, err := db.Exec(query,
		user.UserID,
		user.Name,
		user.Email,
		user.MonthlyIncome,
		user.CreditScore,
		user.EmploymentStatus,
		user.Age,
	)
	
	return err
}





