package db

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
)

// InitPostgres initializes a new PostgreSQL database connection.
func InitPostgres(connStr string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Ping the database to verify the connection
	err = db.Ping()
	if err != nil {
		db.Close() // Close the connection if ping fails
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")
	return db, nil
}
