package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func InitPostgres() (*pgxpool.Pool, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error uis that %v", err)
		fmt.Println("No .env file found, proceeding with environment variables.")
	}

	connStr := os.Getenv("POSTGRES_CONNECTION")
	if connStr == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable not set")
	}

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("error creating connection pool: %w", err)
	}

	// Ping the database to verify the connection
	err = pool.Ping(context.Background())
	if err != nil {
		pool.Close() // Close the pool if ping fails
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")
	return pool, nil
}
