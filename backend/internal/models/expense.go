package models

import "time"

type Expense struct {
    ID         int64     `json:"id"`
    CategoryID int64     `json:"category_id"`
    Amount     float64   `json:"amount"`
    Date       time.Time `json:"date"`
}

// TODO: Add methods for CRUD operations on expenses