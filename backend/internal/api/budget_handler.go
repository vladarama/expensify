package api

import (
	"encoding/json"
	"expense-tracker/internal/models"
	"net/http"
	"strconv"

	"database/sql"
	"fmt"
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

// Get budget by ID
func getBudgetByIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid Budget ID", http.StatusBadRequest)
			return
		}
		budget, err := models.GetBudgetByID(db, id)
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

// Get a budget by category
func getBudgetByCategoryHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		category, err := strconv.ParseInt(r.PathValue("category"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}
		budget, err := models.GetBudgetsByCategoryID(db, category)
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
func updateBudgetHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the budget_id from the query parameters
		// budgetIDStr := r.URL.Query().Get("id") // Adjust based on your router setup
		// if budgetIDStr == "" {
		// 	http.Error(w, "Budget ID is required", http.StatusBadRequest)
		// 	return
		// }
		budgetIDStr := r.PathValue("id")
		// Convert the budgetID from string to int64
		budgetID, err := strconv.ParseInt(budgetIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid Budget ID", http.StatusBadRequest)
			return
		}

		// Fetch the current budget using the budget ID
		existingBudget, err := models.GetBudgetByID(db, budgetID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Budget not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Decode the new budget data from the request body
		var newBudget models.Budget
		err = json.NewDecoder(r.Body).Decode(&newBudget)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Merge the new values with the existing budget
		mergedBudget, err := mergeBudgets(existingBudget, newBudget)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest) // Handle validation errors
			return
		}

		// Update the budget in the database
		updatedBudget, err := models.UpdateBudget(db, mergedBudget)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with the updated budget
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedBudget)
	}
}

// // Update an existing budget by category
// func updateBudgetHandler(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// Get the category_id from the path
// 		categoryIDStr := r.PathValue("category")

// 		// Convert the categoryID from string to int64
// 		categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
// 		if err != nil {
// 			http.Error(w, "Invalid category ID", http.StatusBadRequest)
// 			return
// 		}

// 		var budget models.Budget
// 		err = json.NewDecoder(r.Body).Decode(&budget)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}
// 		// Set the converted categoryID to the budget
// 		budget.CategoryID = categoryID

// 		updatedBudget, err := models.UpdateBudget(db, budget)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(updatedBudget)
// 	}
// }

// Delete a budget by id
func deleteBudgetHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid Budget ID", http.StatusBadRequest)
			return
		}
		err = models.DeleteBudget(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Budget not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// func deleteBudgetHandler(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// Ensure the request method is DELETE
// 		if r.Method != http.MethodDelete {
// 			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
// 			return
// 		}

// 		// Extract the category from the path
// 		category := r.URL.Query().Get("category")
// 		if category == "" {
// 			http.Error(w, "Category is required", http.StatusBadRequest)
// 			return
// 		}

// 		// Attempt to delete the budget from the database
// 		err := models.DeleteBudget(db, category)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		// Respond with no content if successful
// 		w.WriteHeader(http.StatusNoContent)
// 	}
// }

func mergeBudgets(existing, new models.Budget) (models.Budget, error) {
	if new.Amount > 0 {
		existing.Amount = new.Amount
	}
	if new.Spent > 0 {
		existing.Spent = new.Spent
	}
	if !new.StartDate.IsZero() {
		existing.StartDate = new.StartDate
	}
	if !new.EndDate.IsZero() {
		existing.EndDate = new.EndDate
	}

	if new.CategoryID > 0 {
		existing.CategoryID = new.CategoryID
	}

	// Validate that end_date is after start_date
	if !existing.EndDate.IsZero() && !existing.StartDate.IsZero() && existing.EndDate.Before(existing.StartDate) {
		return existing, fmt.Errorf("end date must be after start date")
	}

	return existing, nil
}
