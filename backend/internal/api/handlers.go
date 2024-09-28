package api

import (
    "encoding/json"
    "net/http"

    "expense-tracker/internal/auth"
    "expense-tracker/internal/categories"
    "expense-tracker/internal/expenses"
    "expense-tracker/internal/users"
)

func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Name     string `json:"name"`
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil || req.Name == "" || req.Email == "" || req.Password == "" {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    user, err := users.RegisterUser(s.DB, req.Name, req.Email, req.Password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil || req.Email == "" || req.Password == "" {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    user, err := users.AuthenticateUser(s.DB, req.Email, req.Password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    token, err := auth.GenerateJWT(user.ID, s.JWTSecret)
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (s *Server) ExpensesHandler(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value(auth.UserIDKey).(int64)

    switch r.Method {
    case http.MethodGet:
        s.GetExpenses(w, r, userID)
    case http.MethodPost:
        s.CreateExpense(w, r, userID)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func (s *Server) GetExpenses(w http.ResponseWriter, r *http.Request, userID int64) {
    expensesList, err := expenses.GetExpenses(s.DB, userID)
    if err != nil {
        http.Error(w, "Failed to retrieve expenses", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(expensesList)
}

func (s *Server) CreateExpense(w http.ResponseWriter, r *http.Request, userID int64) {
    var expense expenses.Expense
    err := json.NewDecoder(r.Body).Decode(&expense)
    if err != nil || expense.Amount <= 0 || expense.Date.IsZero() {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    expense.UserID = userID

    err = expenses.CreateExpense(s.DB, &expense)
    if err != nil {
        http.Error(w, "Failed to create expense", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(expense)
}

func (s *Server) CategoriesHandler(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value(auth.UserIDKey).(int64)

    switch r.Method {
    case http.MethodGet:
        s.GetCategories(w, r, userID)
    case http.MethodPost:
        s.CreateCategory(w, r, userID)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func (s *Server) GetCategories(w http.ResponseWriter, r *http.Request, userID int64) {
    categoriesList, err := categories.GetCategories(s.DB, userID)
    if err != nil {
        http.Error(w, "Failed to retrieve categories", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(categoriesList)
}

func (s *Server) CreateCategory(w http.ResponseWriter, r *http.Request, userID int64) {
    var category categories.Category
    err := json.NewDecoder(r.Body).Decode(&category)
    if err != nil || category.Name == "" {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    category.UserID = userID

    err = categories.CreateCategory(s.DB, &category)
    if err != nil {
        http.Error(w, "Failed to create category", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(category)
}
