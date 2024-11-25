package models

import (
	"database/sql"
	"errors"
	"fmt"
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

// // GetBudgetByCategory retrieves a budget by category name.
// func GetBudgetByCategory(db *sql.DB, category string) (Budget, error) {
// 	categoryID, err := strconv.ParseInt(category, 10, 64)
// 	if err != nil {
// 		return Budget{}, err
// 	}

// 	var budget Budget
// 	err = db.QueryRow("SELECT id, category_id, amount, spent, start_date, end_date FROM Budget WHERE category_id = $1",
// 		categoryID).Scan(&budget.ID, &budget.CategoryID, &budget.Amount, &budget.Spent, &budget.StartDate, &budget.EndDate)
// 	if err != nil {
// 		return Budget{}, err
// 	}
// 	return budget, nil
// }

func GetBudgetsByCategoryName(db *sql.DB, categoryName string) ([]Budget, error) {
	// Retrieve the category ID using the category name
	var categoryID int64
	err := db.QueryRow("SELECT id FROM Category WHERE name = $1", categoryName).Scan(&categoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("category not found: %s", categoryName)
		}
		return nil, fmt.Errorf("failed to retrieve category ID: %w", err)
	}

	// Retrieve all budgets associated with the category ID
	rows, err := db.Query(`
		SELECT id, category_id, amount, spent, start_date, end_date 
		FROM Budget 
		WHERE category_id = $1
	`, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve budgets: %w", err)
	}
	defer rows.Close()

	// Collect all budgets into a slice
	var budgets []Budget
	for rows.Next() {
		var budget Budget
		err := rows.Scan(&budget.ID, &budget.CategoryID, &budget.Amount, &budget.Spent, &budget.StartDate, &budget.EndDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan budget: %w", err)
		}
		budgets = append(budgets, budget)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over budgets: %w", err)
	}

	return budgets, nil
}

// GetBudgetByID retrieves a budget by ID.
func GetBudgetByID(db *sql.DB, id int64) (Budget, error) {
	var budget Budget
	err := db.QueryRow("SELECT id, category_id, amount, spent, start_date, end_date FROM Budget WHERE id = $1", id).Scan(&budget.ID, &budget.CategoryID, &budget.Amount, &budget.Spent, &budget.StartDate, &budget.EndDate)
	if err != nil {
		return Budget{}, err
	}
	return budget, nil
}

// GetBudgetByCategory retrieves a budget by category id.
func GetBudgetsByCategoryID(db *sql.DB, categoryID int64) ([]Budget, error) {
	// Retrieve all budgets associated with the category ID
	rows, err := db.Query(`
		SELECT id, category_id, amount, spent, start_date, end_date 
		FROM Budget 
		WHERE category_id = $1
	`, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve budgets: %w", err)
	}
	defer rows.Close()

	// Collect all budgets into a slice
	var budgets []Budget
	for rows.Next() {
		var budget Budget
		err := rows.Scan(&budget.ID, &budget.CategoryID, &budget.Amount, &budget.Spent, &budget.StartDate, &budget.EndDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan budget: %w", err)
		}
		budgets = append(budgets, budget)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over budgets: %w", err)
	}

	return budgets, nil
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

	// Check for overlapping budgets
	overlap, err := DoesBudgetOverlap(db, budget.CategoryID, budget.StartDate, budget.EndDate, 0)
	if err != nil {
		return Budget{}, fmt.Errorf("failed to validate budget overlap: %w", err)
	}
	if overlap {
		return Budget{}, errors.New("budget dates overlap with an existing budget")
	}

	// Calculate the total spent for the category and date range
	totalSpent, err := CalculateTotalSpent(db, budget.CategoryID, budget.StartDate, budget.EndDate)
	if err != nil {
		return Budget{}, err
	}
	budget.Spent = totalSpent

	// Set the calculated spent amount
	budget.Spent = totalSpent

	var id int64
	err = db.QueryRow("INSERT INTO Budget (category_id, amount, spent, start_date, end_date) VALUES ($1, $2, $3, $4, $5) RETURNING id",
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

	// Check for overlapping budgets
	overlap, err := DoesBudgetOverlap(db, budget.CategoryID, budget.StartDate, budget.EndDate, budget.ID)
	if err != nil {
		return Budget{}, fmt.Errorf("failed to validate budget overlap: %w", err)
	}
	if overlap {
		return Budget{}, errors.New("budget dates overlap with an existing budget")
	}

	// Calculate the total spent for the category and date range
	totalSpent, err := CalculateTotalSpent(db, budget.CategoryID, budget.StartDate, budget.EndDate)
	if err != nil {
		return Budget{}, err
	}
	budget.Spent = totalSpent

	_, err = db.Exec("UPDATE Budget SET amount = $1, spent = $2, start_date = $3, end_date = $4 WHERE category_id = $5",
		budget.Amount, budget.Spent, budget.StartDate, budget.EndDate, budget.CategoryID)
	if err != nil {
		return Budget{}, err
	}
	return budget, nil
}

// DeleteBudget removes a budget by category name.
func DeleteBudget(db *sql.DB, id int64) error {
	result, err := db.Exec("DELETE FROM Budget WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func DoesBudgetOverlap(db *sql.DB, categoryID int64, startDate, endDate time.Time, excludeBudgetID int64) (bool, error) {
	// Check if the budget overlaps with any existing budget other than the one being updated
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM Budget
			WHERE category_id = $1
			AND id <> $4
			AND (
				(start_date <= $3 AND end_date >= $2)
			)
		)
	`
	var exists bool
	err := db.QueryRow(query, categoryID, startDate, endDate, excludeBudgetID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check budget overlap: %w", err)
	}
	return exists, nil
}

func CalculateTotalSpent(db *sql.DB, categoryID int64, startDate, endDate time.Time) (float64, error) {
	var totalSpent float64
	err := db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0)
		FROM Expense
		WHERE category_id = $1 AND date >= $2 AND date <= $3
	`, categoryID, startDate, endDate).Scan(&totalSpent)

	if err != nil {
		return 0, fmt.Errorf("failed to calculate spent amount: %w", err)
	}
	return totalSpent, nil
}
