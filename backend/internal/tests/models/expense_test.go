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

	rows := sqlmock.NewRows([]string{"id", "category_id", "amount", "date"}).
		AddRow(1, 1, 100.00, time.Now()).
		AddRow(2, 2, 50.00, time.Now())

	mock.ExpectQuery("SELECT id, category_id, amount, date FROM Expense").WillReturnRows(rows)

	expenses, err := models.GetExpenses(db)

	assert.NoError(t, err)
	assert.Len(t, expenses, 2)
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

	rows := sqlmock.NewRows([]string{"id", "category_id", "amount", "date"}).
		AddRow(1, 1, 100.00, time.Now())

	mock.ExpectQuery("SELECT id, category_id, amount, date FROM Expense WHERE id = \\$1").
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

	expense := models.Expense{
		CategoryID:  1,
		Amount:      100.00,
		Date:        time.Now(),
	}

	mock.ExpectQuery("INSERT INTO Expense \\(category_id, amount, date\\) VALUES \\(\\$1, \\$2, \\$3\\) RETURNING id").
		WithArgs(expense.CategoryID, expense.Amount, expense.Date).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	createdExpense, err := models.CreateExpense(db, expense)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), createdExpense.ID)
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
	}

	mock.ExpectQuery("SELECT id, category_id, amount, date FROM Expense WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "category_id", "amount", "date"}).
			AddRow(currentExpense.ID, currentExpense.CategoryID, currentExpense.Amount, currentExpense.Date))

	updatedExpense := models.Expense{
		ID:          1,
		Amount:      150.00,
	}

	mock.ExpectExec("UPDATE Expense SET amount = \\$1 WHERE id = \\$2").
		WithArgs(updatedExpense.Amount, updatedExpense.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := models.UpdateExpense(db, updatedExpense)

	assert.NoError(t, err)
	assert.Equal(t, updatedExpense.Amount, result.Amount)
	assert.Equal(t, currentExpense.Date, result.Date)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteExpense(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE FROM Expense WHERE id = \\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = models.DeleteExpense(db, 1)

	assert.NoError(t, err)
}