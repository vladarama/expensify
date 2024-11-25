package integration_test

import (
	"database/sql"
	"errors"
	"expense-tracker/internal/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCategoryBudgetIntegration(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	t.Run("Category Deletion Cascade", func(t *testing.T) {
		// 1. Setup - Create Category with Budget
		categoryRows := sqlmock.NewRows([]string{"id"}).AddRow(2)
		mock.ExpectQuery(`INSERT INTO Category \(name, description\) VALUES \(\$1, \$2\) RETURNING id`).
			WithArgs("Entertainment", "Entertainment expenses").
			WillReturnRows(categoryRows)

		category := models.Category{
			Name:        "Entertainment",
			Description: "Entertainment expenses",
		}
		createdCategory, err := models.CreateCategory(db, category)
		assert.NoError(t, err)

		// Mock budget overlap check
		startDate := time.Now()
		endDate := startDate.AddDate(0, 1, 0)
		mock.ExpectQuery(`SELECT EXISTS \( SELECT 1 FROM Budget WHERE category_id = \$1 AND id <> \$4 AND \( \(start_date <= \$3 AND end_date >= \$2\) \) \)`).
			WithArgs(createdCategory.ID, startDate, endDate, 0).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		// Mock total spent calculation
		mock.ExpectQuery(`SELECT COALESCE\(SUM\(amount\), 0\) FROM Expense WHERE category_id = \$1 AND date >= \$2 AND date <= \$3`).
			WithArgs(createdCategory.ID, startDate, endDate).
			WillReturnRows(sqlmock.NewRows([]string{"sum"}).AddRow(0.0))

		// Mock budget creation
		budgetRows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery(`INSERT INTO Budget \(category_id, amount, spent, start_date, end_date\) VALUES \(\$1, \$2, \$3, \$4, \$5\) RETURNING id`).
			WithArgs(createdCategory.ID, 300.0, 0.0, startDate, endDate).
			WillReturnRows(budgetRows)

		budget := models.Budget{
			CategoryID: createdCategory.ID,
			Amount:     300.0,
			StartDate:  startDate,
			EndDate:    endDate,
		}
		_, err = models.CreateBudget(db, budget)
		assert.NoError(t, err)

		// 2. Delete Category (should cascade to budget)
		mock.ExpectExec(`UPDATE Expense SET category_id = 1 WHERE category_id = \$1`).
			WithArgs(createdCategory.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec(`DELETE FROM Category WHERE id = \$1`).
			WithArgs(createdCategory.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = models.DeleteCategory(db, createdCategory.ID)
		assert.NoError(t, err)

		// 3. Verify Budget is deleted (should return no rows)
		mock.ExpectQuery(`SELECT id, category_id, amount, spent, start_date, end_date FROM Budget WHERE category_id = \$1`).
			WithArgs(createdCategory.ID).
			WillReturnError(sql.ErrNoRows)

		_, err = models.GetBudgetsByCategoryID(db, createdCategory.ID)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, sql.ErrNoRows), "Expected sql.ErrNoRows, got: %v", err)
	})
}
