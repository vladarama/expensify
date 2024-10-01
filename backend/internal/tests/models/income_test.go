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

	rows := sqlmock.NewRows([]string{"id", "category_id", "amount", "date", "source"}).
		AddRow(1, 1, 1000.00, time.Now(), "Salary").
		AddRow(2, 2, 500.00, time.Now(), "Freelance")

	mock.ExpectQuery("SELECT id, category_id, amount, date, source FROM incomes").WillReturnRows(rows)

	incomes, err := models.GetIncomes(db)

	assert.NoError(t, err)
	assert.Len(t, incomes, 2)
	assert.Equal(t, int64(1), incomes[0].ID)
	assert.Equal(t, int64(1), incomes[0].CategoryID)
	assert.Equal(t, 1000.00, incomes[0].Amount)
	assert.Equal(t, "Salary", incomes[0].Source)
}

func TestGetIncomeByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "category_id", "amount", "date", "source"}).
		AddRow(1, 1, 1000.00, time.Now(), "Salary")

	mock.ExpectQuery("SELECT id, category_id, amount, date, source FROM incomes WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(rows)

	income, err := models.GetIncomeByID(db, 1)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), income.ID)
	assert.Equal(t, int64(1), income.CategoryID)
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
		CategoryID: 1,
		Amount:     1000.00,
		Date:       time.Now(),
		Source:     "Salary",
	}

	mock.ExpectQuery("INSERT INTO incomes \\(category_id, amount, date, source\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) RETURNING id").
		WithArgs(income.CategoryID, income.Amount, income.Date, income.Source).
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

	income := models.Income{
		ID:         1,
		CategoryID: 1,
		Amount:     1500.00,
		Date:       time.Now(),
		Source:     "Updated Salary",
	}

	mock.ExpectExec("UPDATE incomes SET category_id = \\$1, amount = \\$2, date = \\$3, source = \\$4 WHERE id = \\$5").
		WithArgs(income.CategoryID, income.Amount, income.Date, income.Source, income.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	updatedIncome, err := models.UpdateIncome(db, income)

	assert.NoError(t, err)
	assert.Equal(t, income, updatedIncome)
}

func TestDeleteIncome(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE FROM incomes WHERE id = \\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = models.DeleteIncome(db, 1)

	assert.NoError(t, err)
}