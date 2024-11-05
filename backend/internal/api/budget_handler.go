package api

import (
	"encoding/json"
	"expense-tracker/internal/models"
	"net/http"
	"strconv"

	"database/sql"
)

// Get all budgets
func getBudgetsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		budgets, err := models.GetBudgets(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(budgets)
	}
}

// Get a budget by category
func getBudgetByCategoryHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		category := r.PathValue("category")
		budget, err := models.GetBudgetByCategory(db, category)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Budget not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(budget)
	}
}

// Create a new budget
func createBudgetHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var budget models.Budget
		err := json.NewDecoder(r.Body).Decode(&budget)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		createdBudget, err := models.CreateBudget(db, budget)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdBudget)
	}
}

// Update an existing budget

// Update an existing budget
func updateBudgetHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the category_id from the path
		categoryIDStr := r.PathValue("category")

		// Convert the categoryID from string to int64
		categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		var budget models.Budget
		err = json.NewDecoder(r.Body).Decode(&budget)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Set the converted categoryID to the budget
		budget.CategoryID = categoryID

		updatedBudget, err := models.UpdateBudget(db, budget)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedBudget)
	}
}

// Delete a budget by category
func deleteBudgetHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		category := r.PathValue("category")

		err := models.DeleteBudget(db, category)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
