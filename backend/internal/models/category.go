package models

import (
	"database/sql"
	"errors"
	"fmt"
)

type Category struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func GetCategories(db *sql.DB) ([]Category, error) {
	rows, err := db.Query("SELECT id, name, description FROM Category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.Name, &category.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func GetCategoryByID(db *sql.DB, id int64) (Category, error) {
	var category Category
	err := db.QueryRow("SELECT id, name, description FROM Category WHERE id = $1", id).
		Scan(&category.ID, &category.Name, &category.Description)
	if err != nil {
		return Category{}, err
	}
	return category, nil
}

func CreateCategory(db *sql.DB, category Category) (Category, error) {
	if category.Name == "" {
		return Category{}, errors.New("name must be provided")
	}

	err := db.QueryRow(
		"INSERT INTO Category (name, description) VALUES ($1, $2) RETURNING id",
		category.Name, category.Description,
	).Scan(&category.ID)
	if err != nil {
		return Category{}, err
	}
	return category, nil
}

func UpdateCategory(db *sql.DB, category Category) (Category, error) {
	if category.ID == 0 {
		return Category{}, errors.New("id must be provided")
	}

	_, err := db.Exec(
		"UPDATE Category SET name = $1, description = $2 WHERE id = $3",
		category.Name, category.Description, category.ID,
	)
	if err != nil {
		return Category{}, err
	}
	return category, nil
}

func DeleteCategory(db *sql.DB, id int64) error {
	// Prevent deletion of the "Other" category
	if id == 1 {
		return fmt.Errorf("cannot delete the 'Other' category")
	}

	// Reassign all expenses to the "Other" category
	_, err := db.Exec("UPDATE Expense SET category_id = 1 WHERE category_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to reassign expenses to 'Other': %w", err)
	}

	result, err := db.Exec("DELETE FROM Category WHERE id = $1", id)
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
