package api

import (
	"encoding/json"
	"expense-tracker/internal/models"
	"net/http"
	"strconv"

	"database/sql"
)

func getExpensesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		expenses, err := models.GetExpenses(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expenses)
	}
}

func getExpenseByIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid expense ID", http.StatusBadRequest)
			return
		}

		expense, err := models.GetExpenseByID(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Expense not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expense)
	}
}

func createExpenseHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var expense models.Expense
		err := json.NewDecoder(r.Body).Decode(&expense)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		createdExpense, err := models.CreateExpense(db, expense)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdExpense)
	}
}

func updateExpenseHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid expense ID", http.StatusBadRequest)
			return
		}

		var expense models.Expense
		err = json.NewDecoder(r.Body).Decode(&expense)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		expense.ID = id

		updatedExpense, err := models.UpdateExpense(db, expense)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedExpense)
	}
}

func deleteExpenseHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid expense ID", http.StatusBadRequest)
			return
		}

		err = models.DeleteExpense(db, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}