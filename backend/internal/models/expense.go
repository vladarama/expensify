package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Expense struct {
	ID         int64     `json:"id"`
	CategoryID int64     `json:"category_id"`
	Amount     float64   `json:"amount"`
	Date       time.Time `json:"date"`
}

func GetExpenses(db *sql.DB) ([]Expense, error) {
	rows, err := db.Query("SELECT id, category_id, amount, date FROM Expense")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []Expense
	for rows.Next() {
		var expense Expense
		err := rows.Scan(&expense.ID, &expense.CategoryID, &expense.Amount, &expense.Date)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}

	return expenses, nil
}

func GetExpenseByID(db *sql.DB, id int64) (Expense, error) {
	var expense Expense
	err := db.QueryRow("SELECT id, category_id, amount, date FROM Expense WHERE id = $1", id).
		Scan(&expense.ID, &expense.CategoryID, &expense.Amount, &expense.Date)
	if err != nil {
		return Expense{}, err
	}
	return expense, nil
}

func CreateExpense(db *sql.DB, expense Expense) (Expense, error) {
	if expense.Amount <= 0 {
		return Expense{}, errors.New("amount must be greater than zero")
	}

	if expense.Date.IsZero() {
		return Expense{}, errors.New("date must be provided")
	}
	if expense.Date.After(time.Now()) {
		return Expense{}, errors.New("date cannot be in the future")
	}

	if expense.CategoryID <= 0 {
		return Expense{}, errors.New("category ID must be provided")
	}

	err := db.QueryRow("INSERT INTO Expense (category_id, amount, date) VALUES ($1, $2, $3) RETURNING id",
		expense.CategoryID, expense.Amount, expense.Date).Scan(&expense.ID)
	if err != nil {
		return Expense{}, err
	}
	return expense, nil
}

func UpdateExpense(db *sql.DB, expense Expense) (Expense, error) {
	currentExpense, err := GetExpenseByID(db, expense.ID)
	if err != nil {
		return Expense{}, err
	}

	query := "UPDATE Expense SET"
	args := []interface{}{}
	argCount := 1

	if expense.CategoryID != 0 {
		query += fmt.Sprintf(" category_id = $%d,", argCount)
		args = append(args, expense.CategoryID)
		argCount++
	} else {
		expense.CategoryID = currentExpense.CategoryID
	}

	if expense.Amount != 0 {
		query += fmt.Sprintf(" amount = $%d,", argCount)
		args = append(args, expense.Amount)
		argCount++
	} else {
		expense.Amount = currentExpense.Amount
	}

	if !expense.Date.IsZero() {
		query += fmt.Sprintf(" date = $%d,", argCount)
		args = append(args, expense.Date)
		argCount++
	} else {
		expense.Date = currentExpense.Date
	}

	query = strings.TrimSuffix(query, ",")
	query += fmt.Sprintf(" WHERE id = $%d", argCount)
	args = append(args, expense.ID)

	_, err = db.Exec(query, args...)
	if err != nil {
		return Expense{}, err
	}

	return expense, nil
}

func DeleteExpense(db *sql.DB, id int64) error {
	_, err := db.Exec("DELETE FROM Expense WHERE id = $1", id)
	return err
}