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

	rows := sqlmock.NewRows([]string{"id", "category_id", "amount", "spent", "start_date", "end_date"}).
		AddRow(1, int64(1), 500.0, 200.0, time.Now(), time.Now().Add(30*time.Hour*24))

	mock.ExpectQuery("SELECT id, category_id, amount, spent, start_date, end_date FROM Budget WHERE category_id = \\$1").
		WithArgs(int64(1)).
		WillReturnRows(rows)

	budgets, err := models.GetBudgetsByCategoryID(db, 1)

	assert.NoError(t, err)
	assert.Len(t, budgets, 1)
	assert.Equal(t, int64(1), budgets[0].CategoryID)
	assert.Equal(t, float64(500), budgets[0].Amount)
	assert.Equal(t, float64(200), budgets[0].Spent)
}
func TestCreateBudget(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Mock DoesBudgetOverlap query
	mock.ExpectQuery(`SELECT EXISTS \( SELECT 1 FROM Budget WHERE category_id = \$1 AND id <> \$4 AND \( \(start_date <= \$3 AND end_date >= \$2\) \) \)`).
		WithArgs(int64(1), sqlmock.AnyArg(), sqlmock.AnyArg(), int64(0)).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	// Mock CalculateTotalSpent query
	mock.ExpectQuery("SELECT COALESCE\\(SUM\\(amount\\), 0\\) FROM Expense WHERE category_id = \\$1 AND date >= \\$2 AND date <= \\$3").
		WithArgs(int64(1), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"total_spent"}).AddRow(200.0))

	// Mock Insert query
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
		CategoryID: 0,
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

	// Mock DoesBudgetOverlap query
	mock.ExpectQuery(`SELECT EXISTS \( SELECT 1 FROM Budget WHERE category_id = \$1 AND id <> \$4 AND \( \(start_date <= \$3 AND end_date >= \$2\) \) \)`).
		WithArgs(int64(1), sqlmock.AnyArg(), sqlmock.AnyArg(), int64(1)). // Ensure `id <> $4` matches the budget being updated
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	// Mock CalculateTotalSpent query
	mock.ExpectQuery("SELECT COALESCE\\(SUM\\(amount\\), 0\\) FROM Expense WHERE category_id = \\$1 AND date >= \\$2 AND date <= \\$3").
		WithArgs(int64(1), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"total_spent"}).AddRow(250.0))

	// Mock Update query
	mock.ExpectExec("UPDATE Budget SET amount = \\$1, spent = \\$2, start_date = \\$3, end_date = \\$4 WHERE category_id = \\$5").
		WithArgs(600.0, 250.0, sqlmock.AnyArg(), sqlmock.AnyArg(), int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	budget := models.Budget{
		ID:         1,                                   // Ensure ID matches the updated record
		CategoryID: int64(1),                            // Category ID
		Amount:     600.0,                               // Updated Amount
		Spent:      250.0,                               // Updated Spent
		StartDate:  time.Now(),                          // Updated StartDate
		EndDate:    time.Now().Add(30 * 24 * time.Hour), // Updated EndDate
	}
	updatedBudget, err := models.UpdateBudget(db, budget)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), updatedBudget.CategoryID) // Ensure CategoryID is updated
	assert.Equal(t, float64(600), updatedBudget.Amount) // Ensure Amount is updated
	assert.Equal(t, float64(250), updatedBudget.Spent)  // Ensure Spent is updated
}

func TestUpdateBudgetWithEmptyCategory(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	budget := models.Budget{
		CategoryID: 0,
		Amount:     600.0,
		Spent:      250.0,
		StartDate:  time.Now(),
		EndDate:    time.Now().Add(30 * 24 * time.Hour),
	}
	_, err = models.UpdateBudget(db, budget)

	assert.Error(t, err)
	assert.Equal(t, errors.New("category id must be provided"), err)
}

// func TestDeleteBudget(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()

// 	mock.ExpectExec("DELETE FROM Budget WHERE category_id = \\$1").
// 		WithArgs(int64(1)).
// 		WillReturnResult(sqlmock.NewResult(1, 1))

// 	err = models.DeleteBudget(db, 1)

// 	assert.NoError(t, err)
// }
