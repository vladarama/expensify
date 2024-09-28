package api

import (
    "database/sql"
    "net/http"

    "expense-tracker/internal/auth"
)

type Server struct {
    DB        *sql.DB
    JWTSecret string
}

func NewRouter(db *sql.DB, jwtSecret string) http.Handler {
    server := &Server{DB: db, JWTSecret: jwtSecret}
    mux := http.NewServeMux()

    // Public routes
    mux.HandleFunc("/api/auth/register", server.Register)
    mux.HandleFunc("/api/auth/login", server.Login)

    // Protected routes
    authMiddleware := auth.JWTMiddleware(jwtSecret)
    mux.Handle("/api/expenses", authMiddleware(http.HandlerFunc(server.ExpensesHandler)))
    mux.Handle("/api/expenses/", authMiddleware(http.HandlerFunc(server.ExpensesHandler)))

    mux.Handle("/api/categories", authMiddleware(http.HandlerFunc(server.CategoriesHandler)))
    mux.Handle("/api/categories/", authMiddleware(http.HandlerFunc(server.CategoriesHandler)))

    return mux
}
