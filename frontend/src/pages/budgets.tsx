import { useState, useEffect, useMemo } from "react";
import { Budget } from "../types/budget";
import { Category } from "../types/category";
import { Button } from "@/components/ui/button";
import { MonthWiseBudgetChart } from "@/components/charts/month-budget-chart";
import { CategoryWiseBudgetChart } from "@/components/charts/category-budget-chart";

import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { Trash2, Pencil } from "lucide-react";
import { ErrorBoundary } from "@/components/error-boundary";
import { SortButton } from "@/components/ui/sort-button";
import DatePicker from "react-datepicker";
import "react-datepicker/dist/react-datepicker.css";

// Utility function to format the month-year date
function formatMonth(date: Date): string {
  return date.toLocaleDateString(undefined, {
    month: "long",
    year: "numeric",
  });
}

type SortField = "category" | "amount" | "startDate";
type SortDirection = "asc" | "desc" | null;

interface SortState {
  field: SortField | null;
  direction: SortDirection;
}

export function BudgetsPage() {
  return (
    <ErrorBoundary>
      <Budgets />
    </ErrorBoundary>
  );
}

export function Budgets() {
  const [budgets, setBudgets] = useState<Budget[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [isAddDialogOpen, setIsAddDialogOpen] = useState(false);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [newBudget, setNewBudget] = useState({
    category_id: "",
    amount: "",
    month: new Date(),
  });
  const [editingBudget, setEditingBudget] = useState<Budget | null>(null);
  const [sortState, setSortState] = useState<SortState>({
    field: null,
    direction: null,
  });
  const [error, setError] = useState<string | null>(null); // For error messages in modal
  const [errorGlobal, setErrorGlobal] = useState<string | null>(null); // Global error
  const [selectedChartMonth, setSelectedChartMonth] = useState(new Date());
  const [selectedCategory, setSelectedCategory] = useState<number | null>(null);

  useEffect(() => {
    async function fetchBudgets() {
      try {
        setErrorGlobal(null); // Clear previous errors
        const response = await fetch("http://localhost:8080/budgets");
        if (!response.ok) throw new Error(await response.text());
        const data = await response.json();
        setBudgets(data || []);
      } catch (err: any) {
        setErrorGlobal(err.message);
        console.error("Failed to fetch budgets:", err.message);
      }
    }

    async function fetchCategories() {
      try {
        setErrorGlobal(null); // Clear previous errors
        const response = await fetch("http://localhost:8080/categories");
        if (!response.ok) throw new Error(await response.text());
        const data = await response.json();
        setCategories(data || []);
      } catch (err: any) {
        setErrorGlobal(err.message);
        console.error("Failed to fetch categories:", err.message);
      }
    }

    fetchBudgets();
    fetchCategories();
  }, []);

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

  const sortedBudgets = useMemo(() => {
    if (!sortState.field || !sortState.direction) return budgets;

    return [...budgets].sort((a, b) => {
      if (sortState.field === "category") {
        const categoryA = categories.find((c) => c.id === a.category_id)?.name || "";
        const categoryB = categories.find((c) => c.id === b.category_id)?.name || "";
        return sortState.direction === "asc"
          ? categoryA.localeCompare(categoryB)
          : categoryB.localeCompare(categoryA);
      }

      if (sortState.field === "amount") {
        return sortState.direction === "asc"
          ? a.amount - b.amount
          : b.amount - a.amount;
      }

      if (sortState.field === "startDate") {
        const dateA = new Date(a.start_date).getTime();
        const dateB = new Date(b.start_date).getTime();
        return sortState.direction === "asc" ? dateA - dateB : dateB - dateA;
      }

      return 0;
    });
  }, [sortState, budgets, categories]);

  async function handleAddBudget(e: React.FormEvent) {
    e.preventDefault();
    setError(null); // Clear modal-specific errors
    const startDate = new Date(newBudget.month.getFullYear(), newBudget.month.getMonth(), 1).toISOString();
    const endDate = new Date(newBudget.month.getFullYear(), newBudget.month.getMonth() + 1, 0).toISOString();
    try {
      const response = await fetch("http://localhost:8080/budgets", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          category_id: parseInt(newBudget.category_id, 10),
          amount: parseFloat(newBudget.amount),
          start_date: startDate,
          end_date: endDate,
        }),
      });

      if (!response.ok) throw new Error(await response.text());
      const data = await response.json();
      setBudgets([...budgets, data]);
      setIsAddDialogOpen(false);
      setNewBudget({ category_id: "", amount: "", month: new Date() });
    } catch (err: any) {
      setError(err.message);
      console.error("Failed to add budget:", err.message);
    }
  }

  async function handleEditBudget(e: React.FormEvent) {
    e.preventDefault();
    setError(null); // Clear modal-specific errors
    if (!editingBudget) return;

    try {
      const response = await fetch(`http://localhost:8080/budgets/${editingBudget.id}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          category_id: editingBudget.category_id,
          amount: editingBudget.amount,
          start_date: new Date(editingBudget.start_date).toISOString(),
          end_date: new Date(editingBudget.end_date).toISOString(),
        }),
      });

      if (!response.ok) throw new Error(await response.text());
      const updatedBudget = await response.json();
      setBudgets(
        budgets.map((budget) =>
          budget.id === editingBudget.id ? updatedBudget : budget
        )
      );
      setIsEditDialogOpen(false);
      setEditingBudget(null);
    } catch (err: any) {
      setError(err.message);
      console.error("Failed to edit budget:", err.message);
    }
  }

  async function handleDeleteBudget(id: number) {
    setErrorGlobal(null); // Clear global errors
    try {
      const response = await fetch(`http://localhost:8080/budgets/${id}`, { method: "DELETE" });
      if (!response.ok) throw new Error(await response.text());
      setBudgets(budgets.filter((budget) => budget.id !== id));
    } catch (err: any) {
      setErrorGlobal(err.message);
      console.error("Failed to delete budget:", err.message);
    }
  }

  return (
    <div>
      <h1 className="text-2xl font-bold mb-4">Budgets</h1>

      {errorGlobal && <div className="mb-4 text-red-600">{errorGlobal}</div>}

      {/* Add Budget Dialog */}
      <Dialog open={isAddDialogOpen} onOpenChange={setIsAddDialogOpen}>
        <DialogTrigger asChild>
          <Button className="mb-4">+ Add Budget</Button>
        </DialogTrigger>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Add New Budget</DialogTitle>
          </DialogHeader>
          <form onSubmit={handleAddBudget} className="space-y-4">
            {error && <div className="text-red-600">{error}</div>}
            <div className="space-y-2">
              <Label htmlFor="category">Category</Label>
              <select
                id="category"
                className="w-full rounded-md border border-input bg-background px-3 py-2"
                value={newBudget.category_id}
                onChange={(e) => setNewBudget({ ...newBudget, category_id: e.target.value })}
                required
              >
                <option value="">Select a category</option>
                {categories.map((category) => (
                  <option key={category.id} value={category.id}>
                    {category.name}
                  </option>
                ))}
              </select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="amount">Amount</Label>
              <input
                id="amount"
                type="number"
                className="w-full rounded-md border border-input px-3 py-2"
                step="0.01"
                required
                value={newBudget.amount}
                onChange={(e) => setNewBudget({ ...newBudget, amount: e.target.value })}
              />
            </div>
            <div className="space-y-2">
            <Label htmlFor="month">Month</Label>
            <DatePicker
                id="month"
                selected={newBudget.month}
                onChange={(date) => setNewBudget({ ...newBudget, month: date! })}
                dateFormat="MM/yyyy"
                showMonthYearPicker
                required
                className="w-full rounded-md border border-input px-3 py-2"
            />
            <div className="text-sm text-gray-500 mt-1">
                Select the month and year for this budget.
            </div>
            </div>

            <Button type="submit" className="w-full">
              Add Budget
            </Button>
          </form>
        </DialogContent>
      </Dialog>

      {/* Edit Budget Dialog */}
      <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Edit Budget</DialogTitle>
          </DialogHeader>
          <form onSubmit={handleEditBudget} className="space-y-4">
            {error && <div className="text-red-600">{error}</div>}
            <div className="space-y-2">
              <Label htmlFor="edit-category">Category</Label>
              <select
                id="edit-category"
                className="w-full rounded-md border border-input px-3 py-2"
                value={editingBudget?.category_id || ""}
                onChange={(e) =>
                  setEditingBudget((prev) =>
                    prev ? { ...prev, category_id: parseInt(e.target.value, 10) } : null
                  )
                }
                required
              >
                <option value="">Select a category</option>
                {categories.map((category) => (
                  <option key={category.id} value={category.id}>
                    {category.name}
                  </option>
                ))}
              </select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="edit-amount">Amount</Label>
              <input
                id="edit-amount"
                type="number"
                className="w-full rounded-md border border-input px-3 py-2"
                step="0.01"
                required
                value={editingBudget?.amount || ""}
                onChange={(e) =>
                  setEditingBudget((prev) =>
                    prev ? { ...prev, amount: parseFloat(e.target.value) } : null
                  )
                }
              />
            </div>
            <Button type="submit" className="w-full">
              Save Changes
            </Button>
          </form>
        </DialogContent>
      </Dialog>

      <div className="flex flex-col lg:flex-row gap-8">
      {/* Budgets Table */}
        <div className="w-full lg:w-1/2">
            <table className="min-w-full bg-white">
            <thead>
                <tr>
                <th className="border px-4 py-2">Category</th>
                <th className="border px-4 py-2">
                    <SortButton
                    label="Amount"
                    active={sortState.field === "amount"}
                    direction={sortState.field === "amount" ? sortState.direction : null}
                    onClick={() => handleSort("amount")}
                    />
                </th>
                <th className="border px-4 py-2">
                    <SortButton
                    label="Month"
                    active={sortState.field === "startDate"}
                    direction={sortState.field === "startDate" ? sortState.direction : null}
                    onClick={() => handleSort("startDate")}
                    />
                </th>
                </tr>
            </thead>
            <tbody>
                {sortedBudgets.map((budget) => (
                <tr key={budget.id} className="group hover:bg-gray-50">
                    <td className="border px-4 py-2">
                    {categories.find((c) => c.id === budget.category_id)?.name || "Unknown"}
                    </td>
                    <td className="border px-4 py-2">${budget.amount}</td>
                    <td className="border px-4 py-2 relative">
                    <div className="flex justify-between items-center gap-2">
                        <span>{formatMonth(new Date(budget.start_date))}</span>
                        <div className="flex gap-2 opacity-0 group-hover:opacity-100">
                        <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => {
                            setEditingBudget(budget);
                            setIsEditDialogOpen(true);
                            }}
                            className="h-8 w-8 hover:bg-gray-100/50"
                        >
                            <Pencil className="h-4 w-4 text-blue-500" />
                        </Button>
                        <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => handleDeleteBudget(budget.id)}
                            className="h-8 w-8 hover:bg-gray-100/50"
                        >
                            <Trash2 className="h-4 w-4 text-red-500" />
                        </Button>
                        </div>
                    </div>
                    </td>
                </tr>
                ))}
            </tbody>
            </table>
        </div>

        {/* Charts */}
        <div className="w-full lg:w-1/2">
            {/* Month-Wise Budget Chart */}
            <div className="mb-8">
            <h2 className="text-lg font-semibold mb-2"> Budgets by Month</h2>
            <DatePicker
                selected={selectedChartMonth}
                onChange={(date) => setSelectedChartMonth(date!)}
                dateFormat="MM/yyyy"
                showMonthYearPicker
                className="mb-4 border rounded p-2"
            />
            <MonthWiseBudgetChart budgets={budgets} selectedMonth={selectedChartMonth}  found_categories={categories}/>
            </div>

            {/* Category-Wise Budget Chart */}
            <div>
            <h2 className="text-lg font-semibold mb-2">Budget by Category</h2>
            <select
                value={selectedCategory || ""}
                onChange={(e) => setSelectedCategory(Number(e.target.value))}
                className="mb-4 border rounded p-2 w-full"
            >
                <option value="">Select a category</option>
                {categories.map((category) => (
                <option key={category.id} value={category.id}>
                    {category.name}
                </option>
                ))}
            </select>
            {selectedCategory && (
            <CategoryWiseBudgetChart
                budgets={budgets}
                selectedCategoryId={selectedCategory} 
                found_categories={categories} 
            />
            )}

            </div>
        </div>
        </div>

      {/* Budgets Table
      <table className="min-w-full bg-white">
        <thead>
          <tr>
            <th className="border px-4 py-2">Category</th>
            <th className="border px-4 py-2">
              <SortButton
                label="Amount"
                active={sortState.field === "amount"}
                direction={sortState.field === "amount" ? sortState.direction : null}
                onClick={() => handleSort("amount")}
              />
            </th>
            <th className="border px-4 py-2">
              <SortButton
                label="Month"
                active={sortState.field === "startDate"}
                direction={sortState.field === "startDate" ? sortState.direction : null}
                onClick={() => handleSort("startDate")}
              />
            </th>
          </tr>
        </thead>
        <tbody>
          {sortedBudgets.map((budget) => (
            <tr key={budget.id} className="group hover:bg-gray-50">
              <td className="border px-4 py-2">
                {categories.find((c) => c.id === budget.category_id)?.name || "Unknown"}
              </td>
              <td className="border px-4 py-2">${budget.amount}</td>
              <td className="border px-4 py-2 relative">
                <div className="flex justify-between items-center gap-2">
                  <span>{formatMonth(new Date(budget.start_date))}</span>
                  <div className="flex gap-2 opacity-0 group-hover:opacity-100">
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => {
                        setEditingBudget(budget);
                        setIsEditDialogOpen(true);
                      }}
                      className="h-8 w-8 hover:bg-gray-100/50"
                    >
                      <Pencil className="h-4 w-4 text-blue-500" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => handleDeleteBudget(budget.id)}
                      className="h-8 w-8 hover:bg-gray-100/50"
                    >
                      <Trash2 className="h-4 w-4 text-red-500" />
                    </Button>
                  </div>
                </div>
              </td>
            </tr>
          ))}
        </tbody>
        </table>

        <div className="mt-8">
        <h2 className="text-xl font-bold mb-4">Charts</h2>

        {/* Month-Wise Budget Chart */}
        {/*}
        <div className="mb-8">
            <h3 className="text-lg font-semibold mb-2">Month-wise Budget Chart</h3>
            <DatePicker
            selected={selectedChartMonth}
            onChange={(date) => setSelectedChartMonth(date!)}
            dateFormat="MM/yyyy"
            showMonthYearPicker
            className="mb-4 border rounded p-2"
            />
            <MonthWiseBudgetChart
            budgets={budgets}
            selectedMonth={selectedChartMonth}
            />
        </div>

        {/* Category-Wise Budget Chart */}
        {/*}
        <div>
            <h3 className="text-lg font-semibold mb-2">Category-wise Budget Chart</h3>
            <select
            value={selectedCategory || ""}
            onChange={(e) => setSelectedCategory(Number(e.target.value))}
            className="mb-4 border rounded p-2 w-full"
            >
            <option value="">Select a category</option>
            {categories.map((category) => (
                <option key={category.id} value={category.id}>
                {category.name}
                </option>
            ))}
            </select>
            {selectedCategory && (
            <CategoryWiseBudgetChart
                budgets={budgets}
                categories={categories}
                selectedCategoryId={selectedCategory}
            />
            )}
        </div>
        </div>*/}

    </div> 

  );
}
