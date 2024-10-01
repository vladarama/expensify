package models

import (
	"database/sql"
	"time"
)

type Income struct {
	ID         int64     `json:"id"`
	CategoryID int64     `json:"category_id"`
	Amount     float64   `json:"amount"`
	Date       time.Time `json:"date"`
	Source     string    `json:"source"`
}

func GetIncomes(db *sql.DB) ([]Income, error) {
	rows, err := db.Query("SELECT id, category_id, amount, date, source FROM incomes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incomes []Income
	for rows.Next() {
		var income Income
		err := rows.Scan(&income.ID, &income.CategoryID, &income.Amount, &income.Date, &income.Source)
		if err != nil {
			return nil, err
		}
		incomes = append(incomes, income)
	}

	return incomes, nil
}

func GetIncomeByID(db *sql.DB, id int64) (Income, error) {
	var income Income
	err := db.QueryRow("SELECT id, category_id, amount, date, source FROM incomes WHERE id = $1", id).
		Scan(&income.ID, &income.CategoryID, &income.Amount, &income.Date, &income.Source)
	if err != nil {
		return Income{}, err
	}
	return income, nil
}

func CreateIncome(db *sql.DB, income Income) (Income, error) {
	err := db.QueryRow("INSERT INTO incomes (category_id, amount, date, source) VALUES ($1, $2, $3, $4) RETURNING id",
		income.CategoryID, income.Amount, income.Date, income.Source).Scan(&income.ID)
	if err != nil {
		return Income{}, err
	}
	return income, nil
}

func UpdateIncome(db *sql.DB, income Income) (Income, error) {
	_, err := db.Exec("UPDATE incomes SET category_id = $1, amount = $2, date = $3, source = $4 WHERE id = $5",
		income.CategoryID, income.Amount, income.Date, income.Source, income.ID)
	if err != nil {
		return Income{}, err
	}
	return income, nil
}

func DeleteIncome(db *sql.DB, id int64) error {
	_, err := db.Exec("DELETE FROM incomes WHERE id = $1", id)
	return err
}