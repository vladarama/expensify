package categories

import (
    "database/sql"
    "time"
)

type Category struct {
    ID          int64     `json:"id"`
    UserID      int64     `json:"user_id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

func GetCategories(db *sql.DB, userID int64) ([]Category, error) {
    rows, err := db.Query(`
        SELECT id, user_id, name, description, created_at, updated_at
        FROM categories
        WHERE user_id=$1
        ORDER BY name ASC
    `, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var categoriesList []Category
    for rows.Next() {
        var cat Category
        err := rows.Scan(&cat.ID, &cat.UserID, &cat.Name, &cat.Description, &cat.CreatedAt, &cat.UpdatedAt)
        if err != nil {
            return nil, err
        }
        categoriesList = append(categoriesList, cat)
    }
    return categoriesList, nil
}

func CreateCategory(db *sql.DB, cat *Category) error {
    err := db.QueryRow(`
        INSERT INTO categories (user_id, name, description, created_at, updated_at)
        VALUES ($1, $2, $3, NOW(), NOW())
        RETURNING id, created_at, updated_at
    `, cat.UserID, cat.Name, cat.Description).Scan(&cat.ID, &cat.CreatedAt, &cat.UpdatedAt)
    return err
}
