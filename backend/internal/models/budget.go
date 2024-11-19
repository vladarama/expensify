package models

import (
	"database/sql"
	"errors"
	"strconv"
	"time"
)

// Budget represents the budget for a category.
type Budget struct {
	ID         int64     `json:"id"`
	CategoryID int64     `json:"category_id"`
	Amount     float64   `json:"amount"`
	Spent      float64   `json:"spent"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
}

// GetBudgets retrieves all budgets from the database.
func GetBudgets(db *sql.DB) ([]Budget, error) {
	rows, err := db.Query("SELECT id, category_id, amount, spent, start_date, end_date FROM Budget")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var budgets []Budget
	for rows.Next() {
		var budget Budget
		err := rows.Scan(&budget.ID, &budget.CategoryID, &budget.Amount, &budget.Spent, &budget.StartDate, &budget.EndDate)
		if err != nil {
			return nil, err
		}
		budgets = append(budgets, budget)
	}

	return budgets, nil
}

// GetBudgetByCategory retrieves a budget by category name.
func GetBudgetByCategory(db *sql.DB, category string) (Budget, error) {
	categoryID, err := strconv.ParseInt(category, 10, 64)
	if err != nil {
		return Budget{}, err
	}

	var budget Budget
	err = db.QueryRow("SELECT id, category_id, amount, spent, start_date, end_date FROM Budget WHERE category_id = $1",
		categoryID).Scan(&budget.ID, &budget.CategoryID, &budget.Amount, &budget.Spent, &budget.StartDate, &budget.EndDate)
	if err != nil {
		return Budget{}, err
	}
	return budget, nil
}

// CreateBudget adds a new budget to the database.
func CreateBudget(db *sql.DB, budget Budget) (Budget, error) {
	if budget.CategoryID == 0 {
		return Budget{}, errors.New("category cannot be empty")
	}
	if budget.Amount <= 0 {
		return Budget{}, errors.New("amount must be greater than zero")
	}
	if budget.StartDate.IsZero() {
		return Budget{}, errors.New("start date must be provided")
	}
	if budget.EndDate.IsZero() {
		return Budget{}, errors.New("end date must be provided")
	}
	if budget.EndDate.Before(budget.StartDate) {
		return Budget{}, errors.New("end date must be after start date")
	}

	var id int64
	err := db.QueryRow("INSERT INTO Budget (category_id, amount, spent, start_date, end_date) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		budget.CategoryID, budget.Amount, budget.Spent, budget.StartDate, budget.EndDate).Scan(&id)
	if err != nil {
		return Budget{}, err
	}
	budget.ID = id
	return budget, nil
}

// UpdateBudget updates an existing budget's information.
func UpdateBudget(db *sql.DB, budget Budget) (Budget, error) {
	if budget.CategoryID == 0 {
		return Budget{}, errors.New("category id must be provided")
	}

	if budget.Amount <= 0 {
		return Budget{}, errors.New("amount must be greater than zero")
	}
	if budget.StartDate.IsZero() {
		return Budget{}, errors.New("start date must be provided")
	}
	if budget.EndDate.IsZero() {
		return Budget{}, errors.New("end date must be provided")
	}
	if budget.EndDate.Before(budget.StartDate) {
		return Budget{}, errors.New("end date must be after start date")
	}

	_, err := db.Exec("UPDATE Budget SET amount = $1, spent = $2, start_date = $3, end_date = $4 WHERE category_id = $5",
		budget.Amount, budget.Spent, budget.StartDate, budget.EndDate, budget.CategoryID)
	if err != nil {
		return Budget{}, err
	}
	return budget, nil
}

// DeleteBudget removes a budget by category name.
func DeleteBudget(db *sql.DB, category string) error {
	categoryID, err := strconv.ParseInt(category, 10, 64)
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM Budget WHERE category_id = $1", categoryID)
	if err != nil {
		return err
	}
	return nil
}
