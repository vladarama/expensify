package models

import (
	"database/sql"
	"strings"
)

type Category struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}


func GetCategories(db *sql.DB) ([]Category, error) {
	rows, err := db.Query("SELECT id, name, description FROM categories")
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
	err := db.QueryRow("SELECT id, name, description FROM categories WHERE name = $1", strings.ToLower(name)).Scan(&category.ID, &category.Name, &category.Description)
	if err != nil {
		return Category{}, err
	}
	return category, nil
}

func CreateCategory(db *sql.DB, category Category) (Category, error) {
	var id int64
	err := db.QueryRow("INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id", 
		strings.ToLower(category.Name), strings.ToLower(category.Description)).Scan(&id)
	if err != nil {
		return Category{}, err
	}
	category.ID = id
	return category, nil
}

func UpdateCategory(db *sql.DB, category Category) (Category, error) {
	_, err := db.Exec("UPDATE categories SET description = $1 WHERE name = $2", 
		strings.ToLower(category.Description), strings.ToLower(category.Name))
	if err != nil {
		return Category{}, err
	}
	return category, nil
}

func DeleteCategory(db *sql.DB, name string) error {
	_, err := db.Exec("DELETE FROM categories WHERE name = $1", strings.ToLower(name))
	if err != nil {
		return err
	}
	return nil
}
