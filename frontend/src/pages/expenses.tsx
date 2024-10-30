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
    name: "",
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
        const response = await fetch("http://localhost:8080/expenses");
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
        const response = await fetch("http://localhost:8080/categories");
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

      const response = await fetch("http://localhost:8080/expenses", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          name: newExpense.name,
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
          name: "",
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
      const response = await fetch(`http://localhost:8080/expenses/${id}`, {
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
        `http://localhost:8080/expenses/${editingExpense.id}`,
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
    <div>
      <h1 className="text-2xl font-bold mb-4">Expenses</h1>

      <Dialog open={isOpen} onOpenChange={setIsOpen}>
        <DialogTrigger asChild>
          <Button className="mb-4">+ Add Expense</Button>
        </DialogTrigger>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Add New Expense</DialogTitle>
          </DialogHeader>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="name">Name</Label>
              <Input
                id="name"
                type="text"
                required
                value={newExpense.name}
                onChange={(e) =>
                  setNewExpense({ ...newExpense, name: e.target.value })
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
                  setNewExpense({ ...newExpense, category_id: e.target.value })
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
            <Button type="submit" className="w-full">
              Add Expense
            </Button>
          </form>
        </DialogContent>
      </Dialog>

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
                value={editingExpense?.name || ""}
                onChange={(e) =>
                  setEditingExpense((prev) =>
                    prev ? { ...prev, name: e.target.value } : null
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
            <Button type="submit" className="w-full">
              Save Changes
            </Button>
          </form>
        </DialogContent>
      </Dialog>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Table Section */}
        <div className="overflow-auto">
          <table className="min-w-full bg-white">
            <thead>
              <tr>
                <th className="border px-4 py-2">Name</th>
                <th className="border px-4 py-2">
                  <SortButton
                    label="Category"
                    active={sortState.field === "category"}
                    direction={
                      sortState.field === "category"
                        ? sortState.direction
                        : null
                    }
                    onClick={() => handleSort("category")}
                  />
                </th>
                <th className="border px-4 py-2">
                  <SortButton
                    label="Amount"
                    active={sortState.field === "amount"}
                    direction={
                      sortState.field === "amount" ? sortState.direction : null
                    }
                    onClick={() => handleSort("amount")}
                  />
                </th>
                <th className="border px-4 py-2">
                  <SortButton
                    label="Date"
                    active={sortState.field === "date"}
                    direction={
                      sortState.field === "date" ? sortState.direction : null
                    }
                    onClick={() => handleSort("date")}
                  />
                </th>
              </tr>
            </thead>
            <tbody>
              {sortedExpenses?.length > 0 ? (
                sortedExpenses.map((expense) => (
                  <tr key={expense.id} className="group hover:bg-gray-50">
                    <td className="border px-4 py-2">{expense.name}</td>
                    <td className="border px-4 py-2">
                      {getCategoryName(expense.category_id)}
                    </td>
                    <td className="border px-4 py-2">${expense.amount}</td>
                    <td className="border px-4 py-2 relative min-w-[200px]">
                      <div className="flex justify-between items-center gap-2">
                        <span>
                          {new Date(expense.date).toLocaleDateString()}
                        </span>
                        <div className="flex gap-2 opacity-0 group-hover:opacity-100">
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => handleStartEdit(expense)}
                            className="h-8 w-8 hover:bg-gray-100/50"
                          >
                            <Pencil className="h-4 w-4 text-blue-500" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => handleDelete(expense.id)}
                            className="h-8 w-8 hover:bg-gray-100/50"
                          >
                            <Trash2 className="h-4 w-4 text-red-500" />
                          </Button>
                        </div>
                      </div>
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td
                    colSpan={4}
                    className="border px-4 py-2 text-center text-gray-500"
                  >
                    No expenses found
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>

        {/* Charts Section */}
        <div>
          <MonthlyExpenseChart expenses={expenses} />
        </div>
      </div>
    </div>
  );
}