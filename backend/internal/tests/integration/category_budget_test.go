package integration_test

import (
	"database/sql"
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
		categoryRows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery("INSERT INTO Category").
			WithArgs("Entertainment", "Entertainment expenses").
			WillReturnRows(categoryRows)

		category := models.Category{
			Name:        "Entertainment",
			Description: "Entertainment expenses",
		}
		createdCategory, err := models.CreateCategory(db, category)
		assert.NoError(t, err)

		// Create associated budget
		budgetRows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		startDate := time.Now()
		endDate := startDate.AddDate(0, 1, 0)

		mock.ExpectQuery("INSERT INTO Budget").
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
		mock.ExpectExec("DELETE FROM Category WHERE id = \\$1").
			WithArgs(createdCategory.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = models.DeleteCategory(db, createdCategory.ID)
		assert.NoError(t, err)

		// 3. Verify Budget is deleted (should return no rows)
		mock.ExpectQuery("SELECT (.+) FROM Budget WHERE category_id = \\$1").
			WithArgs(createdCategory.ID).
			WillReturnError(sql.ErrNoRows)

		_, err = models.GetBudgetsByCategoryID(db, 1)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
	})
}
