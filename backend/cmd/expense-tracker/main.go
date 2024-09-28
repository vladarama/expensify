package main

import (
    "database/sql"
    "log"
    "net/http"
    "os"

    _ "github.com/lib/pq"
    "github.com/joho/godotenv"

    "expense-tracker/internal/api"
)

func main() {
    // Load environment variables from .env file
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found or failed to load")
    }

    // Load environment variables
    dbURL := os.Getenv("DATABASE_PUBLIC_URL")
    if dbURL == "" {
        log.Fatal("DATABASE_URL environment variable is required")
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // default port if not specified
    }

    jwtSecret := os.Getenv("JWT_SECRET")
    if jwtSecret == "" {
        log.Fatal("JWT_SECRET environment variable is required")
    }

    // Connect to the database
    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    // Test the database connection
    err = db.Ping()
    if err != nil {
        log.Fatalf("Failed to ping database: %v", err)
    }
    log.Println("Successfully connected to the database")


    // Set up router
    router := api.NewRouter(db, jwtSecret)

    // Start the server
    log.Printf("Server is running on port %s", port)
    if err := http.ListenAndServe(":"+port, router); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}

// runMigrations runs the SQL scripts in the migrations directory
func runMigrations(db *sql.DB) error {
    files, err := os.ReadDir("migrations")
    if err != nil {
        return err
    }

    for _, file := range files {
        if file.IsDir() {
            continue
        }

        content, err := os.ReadFile("migrations/" + file.Name())
        if err != nil {
            return err
        }

        _, err = db.Exec(string(content))
        if err != nil {
            return err
        }
    }

    return nil
}
