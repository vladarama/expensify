package integration_test

import (
	"expense-tracker/internal/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestBudgetExpenseIntegration(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	// Test creating a category, budget, and expense workflow
	t.Run("Complete Budget-Expense Workflow", func(t *testing.T) {
		// 1. Create Category
		categoryRows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery("INSERT INTO Category").
			WithArgs("Groceries", "Food and household items").
			WillReturnRows(categoryRows)

		category := models.Category{
			Name:        "Groceries",
			Description: "Food and household items",
		}
		createdCategory, err := models.CreateCategory(db, category)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), createdCategory.ID)

		// 2. Create Budget for Category
		budgetRows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		startDate := time.Now()
		endDate := startDate.AddDate(0, 1, 0)

		mock.ExpectQuery("INSERT INTO Budget").
			WithArgs(createdCategory.ID, 500.0, 0.0, startDate, endDate).
			WillReturnRows(budgetRows)

		budget := models.Budget{
			CategoryID: createdCategory.ID,
			Amount:    500.0,
			Spent:     0.0,
			StartDate: startDate,
			EndDate:   endDate,
		}
		createdBudget, err := models.CreateBudget(db, budget)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), createdBudget.ID)

		// First expect the INSERT query
		now := time.Now()
		expenseRows := sqlmock.NewRows([]string{"id", "category_id", "amount", "date", "description"}).
			AddRow(1, createdCategory.ID, 100.0, now, "Weekly groceries")

		mock.ExpectQuery("INSERT INTO Expense \\(category_id, amount, date, description\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) RETURNING id, category_id, amount, date, description").
			WithArgs(createdCategory.ID, 100.0, sqlmock.AnyArg(), "Weekly groceries").
			WillReturnRows(expenseRows)

		// Then expect the budget check
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM Budget WHERE category_id = \\$1").
			WithArgs(createdCategory.ID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		// Finally expect the budget update
		mock.ExpectExec("UPDATE Budget SET spent = spent \\+ \\$1 WHERE category_id = \\$2").
			WithArgs(100.0, createdCategory.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		expense := models.Expense{
			CategoryID:  createdCategory.ID,
			Amount:     100.0,
			Date:       time.Now(),
			Description: "Weekly groceries",
		}
		createdExpense, err := models.CreateExpense(db, expense)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), createdExpense.ID)

		// 4. Verify Budget Update
		mock.ExpectQuery("SELECT id, category_id, amount, spent, start_date, end_date FROM Budget WHERE category_id = \\$1").
			WithArgs(createdCategory.ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "category_id", "amount", "spent", "start_date", "end_date"}).
				AddRow(1, createdCategory.ID, 500.0, 100.0, startDate, endDate))

		updatedBudget, err := models.GetBudgetByCategory(db, "1")
		assert.NoError(t, err)
		assert.Equal(t, float64(100.0), updatedBudget.Spent)
	})
} 