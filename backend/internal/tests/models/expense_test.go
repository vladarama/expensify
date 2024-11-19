package models_test

import (
	"expense-tracker/internal/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetExpenses(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "category_id", "amount", "date", "description"}).
		AddRow(1, 1, 100.00, time.Now(), "Groceries").
		AddRow(2, 2, 50.00, time.Now(), "Utilities")

	mock.ExpectQuery("SELECT id, category_id, amount, date, description FROM Expense").WillReturnRows(rows)

	expenses, err := models.GetExpenses(db)

	assert.NoError(t, err)
	assert.Len(t, expenses, 2)
	assert.Equal(t, "Groceries", expenses[0].Description)
	assert.Equal(t, int64(1), expenses[0].ID)
	assert.Equal(t, int64(1), expenses[0].CategoryID)
	assert.Equal(t, 100.00, expenses[0].Amount)
}

func TestGetExpenseByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "category_id", "amount", "date", "description"}).
		AddRow(1, 1, 100.00, time.Now(), "Groceries")

	mock.ExpectQuery("SELECT id, category_id, amount, date, description FROM Expense WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(rows)

	expense, err := models.GetExpenseByID(db, 1)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), expense.ID)
	assert.Equal(t, int64(1), expense.CategoryID)
	assert.Equal(t, 100.00, expense.Amount)
}

func TestCreateExpense(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	now := time.Now()
	expense := models.Expense{
		Description: "Groceries",
		CategoryID:  1,
		Amount:      100.00,
		 Date:        now,
	}

	mock.ExpectQuery("INSERT INTO Expense \\(category_id, amount, date, description\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) RETURNING id, category_id, amount, date, description").
		WithArgs(expense.CategoryID, expense.Amount, expense.Date, expense.Description).
		WillReturnRows(sqlmock.NewRows([]string{"id", "category_id", "amount", "date", "description"}).
			AddRow(1, expense.CategoryID, expense.Amount, expense.Date, expense.Description))

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM Budget WHERE category_id = \\$1").
		WithArgs(expense.CategoryID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectExec("UPDATE Budget SET spent = spent \\+ \\$1 WHERE category_id = \\$2").
		WithArgs(expense.Amount, expense.CategoryID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	createdExpense, err := models.CreateExpense(db, expense)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), createdExpense.ID)
	assert.Equal(t, expense.CategoryID, createdExpense.CategoryID)
	assert.Equal(t, expense.Amount, createdExpense.Amount)
	assert.Equal(t, expense.Date, createdExpense.Date)
	assert.Equal(t, expense.Description, createdExpense.Description)

	// Verify that all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateExpense(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	currentExpense := models.Expense{
		ID:          1,
		CategoryID:  1,
		Amount:      100.00,
		Date:        time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		Description: "Groceries",
	}

	// First expect the GetExpenseByID query
	mock.ExpectQuery("SELECT id, category_id, amount, date, description FROM Expense WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "category_id", "amount", "date", "description"}).
			AddRow(currentExpense.ID, currentExpense.CategoryID, currentExpense.Amount, currentExpense.Date, currentExpense.Description))

	updatedExpense := models.Expense{
		ID:          1,
		Amount:      150.00,
		Description: "Updated Groceries",
	}

	// Then expect the update query
	mock.ExpectExec("UPDATE Expense SET amount = \\$1, description = \\$2 WHERE id = \\$3").
		WithArgs(updatedExpense.Amount, updatedExpense.Description, updatedExpense.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Then check if budget exists
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM Budget WHERE category_id = \\$1").
		WithArgs(currentExpense.CategoryID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// Finally update budget spent amount (150 - 100 = 50 difference)
	mock.ExpectExec("UPDATE Budget SET spent = spent \\+ \\$1 WHERE category_id = \\$2").
		WithArgs(50.0, currentExpense.CategoryID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := models.UpdateExpense(db, updatedExpense)

	assert.NoError(t, err)
	assert.Equal(t, updatedExpense.Amount, result.Amount)
	assert.Equal(t, updatedExpense.Description, result.Description)
	assert.Equal(t, currentExpense.CategoryID, result.CategoryID)
	assert.Equal(t, currentExpense.Date, result.Date)

	// Verify that all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteExpense(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Get the expense details first
	mock.ExpectQuery("SELECT id, category_id, amount, date, description FROM Expense WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "category_id", "amount", "date", "description"}).
			AddRow(1, 1, 100.00, time.Now(), "Test Expense"))

	// Delete the expense first
	mock.ExpectExec("DELETE FROM Expense WHERE id = \\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Then check if budget exists
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM Budget WHERE category_id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// Finally update budget spent amount (using negative amount to decrease the spent value)
	mock.ExpectExec("UPDATE Budget SET spent = spent \\+ \\$1 WHERE category_id = \\$2").
		WithArgs(-100.00, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = models.DeleteExpense(db, 1)

	assert.NoError(t, err)
}