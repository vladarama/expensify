package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"expense-tracker/internal/api"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get database URL from environment variable
	dbUrl := os.Getenv("DATABASE_PUBLIC_URL")
	if dbUrl == "" {
		log.Fatal("DATABASE_PUBLIC_URL is not set")
	}

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	
	// Open connection to the database
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	// Ensure database connection is closed when the program exits
	defer db.Close()

	// Verify database connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database:", err)
	}

	// Initialize router with database connection
	router := api.NewRouter(db)

	// Start the HTTP server
    log.Printf("Server is running on port %s", port)
    if err := http.ListenAndServe(":"+port, router); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
