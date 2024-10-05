package models_test

import (
	"expense-tracker/internal/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetIncomes(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "amount", "date", "source"}).
		AddRow(1, 1000.00, time.Now(), "Salary").
		AddRow(2, 500.00, time.Now(), "Freelance")

	mock.ExpectQuery("SELECT id, amount, date, source FROM Income").WillReturnRows(rows)

	incomes, err := models.GetIncomes(db)

	assert.NoError(t, err)
	assert.Len(t, incomes, 2)
	assert.Equal(t, int64(1), incomes[0].ID)
	assert.Equal(t, 1000.00, incomes[0].Amount)
	assert.Equal(t, "Salary", incomes[0].Source)
}

func TestGetIncomeByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "amount", "date", "source"}).
		AddRow(1, 1000.00, time.Now(), "Salary")

	mock.ExpectQuery("SELECT id, amount, date, source FROM Income WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(rows)

	income, err := models.GetIncomeByID(db, 1)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), income.ID)
	assert.Equal(t, 1000.00, income.Amount)
	assert.Equal(t, "Salary", income.Source)
}

func TestCreateIncome(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	income := models.Income{
		Amount: 1000.00,
		Date:   time.Now(),
		Source: "Salary",
	}

	mock.ExpectQuery("INSERT INTO Income \\(amount, date, source\\) VALUES \\(\\$1, \\$2, \\$3\\) RETURNING id").
		WithArgs(income.Amount, income.Date, income.Source).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	createdIncome, err := models.CreateIncome(db, income)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), createdIncome.ID)
}

func TestUpdateIncome(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Mock the current income data
	currentIncome := models.Income{
		ID:     1,
		Amount: 1000.00,
		Date:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		Source: "Original Salary",
	}

	// Mock the GetIncomeByID call
	mock.ExpectQuery("SELECT id, amount, date, source FROM Income WHERE id = \\$1").
		WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "amount", "date", "source"}).
				AddRow(currentIncome.ID, currentIncome.Amount, currentIncome.Date, currentIncome.Source))

	// Test case 1: Update only amount
	updatedIncome := models.Income{
		ID:     1,
		Amount: 1500.00,
	}

	mock.ExpectExec("UPDATE Income SET amount = \\$1 WHERE id = \\$2").
		WithArgs(updatedIncome.Amount, updatedIncome.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := models.UpdateIncome(db, updatedIncome)

	assert.NoError(t, err)
	assert.Equal(t, updatedIncome.Amount, result.Amount)
	assert.Equal(t, currentIncome.Date, result.Date)
	assert.Equal(t, currentIncome.Source, result.Source)

	// Test case 2: Update amount and source
	updatedIncome = models.Income{
		ID:     1,
		Amount: 2000.00,
		Source: "Updated Salary",
	}

	mock.ExpectQuery("SELECT id, amount, date, source FROM Income WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "amount", "date", "source"}).
			AddRow(currentIncome.ID, currentIncome.Amount, currentIncome.Date, currentIncome.Source))

	mock.ExpectExec("UPDATE Income SET amount = \\$1, source = \\$2 WHERE id = \\$3").
		WithArgs(updatedIncome.Amount, updatedIncome.Source, updatedIncome.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err = models.UpdateIncome(db, updatedIncome)

	assert.NoError(t, err)
	assert.Equal(t, updatedIncome.Amount, result.Amount)
	assert.Equal(t, currentIncome.Date, result.Date)
	assert.Equal(t, updatedIncome.Source, result.Source)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteIncome(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE FROM Income WHERE id = \\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = models.DeleteIncome(db, 1)

	assert.NoError(t, err)
}