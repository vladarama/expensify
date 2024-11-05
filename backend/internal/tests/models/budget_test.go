package models_test

import (
	"errors"
	"expense-tracker/internal/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetBudgets(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "category_id", "amount", "spent", "start_date", "end_date"}).
		AddRow(1, int64(1), 500.0, 200.0, time.Now(), time.Now().Add(30*time.Hour*24)).
		AddRow(2, int64(2), 300.0, 100.0, time.Now(), time.Now().Add(30*time.Hour*24))

	mock.ExpectQuery("SELECT id, category_id, amount, spent, start_date, end_date FROM Budget").WillReturnRows(rows)

	budgets, err := models.GetBudgets(db)

	assert.NoError(t, err)
	assert.Len(t, budgets, 2)
	assert.Equal(t, int64(1), budgets[0].CategoryID) // Updated to int64
	assert.Equal(t, float64(500), budgets[0].Amount)
	assert.Equal(t, float64(200), budgets[0].Spent)
	assert.Equal(t, int64(2), budgets[1].CategoryID) // Updated to int64
	assert.Equal(t, float64(300), budgets[1].Amount)
}

func TestGetBudgetByCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "category_id", "Amount", "spent", "start_date", "end_date"}).
		AddRow(1, int64(1), 500.0, 200.0, time.Now(), time.Now().Add(30*time.Hour*24))

	mock.ExpectQuery("SELECT id, category_id, Amount, spent, start_date, end_date FROM Budget WHERE category_id = \\$1").
		WithArgs(int64(1)).
		WillReturnRows(rows)

	budget, err := models.GetBudgetByCategory(db, "1") // Pass CategoryID as string

	assert.NoError(t, err)
	assert.Equal(t, int64(1), budget.CategoryID) // Use int64 for CategoryID
	assert.Equal(t, float64(500), budget.Amount)
	assert.Equal(t, float64(200), budget.Spent)
}

func TestCreateBudget(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("INSERT INTO Budget \\(category_id, amount, spent, start_date, end_date\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5\\) RETURNING id").
		WithArgs(int64(1), 500.0, 200.0, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	budget := models.Budget{
		CategoryID: int64(1), // Change to int64
		Amount:     500.0,
		Spent:      200.0,
		StartDate:  time.Now(),
		EndDate:    time.Now().Add(30 * 24 * time.Hour),
	}
	createdBudget, err := models.CreateBudget(db, budget)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), createdBudget.ID)
	assert.Equal(t, int64(1), createdBudget.CategoryID)
	assert.Equal(t, float64(500), createdBudget.Amount)
	assert.Equal(t, float64(200), createdBudget.Spent)
}

func TestCreateBudgetWithEmptyCategory(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	budget := models.Budget{
		CategoryID: 0, // Assuming you want to allow 0 but handle it appropriately in your model logic
		Amount:     500.0,
		Spent:      200.0,
		StartDate:  time.Now(),
		EndDate:    time.Now().Add(30 * 24 * time.Hour),
	}
	_, err = models.CreateBudget(db, budget)

	assert.Error(t, err)
	assert.Equal(t, errors.New("category cannot be empty"), err)
}

func TestUpdateBudget(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("UPDATE Budget SET amount = \\$1, spent = \\$2, start_date = \\$3, end_date = \\$4 WHERE category_id = \\$5").
		WithArgs(600.0, 250.0, sqlmock.AnyArg(), sqlmock.AnyArg(), int64(1)). // Change to int64
		WillReturnResult(sqlmock.NewResult(1, 1))

	budget := models.Budget{
		CategoryID: int64(1), // Change to int64
		Amount:     600.0,
		Spent:      250.0,
		StartDate:  time.Now(),
		EndDate:    time.Now().Add(30 * 24 * time.Hour),
	}
	updatedBudget, err := models.UpdateBudget(db, budget)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), updatedBudget.CategoryID) // Updated to int64
	assert.Equal(t, float64(600), updatedBudget.Amount)
	assert.Equal(t, float64(250), updatedBudget.Spent)
}

func TestUpdateBudgetWithEmptyCategory(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	budget := models.Budget{
		CategoryID: 0, // Assuming you want to allow 0 but handle it appropriately in your model logic
		Amount:     600.0,
		Spent:      250.0,
		StartDate:  time.Now(),
		EndDate:    time.Now().Add(30 * 24 * time.Hour),
	}
	_, err = models.UpdateBudget(db, budget)

	assert.Error(t, err)
	assert.Equal(t, errors.New("category cannot be empty"), err)
}

func TestDeleteBudget(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE FROM Budget WHERE category_id = \\$1").
		WithArgs(int64(1)). // Ensure you're using the correct SQL type here
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Pass the CategoryID as a string if models.DeleteBudget requires a string
	err = models.DeleteBudget(db, "1") // Pass CategoryID as string

	assert.NoError(t, err)
}
