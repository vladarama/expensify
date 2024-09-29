package models_test

import (
	"testing"

	"expense-tracker/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetCategories(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow(1, "Food", "Food expenses").
		AddRow(2, "Transport", "Transportation expenses")

	mock.ExpectQuery("SELECT id, name, description FROM categories").WillReturnRows(rows)

	categories, err := models.GetCategories(db)

	assert.NoError(t, err)
	assert.Len(t, categories, 2)
	assert.Equal(t, int64(1), categories[0].ID)
	assert.Equal(t, "Food", categories[0].Name)
	assert.Equal(t, "Food expenses", categories[0].Description)
	assert.Equal(t, int64(2), categories[1].ID)
	assert.Equal(t, "Transport", categories[1].Name)
	assert.Equal(t, "Transportation expenses", categories[1].Description)
}

func TestGetCategoryByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow(1, "Food", "Food expenses")

	mock.ExpectQuery("SELECT id, name, description FROM categories WHERE name = \\$1").
		WithArgs("food").
		WillReturnRows(rows)

	category, err := models.GetCategoryByName(db, "Food")

	assert.NoError(t, err)
	assert.Equal(t, int64(1), category.ID)
	assert.Equal(t, "Food", category.Name)
	assert.Equal(t, "Food expenses", category.Description)
}

func TestCreateCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("INSERT INTO categories \\(name, description\\) VALUES \\(\\$1, \\$2\\) RETURNING id").
		WithArgs("entertainment", "entertainment expenses").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	category := models.Category{Name: "Entertainment", Description: "Entertainment expenses"}
	createdCategory, err := models.CreateCategory(db, category)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), createdCategory.ID)
	assert.Equal(t, "Entertainment", createdCategory.Name)
	assert.Equal(t, "Entertainment expenses", createdCategory.Description)
}

func TestUpdateCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("UPDATE categories SET description = \\$1 WHERE name = \\$2").
		WithArgs("updated food expenses", "food").
		WillReturnResult(sqlmock.NewResult(1, 1))

	category := models.Category{Name: "Food", Description: "Updated food expenses"}
	updatedCategory, err := models.UpdateCategory(db, category)

	assert.NoError(t, err)
	assert.Equal(t, "Food", updatedCategory.Name)
	assert.Equal(t, "Updated food expenses", updatedCategory.Description)
}

func TestDeleteCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE FROM categories WHERE name = \\$1").
		WithArgs("food").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = models.DeleteCategory(db, "Food")

	assert.NoError(t, err)
}