package api

import (
	"database/sql"
	"expense-tracker/internal/api/middleware"
	"net/http"
)

func NewRouter(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	// Category routes
	mux.HandleFunc("GET /categories", getCategoriesHandler(db))
	mux.HandleFunc("GET /categories/{id}", getCategoryByIDHandler(db))
	mux.HandleFunc("POST /categories", createCategoryHandler(db))
	mux.HandleFunc("PUT /categories/{id}", updateCategoryHandler(db))
	mux.HandleFunc("DELETE /categories/{id}", deleteCategoryHandler(db))

	// Income routes
	mux.HandleFunc("GET /incomes", getIncomesHandler(db))
	mux.HandleFunc("GET /incomes/{id}", getIncomeByIDHandler(db))
	mux.HandleFunc("POST /incomes", createIncomeHandler(db))
	mux.HandleFunc("PUT /incomes/{id}", updateIncomeHandler(db))
	mux.HandleFunc("DELETE /incomes/{id}", deleteIncomeHandler(db))

	// Expense routes
	mux.HandleFunc("GET /expenses", getExpensesHandler(db))
	mux.HandleFunc("GET /expenses/{id}", getExpenseByIDHandler(db))
	mux.HandleFunc("POST /expenses", createExpenseHandler(db))
	mux.HandleFunc("PUT /expenses/{id}", updateExpenseHandler(db))
	mux.HandleFunc("DELETE /expenses/{id}", deleteExpenseHandler(db))

	// Wrap the mux with CORS middleware
	return middleware.CORS(mux)
}

