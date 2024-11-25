import { useState, useEffect, useMemo, useCallback } from "react";
import { Expense } from "../types/expense";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Trash2, Pencil } from "lucide-react";
import { ErrorBoundary } from "@/components/error-boundary";
import { Category } from "../types/category";
import { MonthlyExpenseChart } from "@/components/charts/monthly-expense-chart";
import { SortButton } from "@/components/ui/sort-button";
import { cn } from "@/lib/utils";
import { API_ENDPOINTS } from "@/config/api";

type SortField = "category" | "amount" | "date";
type SortDirection = "asc" | "desc" | null;

interface SortState {
  field: SortField | null;
  direction: SortDirection;
}

export function ExpensesPage() {
  return (
    <ErrorBoundary>
      <Expenses />
    </ErrorBoundary>
  );
}

export function Expenses() {
  const [expenses, setExpenses] = useState<Expense[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isOpen, setIsOpen] = useState(false);
  const [newExpense, setNewExpense] = useState({
    description: "",
    amount: "",
    date: new Date().toISOString().split("T")[0],
    category_id: "",
  });
  const [editingExpense, setEditingExpense] = useState<Expense | null>(null);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [categories, setCategories] = useState<Category[]>([]);
  const [sortState, setSortState] = useState<SortState>({
    field: null,
    direction: null,
  });

  useEffect(() => {
    let isMounted = true;

    async function fetchExpenses() {
      try {
        setIsLoading(true);
        const response = await fetch(API_ENDPOINTS.expenses.getAll);
        if (!response.ok)
          throw new Error(`HTTP error! status: ${response.status}`);
        const data = await response.json();
        if (isMounted) {
          setExpenses(data || []);
        }
      } catch (error) {
        console.error("Failed to fetch expenses:", error);
        if (isMounted) {
          setExpenses([]);
        }
      } finally {
        if (isMounted) {
          setIsLoading(false);
        }
      }
    }

    fetchExpenses();

    return () => {
      isMounted = false;
    };
  }, []);

  useEffect(() => {
    let isMounted = true;

    async function fetchCategories() {
      try {
        const response = await fetch(API_ENDPOINTS.categories.getAll);
        if (!response.ok)
          throw new Error(`HTTP error! status: ${response.status}`);
        const data = await response.json();
        if (isMounted) {
          setCategories(data || []);
        }
      } catch (error) {
        console.error("Failed to fetch categories:", error);
      }
    }

    fetchCategories();
    return () => {
      isMounted = false;
    };
  }, []);

  const getCategoryName = useCallback(
    (categoryId: number) => {
      const category = categories.find((cat) => cat.id === categoryId);
      return category?.name || "Unknown Category";
    },
    [categories]
  );

  const handleSort = (field: SortField) => {
    setSortState((prev) => ({
      field,
      direction:
        prev.field === field
          ? prev.direction === null
            ? "asc"
            : prev.direction === "asc"
            ? "desc"
            : null
          : "asc",
    }));
  };

  const sortedExpenses = useMemo(() => {
    if (!sortState.field || !sortState.direction) {
      return expenses;
    }

    return [...expenses].sort((a, b) => {
      if (sortState.field === "category") {
        const categoryA = getCategoryName(a.category_id).toLowerCase();
        const categoryB = getCategoryName(b.category_id).toLowerCase();
        return sortState.direction === "asc"
          ? categoryA.localeCompare(categoryB)
          : categoryB.localeCompare(categoryA);
      }

      if (sortState.field === "date") {
        const dateA = new Date(a.date).getTime();
        const dateB = new Date(b.date).getTime();
        return sortState.direction === "asc" ? dateA - dateB : dateB - dateA;
      }

      if (sortState.field === "amount") {
        return sortState.direction === "asc"
          ? a.amount - b.amount
          : b.amount - a.amount;
      }

      return 0; // Default case
    });
  }, [sortState.field, sortState.direction, expenses, getCategoryName]);

  if (isLoading) {
    return <div className="p-4">Loading expenses...</div>;
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    try {
      const formattedDate = new Date(newExpense.date).toISOString();
      const category = categories.find(
        (cat) => cat.name === newExpense.category_id
      );

      if (!category) {
        console.error("Category not found");
        return;
      }

      const response = await fetch(API_ENDPOINTS.expenses.create, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          description: newExpense.description,
          amount: parseFloat(newExpense.amount),
          date: formattedDate,
          category_id: category.id,
        }),
      });

      if (response.ok) {
        const data = await response.json();
        setExpenses([...expenses, data]);
        setIsOpen(false);
        setNewExpense({
          description: "",
          amount: "",
          date: new Date().toISOString().split("T")[0],
          category_id: "",
        });
      } else {
        const errorData = await response.text();
        console.error("Failed to add expense:", errorData);
      }
    } catch (error) {
      console.error("Failed to add expense:", error);
    }
  }

  async function handleDelete(id: number) {
    try {
      const response = await fetch(API_ENDPOINTS.expenses.delete(id), {
        method: "DELETE",
      });

      if (response.ok) {
        setExpenses(expenses.filter((expense) => expense.id !== id));
      }
    } catch (error) {
      console.error("Failed to delete expense:", error);
    }
  }

  async function handleEdit(e: React.FormEvent) {
    e.preventDefault();
    if (!editingExpense) return;

    try {
      const formattedDate = new Date(editingExpense.date).toISOString();

      const response = await fetch(
        API_ENDPOINTS.expenses.update(editingExpense.id),
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            amount: editingExpense.amount,
            date: formattedDate,
            category_id: editingExpense.category_id,
          }),
        }
      );

      if (response.ok) {
        const updatedExpense = await response.json();
        setExpenses(
          expenses.map((expense) =>
            expense.id === editingExpense.id ? updatedExpense : expense
          )
        );
        setIsEditDialogOpen(false);
        setEditingExpense(null);
      }
    } catch (error) {
      console.error("Failed to update expense:", error);
    }
  }

  function handleStartEdit(expense: Expense) {
    setEditingExpense({
      ...expense,
      date: new Date(expense.date).toISOString().split("T")[0],
    });
    setIsEditDialogOpen(true);
  }

  return (
    <div className="min-h-screen p-4 sm:p-6">
      <div className="max-w-7xl mx-auto space-y-6 sm:space-y-8">
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 sm:gap-0">
          <h1 className="text-2xl sm:text-3xl font-bold bg-gradient-to-r from-red-600 to-rose-500 bg-clip-text text-transparent">
            Expenses Dashboard
          </h1>
          <Dialog open={isOpen} onOpenChange={setIsOpen}>
            <DialogTrigger asChild>
              <Button
                className={cn(
                  "w-full sm:w-auto",
                  "bg-gradient-to-r from-red-600 to-rose-500",
                  "hover:from-red-700 hover:to-rose-600",
                  "text-white shadow-md hover:shadow-lg",
                  "transition-all duration-200"
                )}
              >
                + Add Expense
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Add New Expense</DialogTitle>
              </DialogHeader>
              <form onSubmit={handleSubmit} className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="description">Description</Label>
                  <Input
                    id="description"
                    type="text"
                    required
                    value={newExpense.description}
                    onChange={(e) =>
                      setNewExpense({
                        ...newExpense,
                        description: e.target.value,
                      })
                    }
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="amount">Amount</Label>
                  <Input
                    id="amount"
                    type="number"
                    step="0.01"
                    required
                    value={newExpense.amount}
                    onChange={(e) =>
                      setNewExpense({ ...newExpense, amount: e.target.value })
                    }
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="date">Date</Label>
                  <Input
                    id="date"
                    type="date"
                    required
                    value={newExpense.date}
                    onChange={(e) =>
                      setNewExpense({ ...newExpense, date: e.target.value })
                    }
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="category_id">Category</Label>
                  <select
                    id="category_id"
                    className="w-full rounded-md border border-input bg-background px-3 py-2"
                    value={newExpense.category_id}
                    onChange={(e) =>
                      setNewExpense({
                        ...newExpense,
                        category_id: e.target.value,
                      })
                    }
                    required
                  >
                    <option value="">Select a category</option>
                    {categories.map((category) => (
                      <option key={category.id} value={category.name}>
                        {category.name}
                      </option>
                    ))}
                  </select>
                </div>
                <Button
                  type="submit"
                  className={cn(
                    "w-full",
                    "bg-gradient-to-r from-red-600 to-rose-500",
                    "hover:from-red-700 hover:to-rose-600",
                    "text-white shadow-md hover:shadow-lg",
                    "transition-all duration-200"
                  )}
                >
                  Add Expense
                </Button>
              </form>
            </DialogContent>
          </Dialog>
        </div>

        <div className="grid grid-cols-1 xl:grid-cols-2 gap-4 sm:gap-8">
          <div className="bg-white rounded-xl shadow-lg p-4 sm:p-6 overflow-hidden flex flex-col min-h-[500px] xl:h-[600px]">
            <h2 className="text-lg sm:text-xl font-semibold mb-4 text-gray-800">
              Expense History
            </h2>
            <div className="overflow-x-auto">
              <div className="space-y-4">
                {sortedExpenses.map((expense) => (
                  <div
                    key={expense.id}
                    className="group bg-white p-4 rounded-lg shadow-sm hover:bg-red-50/50 transition-colors lg:hidden"
                  >
                    <div className="flex justify-between items-start mb-2">
                      <div>
                        <div className="font-medium text-gray-900">
                          {expense.description}
                        </div>
                        <div className="text-gray-600 text-sm mt-1">
                          {getCategoryName(expense.category_id)}
                        </div>
                        <div className="text-red-600 font-semibold mt-1">
                          ${expense.amount.toLocaleString()}
                        </div>
                      </div>
                      <div className="flex gap-2">
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => handleStartEdit(expense)}
                          className="h-8 w-8 hover:bg-red-100 rounded-full"
                        >
                          <Pencil className="h-4 w-4 text-red-600" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => handleDelete(expense.id)}
                          className="h-8 w-8 hover:bg-red-100 rounded-full"
                        >
                          <Trash2 className="h-4 w-4 text-red-500" />
                        </Button>
                      </div>
                    </div>
                    <div className="text-gray-600 text-sm">
                      {new Date(expense.date).toLocaleDateString()}
                    </div>
                  </div>
                ))}

                <div className="hidden lg:block">
                  <table className="w-full">
                    <thead>
                      <tr className="border-b border-gray-200">
                        <th className="text-left py-3 px-4 font-semibold text-gray-600 w-1/5">
                          Description
                        </th>
                        <th className="text-left py-3 px-4 font-semibold text-gray-600 w-1/5">
                          Category
                        </th>
                        <th className="text-left py-3 px-4 font-semibold text-gray-600 w-1/5">
                          <SortButton
                            label="Amount"
                            active={sortState.field === "amount"}
                            direction={
                              sortState.field === "amount"
                                ? sortState.direction
                                : null
                            }
                            onClick={() => handleSort("amount")}
                            className="hover:text-red-600 hover:bg-transparent"
                          />
                        </th>
                        <th className="text-left py-3 px-4 font-semibold text-gray-600 w-2/5 min-w-[200px]">
                          <div className="flex justify-between items-center">
                            <SortButton
                              label="Date"
                              active={sortState.field === "date"}
                              direction={
                                sortState.field === "date"
                                  ? sortState.direction
                                  : null
                              }
                              onClick={() => handleSort("date")}
                              className="hover:text-red-600 hover:bg-transparent"
                            />
                          </div>
                        </th>
                      </tr>
                    </thead>
                    <tbody>
                      {sortedExpenses.map((expense) => (
                        <tr
                          key={expense.id}
                          className="group border-b border-gray-100 last:border-none hover:bg-red-50/50 transition-colors"
                        >
                          <td className="py-3 px-4">
                            <span className="font-medium text-gray-900">
                              {expense.description}
                            </span>
                          </td>
                          <td className="py-3 px-4">
                            <span className="text-gray-600">
                              {getCategoryName(expense.category_id)}
                            </span>
                          </td>
                          <td className="py-3 px-4">
                            <span className="text-red-600 font-semibold">
                              ${expense.amount.toLocaleString()}
                            </span>
                          </td>
                          <td className="py-3 px-4">
                            <div className="flex justify-between items-center gap-2">
                              <span className="text-gray-600 text-sm sm:text-base">
                                {new Date(expense.date).toLocaleDateString()}
                              </span>
                              <div className="flex gap-1 sm:gap-2 transition-opacity">
                                <Button
                                  variant="ghost"
                                  size="icon"
                                  onClick={() => handleStartEdit(expense)}
                                  className="h-7 w-7 sm:h-8 sm:w-8 hover:bg-red-100 rounded-full"
                                >
                                  <Pencil className="h-3 w-3 sm:h-4 sm:w-4 text-red-600" />
                                </Button>
                                <Button
                                  variant="ghost"
                                  size="icon"
                                  onClick={() => handleDelete(expense.id)}
                                  className="h-7 w-7 sm:h-8 sm:w-8 hover:bg-red-100 rounded-full"
                                >
                                  <Trash2 className="h-3 w-3 sm:h-4 sm:w-4 text-red-500" />
                                </Button>
                              </div>
                            </div>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
          </div>

          <div className="bg-white rounded-xl shadow-lg p-4 sm:p-6 flex flex-col min-h-[500px] xl:h-[600px]">
            <div className="flex-1">
              <MonthlyExpenseChart expenses={expenses} />
            </div>
          </div>
        </div>

        <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Edit Expense</DialogTitle>
            </DialogHeader>
            <form onSubmit={handleEdit} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="edit-name">Name</Label>
                <Input
                  id="edit-name"
                  type="text"
                  required
                  value={editingExpense?.description || ""}
                  onChange={(e) =>
                    setEditingExpense((prev) =>
                      prev ? { ...prev, description: e.target.value } : null
                    )
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="edit-amount">Amount</Label>
                <Input
                  id="edit-amount"
                  type="number"
                  step="0.01"
                  required
                  value={editingExpense?.amount || ""}
                  onChange={(e) =>
                    setEditingExpense((prev) =>
                      prev
                        ? { ...prev, amount: parseFloat(e.target.value) }
                        : null
                    )
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="edit-date">Date</Label>
                <Input
                  id="edit-date"
                  type="date"
                  required
                  value={editingExpense?.date || ""}
                  onChange={(e) =>
                    setEditingExpense((prev) =>
                      prev ? { ...prev, date: e.target.value } : null
                    )
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="edit-category">Category</Label>
                <select
                  id="edit-category"
                  className="w-full rounded-md border border-input bg-background px-3 py-2"
                  value={getCategoryName(editingExpense?.category_id || 0)}
                  onChange={(e) => {
                    const category = categories.find(
                      (cat) => cat.name === e.target.value
                    );
                    if (category && editingExpense) {
                      setEditingExpense({
                        ...editingExpense,
                        category_id: category.id,
                      });
                    }
                  }}
                  required
                >
                  <option value="">Select a category</option>
                  {categories.map((category) => (
                    <option key={category.id} value={category.name}>
                      {category.name}
                    </option>
                  ))}
                </select>
              </div>
              <Button
                type="submit"
                className={cn(
                  "w-full",
                  "bg-gradient-to-r from-red-600 to-rose-500",
                  "hover:from-red-700 hover:to-rose-600",
                  "text-white shadow-md hover:shadow-lg",
                  "transition-all duration-200"
                )}
              >
                Save Changes
              </Button>
            </form>
          </DialogContent>
        </Dialog>
      </div>
    </div>
  );
}
