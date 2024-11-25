import { useState, useEffect, useMemo } from "react";
import { Income as IncomeType } from "../types/income";
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
import { IncomeSourceChart } from "@/components/charts/income-source-chart";
import { SortButton } from "@/components/ui/sort-button";
import { cn } from "@/lib/utils";
import { API_ENDPOINTS } from "@/config/api";

type SortField = "amount" | "date";
type SortDirection = "asc" | "desc" | null;

interface SortState {
  field: SortField | null;
  direction: SortDirection;
}

export function Income() {
  const [incomes, setIncomes] = useState<IncomeType[]>([]);
  const [isOpen, setIsOpen] = useState(false);
  const [newIncome, setNewIncome] = useState({
    amount: "",
    date: new Date().toISOString().split("T")[0],
    source: "",
  });
  const [editingIncome, setEditingIncome] = useState<IncomeType | null>(null);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [sortState, setSortState] = useState<SortState>({
    field: null,
    direction: null,
  });

  useEffect(() => {
    fetch(API_ENDPOINTS.incomes.getAll)
      .then((response) => response.json())
      .then((data) => setIncomes(data));
  }, []);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    try {
      const formattedDate = new Date(newIncome.date).toISOString();

      const response = await fetch(API_ENDPOINTS.incomes.create, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          amount: parseFloat(newIncome.amount),
          date: formattedDate,
          source: newIncome.source,
        }),
      });

      if (response.ok) {
        const data = await response.json();
        setIncomes([...incomes, data]);
        setIsOpen(false);
        setNewIncome({
          amount: "",
          date: new Date().toISOString().split("T")[0],
          source: "",
        });
      }
    } catch (error) {
      console.error("Failed to add income:", error);
    }
  }

  async function handleDelete(id: number) {
    try {
      const response = await fetch(API_ENDPOINTS.incomes.delete(id), {
        method: "DELETE",
      });

      if (response.ok) {
        setIncomes(incomes.filter((income) => income.id !== id));
      }
    } catch (error) {
      console.error("Failed to delete income:", error);
    }
  }

  async function handleEdit(e: React.FormEvent) {
    e.preventDefault();
    if (!editingIncome) return;

    try {
      const formattedDate = new Date(editingIncome.date).toISOString();

      const response = await fetch(
        API_ENDPOINTS.incomes.update(editingIncome.id),
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            amount: editingIncome.amount,
            date: formattedDate,
            source: editingIncome.source,
          }),
        }
      );

      if (response.ok) {
        const updatedIncome = await response.json();
        setIncomes(
          incomes.map((income) =>
            income.id === editingIncome.id ? updatedIncome : income
          )
        );
        setIsEditDialogOpen(false);
        setEditingIncome(null);
      }
    } catch (error) {
      console.error("Failed to update income:", error);
    }
  }

  function handleStartEdit(income: IncomeType) {
    setEditingIncome({
      ...income,
      date: new Date(income.date).toISOString().split("T")[0],
    });
    setIsEditDialogOpen(true);
  }

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

  const sortedIncomes = useMemo(() => {
    if (!sortState.field || !sortState.direction) {
      return incomes;
    }

    return [...incomes].sort((a, b) => {
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

      return 0;
    });
  }, [sortState.field, sortState.direction, incomes]);

  return (
    <div className="min-h-screen p-4 sm:p-6">
      <div className="max-w-7xl mx-auto space-y-6 sm:space-y-8">
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 sm:gap-0">
          <h1 className="text-2xl sm:text-3xl font-bold bg-gradient-to-r from-emerald-600 to-green-500 bg-clip-text text-transparent">
            Income Dashboard
          </h1>
          <Dialog open={isOpen} onOpenChange={setIsOpen}>
            <DialogTrigger asChild>
              <Button
                className={cn(
                  "w-full sm:w-auto",
                  "bg-gradient-to-r from-emerald-600 to-green-500",
                  "hover:from-emerald-700 hover:to-green-600",
                  "text-white shadow-md hover:shadow-lg",
                  "transition-all duration-200"
                )}
              >
                + Add Income
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Add New Income</DialogTitle>
              </DialogHeader>
              <form onSubmit={handleSubmit} className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="amount">Amount</Label>
                  <Input
                    id="amount"
                    type="number"
                    step="0.01"
                    required
                    value={newIncome.amount}
                    onChange={(e) =>
                      setNewIncome({ ...newIncome, amount: e.target.value })
                    }
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="date">Date</Label>
                  <Input
                    id="date"
                    type="date"
                    required
                    value={newIncome.date}
                    onChange={(e) =>
                      setNewIncome({ ...newIncome, date: e.target.value })
                    }
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="source">Source</Label>
                  <Input
                    id="source"
                    required
                    value={newIncome.source}
                    onChange={(e) =>
                      setNewIncome({ ...newIncome, source: e.target.value })
                    }
                  />
                </div>
                <Button
                  type="submit"
                  className={cn(
                    "w-full",
                    "bg-gradient-to-r from-emerald-600 to-green-500",
                    "hover:from-emerald-700 hover:to-green-600",
                    "text-white shadow-md hover:shadow-lg",
                    "transition-all duration-200"
                  )}
                >
                  Add Income
                </Button>
              </form>
            </DialogContent>
          </Dialog>
        </div>

        <div className="grid grid-cols-1 xl:grid-cols-2 gap-4 sm:gap-8">
          <div className="bg-white rounded-xl shadow-lg p-4 sm:p-6 overflow-hidden flex flex-col min-h-[500px] xl:h-[600px]">
            <h2 className="text-lg sm:text-xl font-semibold mb-4 text-gray-800">
              Income History
            </h2>
            <div className="overflow-x-auto flex-1">
              <div className="overflow-y-auto h-full">
                <div className="space-y-4">
                  {sortedIncomes.map((income) => (
                    <div
                      key={income.id}
                      className="group bg-white p-4 rounded-lg shadow-sm hover:bg-emerald-50/50 transition-colors lg:hidden"
                    >
                      <div className="flex justify-between items-start mb-2">
                        <div>
                          <div className="font-medium text-gray-900">
                            {income.source}
                          </div>
                          <div className="text-emerald-600 font-semibold mt-1">
                            ${income.amount.toLocaleString()}
                          </div>
                        </div>
                        <div className="flex gap-2">
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => handleStartEdit(income)}
                            className="h-8 w-8 hover:bg-emerald-100 rounded-full"
                          >
                            <Pencil className="h-4 w-4 text-emerald-600" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => handleDelete(income.id)}
                            className="h-8 w-8 hover:bg-red-100 rounded-full"
                          >
                            <Trash2 className="h-4 w-4 text-red-500" />
                          </Button>
                        </div>
                      </div>
                      <div className="text-gray-600 text-sm">
                        {new Date(income.date).toLocaleDateString()}
                      </div>
                    </div>
                  ))}

                  {/* Traditional table for larger screens */}
                  <div className="hidden lg:block">
                    <table className="w-full">
                      <thead>
                        <tr className="border-b border-gray-200">
                          <th className="text-left py-3 px-4 font-semibold text-gray-600 w-1/4">
                            Source
                          </th>
                          <th className="text-left py-3 px-4 font-semibold text-gray-600 w-1/4">
                            <SortButton
                              label="Amount"
                              active={sortState.field === "amount"}
                              direction={
                                sortState.field === "amount"
                                  ? sortState.direction
                                  : null
                              }
                              onClick={() => handleSort("amount")}
                              className="hover:text-emerald-600 hover:bg-transparent"
                            />
                          </th>
                          <th className="text-left py-3 px-4 font-semibold text-gray-600 w-2/4 min-w-[200px]">
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
                                className="hover:text-emerald-600 hover:bg-transparent"
                              />
                            </div>
                          </th>
                        </tr>
                      </thead>
                      <tbody>
                        {sortedIncomes.map((income) => (
                          <tr
                            key={income.id}
                            className="group border-b border-gray-100 last:border-none hover:bg-emerald-50/50 transition-colors"
                          >
                            <td className="py-3 px-4">
                              <span className="font-medium text-gray-900">
                                {income.source}
                              </span>
                            </td>
                            <td className="py-3 px-4">
                              <span className="text-emerald-600 font-semibold">
                                ${income.amount.toLocaleString()}
                              </span>
                            </td>
                            <td className="py-3 px-4">
                              <div className="flex justify-between items-center gap-2">
                                <span className="text-gray-600 text-sm sm:text-base">
                                  {new Date(income.date).toLocaleDateString()}
                                </span>
                                <div className="flex gap-1 sm:gap-2 transition-opacity">
                                  <Button
                                    variant="ghost"
                                    size="icon"
                                    onClick={() => handleStartEdit(income)}
                                    className="h-7 w-7 sm:h-8 sm:w-8 hover:bg-emerald-100 rounded-full"
                                  >
                                    <Pencil className="h-3 w-3 sm:h-4 sm:w-4 text-emerald-600" />
                                  </Button>
                                  <Button
                                    variant="ghost"
                                    size="icon"
                                    onClick={() => handleDelete(income.id)}
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
          </div>

          <div className="bg-white rounded-xl shadow-lg p-4 sm:p-6 flex flex-col min-h-[500px] xl:h-[600px]">
            <h2 className="text-lg sm:text-xl font-semibold mb-4 text-gray-800">
              Income by Source
            </h2>
            <div className="flex-1">
              <IncomeSourceChart incomes={incomes} />
            </div>
          </div>
        </div>

        <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Edit Income</DialogTitle>
            </DialogHeader>
            <form onSubmit={handleEdit} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="edit-amount">Amount</Label>
                <Input
                  id="edit-amount"
                  type="number"
                  step="0.01"
                  required
                  value={editingIncome?.amount || ""}
                  onChange={(e) =>
                    setEditingIncome((prev) =>
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
                  value={editingIncome?.date || ""}
                  onChange={(e) =>
                    setEditingIncome((prev) =>
                      prev ? { ...prev, date: e.target.value } : null
                    )
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="edit-source">Source</Label>
                <Input
                  id="edit-source"
                  required
                  value={editingIncome?.source || ""}
                  onChange={(e) =>
                    setEditingIncome((prev) =>
                      prev ? { ...prev, source: e.target.value } : null
                    )
                  }
                />
              </div>
              <Button
                type="submit"
                className={cn(
                  "w-full",
                  "bg-gradient-to-r from-emerald-600 to-green-500",
                  "hover:from-emerald-700 hover:to-green-600",
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
