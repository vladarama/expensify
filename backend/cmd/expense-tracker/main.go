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

	// Initialize database (only run once)
	err = initializeDatabase(db)
	if err != nil {
		log.Fatal("Error initializing database:", err)
	}

	// Initialize router with database connection
	router := api.NewRouter(db)

	// Start the HTTP server
    log.Printf("Server is running on port %s", port)
    if err := http.ListenAndServe(":"+port, router); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}


func initializeDatabase(db *sql.DB) error {
    schema := `
    -- Table: Category
    CREATE TABLE IF NOT EXISTS Category (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        description TEXT
    );

    -- Table: Budget
    CREATE TABLE IF NOT EXISTS Budget (
        id SERIAL PRIMARY KEY,
        category_id INT REFERENCES Category(id) ON DELETE CASCADE,
        amount NUMERIC(10, 2) NOT NULL
    );

    -- Table: Expense
    CREATE TABLE IF NOT EXISTS Expense (
        id SERIAL PRIMARY KEY,
        category_id INT REFERENCES Category(id) ON DELETE SET NULL,
        amount NUMERIC(10, 2) NOT NULL,
        date DATE NOT NULL
    );

    -- Table: Income
    CREATE TABLE IF NOT EXISTS Income (
        id SERIAL PRIMARY KEY,
        amount NUMERIC(10, 2) NOT NULL,
        date DATE NOT NULL,
        source VARCHAR(255) NOT NULL
    );

    -- Table: Report
    CREATE TABLE IF NOT EXISTS Report (
        id SERIAL PRIMARY KEY,
        expense_id INT REFERENCES Expense(id) ON DELETE CASCADE,
        income_id INT REFERENCES Income(id) ON DELETE CASCADE
    );`

    // Execute the schema
    _, err := db.Exec(schema)
    return err
}