package postgresql

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // Import driver PostgreSQL
)

var database *sql.DB

// InitDB menginisialisasi koneksi database
func InitDB(host string, port int32, username string, password string, dbname string, sslmode string) error {
	var err error
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, username, password, dbname, sslmode,
	)

	database, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	// Cek koneksi database
	if err = database.Ping(); err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	log.Println("Connected to PostgreSQL successfully!")
	return nil
}

// GetDB mengembalikan instance database
func GetDB() *sql.DB {
	return database
}
