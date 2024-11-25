package models_test

import (
	"database/sql"
	"errors"
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

	mock.ExpectQuery("SELECT id, name, description FROM Category").WillReturnRows(rows)

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

func TestGetCategoryByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow(1, "Food", "Food expenses")

	mock.ExpectQuery("SELECT id, name, description FROM Category WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(rows)

	category, err := models.GetCategoryByID(db, 1)

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

	mock.ExpectQuery("INSERT INTO Category \\(name, description\\) VALUES \\(\\$1, \\$2\\) RETURNING id").
		WithArgs("Entertainment", "Entertainment expenses").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	category := models.Category{Name: "Entertainment", Description: "Entertainment expenses"}
	createdCategory, err := models.CreateCategory(db, category)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), createdCategory.ID)
	assert.Equal(t, "Entertainment", createdCategory.Name)
	assert.Equal(t, "Entertainment expenses", createdCategory.Description)
}

func TestCreateCategoryWithEmptyName(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	category := models.Category{Name: "", Description: "Empty category"}
	_, err = models.CreateCategory(db, category)

	assert.Error(t, err)
	assert.Equal(t, errors.New("name must be provided"), err)
}

func TestUpdateCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("UPDATE Category SET name = \\$1, description = \\$2 WHERE id = \\$3").
		WithArgs("Food", "Updated food expenses", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	category := models.Category{
		ID:          1,
		Name:        "Food",
		Description: "Updated food expenses",
	}
	updatedCategory, err := models.UpdateCategory(db, category)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), updatedCategory.ID)
	assert.Equal(t, "Food", updatedCategory.Name)
	assert.Equal(t, "Updated food expenses", updatedCategory.Description)
}

func TestUpdateCategoryWithoutID(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	category := models.Category{Name: "Food", Description: "Updated category"}
	_, err = models.UpdateCategory(db, category)

	assert.Error(t, err)
	assert.Equal(t, errors.New("id must be provided"), err)
}
func TestDeleteCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Test case: Prevent deletion of the "Other" category
	err = models.DeleteCategory(db, 1)
	assert.Error(t, err)
	assert.Equal(t, "cannot delete the 'Other' category", err.Error())

	// Test case: Reassign expenses to "Other" and delete a category
	mock.ExpectExec(`UPDATE Expense SET category_id = 1 WHERE category_id = \$1`).
		WithArgs(int64(2)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(`DELETE FROM Category WHERE id = \$1`).
		WithArgs(int64(2)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = models.DeleteCategory(db, 2)
	assert.NoError(t, err)

	// Test case: Category not found
	mock.ExpectExec(`UPDATE Expense SET category_id = 1 WHERE category_id = \$1`).
		WithArgs(int64(3)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectExec(`DELETE FROM Category WHERE id = \$1`).
		WithArgs(int64(3)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = models.DeleteCategory(db, 3)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
}
