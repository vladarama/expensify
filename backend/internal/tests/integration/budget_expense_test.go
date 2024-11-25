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

	t.Run("Complete Budget-Expense Workflow", func(t *testing.T) {
		// Step 1: Create Category
		mock.ExpectQuery(`INSERT INTO Category \(name, description\) VALUES \(\$1, \$2\) RETURNING id`).
			WithArgs("Groceries", "Food and household items").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		category := models.Category{Name: "Groceries", Description: "Food and household items"}
		createdCategory, err := models.CreateCategory(db, category)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), createdCategory.ID)

		// Step 2: Create Budget
		startDate := time.Now().AddDate(0, -1, 0)
		endDate := time.Now().AddDate(0, 1, 0)

		// Mock budget overlap check
		mock.ExpectQuery(`SELECT EXISTS \( SELECT 1 FROM Budget WHERE category_id = \$1 AND id <> \$4 AND \( \(start_date <= \$3 AND end_date >= \$2\) \) \)`).
			WithArgs(createdCategory.ID, startDate, endDate, 0).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		// Mock total spent calculation (initially 0)
		mock.ExpectQuery(`SELECT COALESCE\(SUM\(amount\), 0\) FROM Expense WHERE category_id = \$1 AND date >= \$2 AND date <= \$3`).
			WithArgs(createdCategory.ID, startDate, endDate).
			WillReturnRows(sqlmock.NewRows([]string{"sum"}).AddRow(0.0))

		// Mock budget creation
		mock.ExpectQuery(`INSERT INTO Budget \(category_id, amount, spent, start_date, end_date\) VALUES \(\$1, \$2, \$3, \$4, \$5\) RETURNING id`).
			WithArgs(createdCategory.ID, 500.0, 0.0, startDate, endDate).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		budget := models.Budget{
			CategoryID: createdCategory.ID,
			Amount:     500.0,
			StartDate:  startDate,
			EndDate:    endDate,
		}
		createdBudget, err := models.CreateBudget(db, budget)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), createdBudget.ID)
		assert.Equal(t, float64(0.0), createdBudget.Spent)

		// Step 3: Create Expense
		mock.ExpectQuery(`INSERT INTO Expense \(category_id, amount, date, description\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING id, category_id, amount, date, description`).
			WithArgs(createdCategory.ID, 100.0, sqlmock.AnyArg(), "Weekly groceries").
			WillReturnRows(sqlmock.NewRows([]string{"id", "category_id", "amount", "date", "description"}).
				AddRow(1, createdCategory.ID, 100.0, time.Now(), "Weekly groceries"))

		// Mock budget existence check
		mock.ExpectQuery(`SELECT COUNT\(\*\) FROM Budget WHERE category_id = \$1`).
			WithArgs(createdCategory.ID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		// Mock budget spent update
		mock.ExpectExec(`UPDATE Budget SET spent = spent \+ \$1 WHERE category_id = \$2`).
			WithArgs(100.0, createdCategory.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		expense := models.Expense{
			CategoryID:  createdCategory.ID,
			Amount:      100.0,
			Date:        time.Now(),
			Description: "Weekly groceries",
		}
		createdExpense, err := models.CreateExpense(db, expense)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), createdExpense.ID)

		// Step 4: Verify Budget Update
		mock.ExpectQuery(`SELECT id, category_id, amount, spent, start_date, end_date FROM Budget WHERE category_id = \$1`).
			WithArgs(createdCategory.ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "category_id", "amount", "spent", "start_date", "end_date"}).
				AddRow(1, createdCategory.ID, 500.0, 100.0, startDate, endDate))

		updatedBudgets, err := models.GetBudgetsByCategoryID(db, createdCategory.ID)
		assert.NoError(t, err)
		assert.Len(t, updatedBudgets, 1)
		assert.Equal(t, float64(100.0), updatedBudgets[0].Spent)
	})

	// Verify all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
