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

	// Income routes
	mux.HandleFunc("GET /incomes", getIncomesHandler(db))
	mux.HandleFunc("GET /incomes/{id}", getIncomeByIDHandler(db))
	mux.HandleFunc("POST /incomes", createIncomeHandler(db))
	mux.HandleFunc("PUT /incomes/{id}", updateIncomeHandler(db))
	mux.HandleFunc("DELETE /incomes/{id}", deleteIncomeHandler(db))

	return mux
}

