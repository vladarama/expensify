package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Expense struct {
	ID          int64     `json:"id"`
	CategoryID  int64     `json:"category_id"`
	Amount      float64   `json:"amount"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
}

func GetExpenses(db *sql.DB) ([]Expense, error) {
	rows, err := db.Query("SELECT id, category_id, amount, date, description FROM Expense")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []Expense
	for rows.Next() {
		var expense Expense
		if err := rows.Scan(&expense.ID, &expense.CategoryID, &expense.Amount, &expense.Date, &expense.Description); err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}

	return expenses, nil
}

func GetExpenseByID(db *sql.DB, id int64) (Expense, error) {
	var expense Expense
	err := db.QueryRow("SELECT id, category_id, amount, date, description FROM Expense WHERE id = $1", id).
		Scan(&expense.ID, &expense.CategoryID, &expense.Amount, &expense.Date, &expense.Description)
	if err != nil {
		return Expense{}, err
	}
	return expense, nil
}

// validateCreateExpense validates fields for creating an expense.
func validateCreateExpense(expense Expense) error {
	if expense.Description == "" {
		return errors.New("description is required")
	}
	if expense.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	if expense.Date.IsZero() || expense.Date.After(time.Now()) {
		return errors.New("date must be provided and cannot be in the future")
	}
	if expense.CategoryID <= 0 {
		return errors.New("category ID must be provided")
	}
	return nil
}

// CreateExpense adds a new expense to the database and updates the associated budget.
func CreateExpense(db *sql.DB, expense Expense) (Expense, error) {
	if err := validateCreateExpense(expense); err != nil {
		return Expense{}, err
	}

	// Insert the new expense and get the ID using RETURNING
	err := db.QueryRow(
		"INSERT INTO Expense (category_id, amount, date, description) VALUES ($1, $2, $3, $4) RETURNING id, category_id, amount, date, description",
		expense.CategoryID, expense.Amount, expense.Date, expense.Description,
	).Scan(&expense.ID, &expense.CategoryID, &expense.Amount, &expense.Date, &expense.Description)
	if err != nil {
		return Expense{}, err
	}

	// Check if budget exists before updating it
	budgetExists, err := checkBudgetExists(db, expense.CategoryID)
	if err != nil {
		return Expense{}, err
	}

	if budgetExists {
		if err := updateBudgetSpent(db, expense.CategoryID, expense.Amount); err != nil {
			return Expense{}, err
		}
	}

	return expense, nil
}

// validateUpdateExpense validates fields for updating an expense.
func validateUpdateExpense(expense Expense, existingExpense Expense) error {
	if expense.Amount <= 0 && expense.Amount != existingExpense.Amount {
		return errors.New("amount must be greater than zero if provided")
	}
	if !expense.Date.IsZero() && (expense.Date.Before(existingExpense.Date) || expense.Date.After(time.Now())) {
		return errors.New("date must be within the budget period and cannot be in the future")
	}
	return nil
}

// UpdateExpense updates an existing expense in the database and updates the associated budget.
func UpdateExpense(db *sql.DB, expense Expense) (Expense, error) {
	// Get current expense first
	currentExpense, err := GetExpenseByID(db, expense.ID)
	if err != nil {
		return Expense{}, err
	}

	// Validate fields
	if err := validateUpdateExpense(expense, currentExpense); err != nil {
		return Expense{}, err
	}

	// Build update query
	query := "UPDATE Expense SET"
	args := []interface{}{}
	argCount := 1
	updates := []string{}

	if expense.Amount != 0 {
		updates = append(updates, fmt.Sprintf("amount = $%d", argCount))
		args = append(args, expense.Amount)
		argCount++
	}

	if expense.Description != "" {
		updates = append(updates, fmt.Sprintf("description = $%d", argCount))
		args = append(args, expense.Description)
		argCount++
	}

	if !expense.Date.IsZero() {
		updates = append(updates, fmt.Sprintf("date = $%d", argCount))
		args = append(args, expense.Date)
		argCount++
	}

	if len(updates) == 0 {
		return currentExpense, nil // No changes made
	}

	// Finalize the query
	query += " " + strings.Join(updates, ", ")
	query += fmt.Sprintf(" WHERE id = $%d", argCount)
	args = append(args, expense.ID)

	if _, err := db.Exec(query, args...); err != nil {
		return Expense{}, err
	}

	// Update the budget if amount changed
	if expense.Amount != 0 && expense.Amount != currentExpense.Amount {
		if budgetExists, err := checkBudgetExists(db, currentExpense.CategoryID); err != nil {
			return Expense{}, err
		} else if budgetExists {
			if err := updateBudgetSpent(db, currentExpense.CategoryID, expense.Amount-currentExpense.Amount); err != nil {
				return Expense{}, err
			}
		}
	}

	// Merge the updated fields with current expense
	updatedExpense := currentExpense
	if expense.Amount != 0 {
		updatedExpense.Amount = expense.Amount
	}
	if expense.Description != "" {
		updatedExpense.Description = expense.Description
	}
	if !expense.Date.IsZero() {
		updatedExpense.Date = expense.Date
	}

	return updatedExpense, nil
}

// DeleteExpense removes an expense from the database and updates the associated budget.
func DeleteExpense(db *sql.DB, id int64) error {
	currentExpense, err := GetExpenseByID(db, id)
	if err != nil {
		return err
	}

	// Delete the expense from the database
	_, err = db.Exec("DELETE FROM Expense WHERE id = $1", id)
	if err != nil {
		return err
	}

	// Check if budget exists before updating it
	budgetExists, err := checkBudgetExists(db, currentExpense.CategoryID)
	if err != nil {
		return err
	}

	// Update the associated budget by deducting the amount
	if budgetExists {
		if err := updateBudgetSpent(db, currentExpense.CategoryID, -currentExpense.Amount); err != nil {
			return err // Deduct amount from the budget
		}
	}

	return nil
}

// updateBudgetSpent updates the spent amount for the associated budget category.
func updateBudgetSpent(db *sql.DB, categoryID int64, amount float64) error {
	_, err := db.Exec("UPDATE Budget SET spent = spent + $1 WHERE category_id = $2", amount, categoryID)
	return err
}

// checkBudgetExists checks if a budget exists for the given category ID.
func checkBudgetExists(db *sql.DB, categoryID int64) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Budget WHERE category_id = $1", categoryID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
