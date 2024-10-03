package models

import (
	"database/sql"
	"errors"
	"strings"
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

func GetCategoryByName(db *sql.DB, name string) (Category, error) {
	var category Category
	err := db.QueryRow("SELECT id, name, description FROM Category WHERE name = $1", strings.ToLower(name)).Scan(&category.ID, &category.Name, &category.Description)
	if err != nil {
		return Category{}, err
	}
	return category, nil
}

func CreateCategory(db *sql.DB, category Category) (Category, error) {
	if strings.TrimSpace(category.Name) == "" {
		return Category{}, errors.New("category name cannot be empty")
	}

	var id int64
	err := db.QueryRow("INSERT INTO Category (name, description) VALUES ($1, $2) RETURNING id", 
		strings.ToLower(strings.TrimSpace(category.Name)), strings.ToLower(category.Description)).Scan(&id)
	if err != nil {
		return Category{}, err
	}
	category.ID = id
	return category, nil
}

func UpdateCategory(db *sql.DB, category Category) (Category, error) {
	if strings.TrimSpace(category.Name) == "" {
		return Category{}, errors.New("category name cannot be empty")
	}

	_, err := db.Exec("UPDATE Category SET description = $1 WHERE name = $2", 
		strings.ToLower(category.Description), strings.ToLower(strings.TrimSpace(category.Name)))
	if err != nil {
		return Category{}, err
	}
	return category, nil
}

func DeleteCategory(db *sql.DB, name string) error {
	_, err := db.Exec("DELETE FROM Category WHERE name = $1", strings.ToLower(name))
	if err != nil {
		return err
	}
	return nil
}
