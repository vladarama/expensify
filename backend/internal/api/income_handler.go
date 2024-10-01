package api

import (
	"encoding/json"
	"expense-tracker/internal/models"
	"net/http"
	"strconv"

	"database/sql"
)

func getIncomesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		incomes, err := models.GetIncomes(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(incomes)
	}
}

func getIncomeByIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid income ID", http.StatusBadRequest)
			return
		}

		income, err := models.GetIncomeByID(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Income not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(income)
	}
}

func createIncomeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var income models.Income
		err := json.NewDecoder(r.Body).Decode(&income)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		createdIncome, err := models.CreateIncome(db, income)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdIncome)
	}
}

func updateIncomeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid income ID", http.StatusBadRequest)
			return
		}

		var income models.Income
		err = json.NewDecoder(r.Body).Decode(&income)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		income.ID = id

		updatedIncome, err := models.UpdateIncome(db, income)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedIncome)
	}
}

func deleteIncomeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid income ID", http.StatusBadRequest)
			return
		}

		err = models.DeleteIncome(db, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}