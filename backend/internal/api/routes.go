package api

import (
	"database/sql"
	"net/http"
)

func NewRouter(db *sql.DB) http.Handler {
    mux := http.NewServeMux()

    // Category routes
    mux.HandleFunc("GET /categories", getCategoriesHandler(db))
    mux.HandleFunc("GET /categories/{name}", getCategoryByNameHandler(db))
    mux.HandleFunc("POST /categories", createCategoryHandler(db))
    mux.HandleFunc("PUT /categories/{name}", updateCategoryHandler(db))
    mux.HandleFunc("DELETE /categories/{name}", deleteCategoryHandler(db))

    return mux
}

