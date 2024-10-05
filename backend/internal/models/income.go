package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Income struct {
	ID     int64     `json:"id"`
	Amount float64   `json:"amount"`
	Date   time.Time `json:"date"`
	Source string    `json:"source"`
}

func GetIncomes(db *sql.DB) ([]Income, error) {
	rows, err := db.Query("SELECT id, amount, date, source FROM Income")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incomes []Income
	for rows.Next() {
		var income Income
		err := rows.Scan(&income.ID, &income.Amount, &income.Date, &income.Source)
		if err != nil {
			return nil, err
		}
		incomes = append(incomes, income)
	}

	return incomes, nil
}

func GetIncomeByID(db *sql.DB, id int64) (Income, error) {
	var income Income
	err := db.QueryRow("SELECT id, amount, date, source FROM Income WHERE id = $1", id).
		Scan(&income.ID, &income.Amount, &income.Date, &income.Source)
	if err != nil {
		return Income{}, err
	}
	return income, nil
}

func CreateIncome(db *sql.DB, income Income) (Income, error) {
	// Validate Amount
	if income.Amount <= 0 {
		return Income{}, errors.New("amount must be greater than zero")
	}

	// Validate Date
	if income.Date.IsZero() {
		return Income{}, errors.New("date must be provided")
	}
	if income.Date.After(time.Now()) {
		return Income{}, errors.New("date cannot be in the future")
	}

	// Validate Source
	if income.Source == "" {
		return Income{}, errors.New("source must be provided")
	}
	if len(income.Source) > 255 {  // Assuming a reasonable max length for the source field
		return Income{}, errors.New("source is too long (max 255 characters)")
	}

	// If all validations pass, insert into database
	err := db.QueryRow("INSERT INTO Income (amount, date, source) VALUES ($1, $2, $3) RETURNING id",
		income.Amount, income.Date, income.Source).Scan(&income.ID)
	if err != nil {
		return Income{}, err
	}
	return income, nil
}

func UpdateIncome(db *sql.DB, income Income) (Income, error) {
	// Fetch the current income data
	currentIncome, err := GetIncomeByID(db, income.ID)
	if err != nil {
		return Income{}, err
	}

	// Prepare the update query and arguments
	query := "UPDATE Income SET"
	args := []interface{}{}
	argCount := 1

	// Check and update each field if it's provided
	if income.Amount != 0 {
		query += fmt.Sprintf(" amount = $%d,", argCount)
		args = append(args, income.Amount)
		argCount++
	} else {
		income.Amount = currentIncome.Amount
	}

	if !income.Date.IsZero() {
		query += fmt.Sprintf(" date = $%d,", argCount)
		args = append(args, income.Date)
		argCount++
	} else {
		income.Date = currentIncome.Date
	}

	if income.Source != "" {
		query += fmt.Sprintf(" source = $%d,", argCount)
		args = append(args, income.Source)
		argCount++
	} else {
		income.Source = currentIncome.Source
	}

	// Remove the trailing comma and add the WHERE clause
	query = strings.TrimSuffix(query, ",")
	query += fmt.Sprintf(" WHERE id = $%d", argCount)
	args = append(args, income.ID)

	// Execute the update query
	_, err = db.Exec(query, args...)
	if err != nil {
		return Income{}, err
	}

	return income, nil
}

func DeleteIncome(db *sql.DB, id int64) error {
	_, err := db.Exec("DELETE FROM Income WHERE id = $1", id)
	return err
}