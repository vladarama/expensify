package models

import "time"

type Income struct {
    ID         int64     `json:"id"`
    CategoryID int64     `json:"category_id"`
    Amount     float64   `json:"amount"`
    Date       time.Time `json:"date"`
    Source     string    `json:"source"`
}

// TODO: Add methods for CRUD operations on incomes