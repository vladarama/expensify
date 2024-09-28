package expenses

import (
    "database/sql"
    "time"
)

type Expense struct {
    ID         int64     `json:"id"`
    UserID     int64     `json:"user_id"`
    CategoryID int64     `json:"category_id"`
    Amount     float64   `json:"amount"`
    Date       time.Time `json:"date"`
    Note       string    `json:"note"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}

func GetExpenses(db *sql.DB, userID int64) ([]Expense, error) {
    rows, err := db.Query(`
        SELECT id, user_id, category_id, amount, date, note, created_at, updated_at
        FROM expenses
        WHERE user_id=$1
        ORDER BY date DESC
    `, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var expensesList []Expense
    for rows.Next() {
        var exp Expense
        err := rows.Scan(&exp.ID, &exp.UserID, &exp.CategoryID, &exp.Amount, &exp.Date, &exp.Note, &exp.CreatedAt, &exp.UpdatedAt)
        if err != nil {
            return nil, err
        }
        expensesList = append(expensesList, exp)
    }
    return expensesList, nil
}

func CreateExpense(db *sql.DB, exp *Expense) error {
    err := db.QueryRow(`
        INSERT INTO expenses (user_id, category_id, amount, date, note, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
        RETURNING id, created_at, updated_at
    `, exp.UserID, exp.CategoryID, exp.Amount, exp.Date, exp.Note).Scan(&exp.ID, &exp.CreatedAt, &exp.UpdatedAt)
    return err
}
