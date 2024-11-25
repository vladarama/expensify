const API_BASE_URL =
  "https://expense-tracker-golang-backend-eyhfeyfnhqexcxe2.canadacentral-01.azurewebsites.net";

export const API_ENDPOINTS = {
  budgets: {
    getAll: `${API_BASE_URL}/budgets`,
    create: `${API_BASE_URL}/budgets`,
    update: (id: number) => `${API_BASE_URL}/budgets/${id}`,
    delete: (id: number) => `${API_BASE_URL}/budgets/${id}`,
  },
  expenses: {
    getAll: `${API_BASE_URL}/expenses`,
    create: `${API_BASE_URL}/expenses`,
    update: (id: number) => `${API_BASE_URL}/expenses/${id}`,
    delete: (id: number) => `${API_BASE_URL}/expenses/${id}`,
  },
  incomes: {
    getAll: `${API_BASE_URL}/incomes`,
    create: `${API_BASE_URL}/incomes`,
    update: (id: number) => `${API_BASE_URL}/incomes/${id}`,
    delete: (id: number) => `${API_BASE_URL}/incomes/${id}`,
  },
  categories: {
    getAll: `${API_BASE_URL}/categories`,
    create: `${API_BASE_URL}/categories`,
    update: (id: number) => `${API_BASE_URL}/categories/${id}`,
    delete: (id: number) => `${API_BASE_URL}/categories/${id}`,
  },
};
