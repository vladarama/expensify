import { useState, useEffect } from "react";
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

  useEffect(() => {
    fetch("http://localhost:8080/incomes")
      .then((response) => response.json())
      .then((data) => setIncomes(data));
  }, []);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    try {
      const formattedDate = new Date(newIncome.date).toISOString();

      const response = await fetch("http://localhost:8080/incomes", {
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
      const response = await fetch(`http://localhost:8080/incomes/${id}`, {
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
        `http://localhost:8080/incomes/${editingIncome.id}`,
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

  return (
    <div>
      <h1 className="text-2xl font-bold mb-4">Income</h1>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <div className="overflow-auto">
          <Dialog open={isOpen} onOpenChange={setIsOpen}>
            <DialogTrigger asChild>
              <Button className="mb-4">+ Add Income</Button>
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
                <Button type="submit" className="w-full">
                  Add Income
                </Button>
              </form>
            </DialogContent>
          </Dialog>

          <table className="min-w-full bg-white">
            <thead>
              <tr>
                <th className="border px-4 py-2">Amount</th>
                <th className="border px-4 py-2">Date</th>
                <th className="border px-4 py-2">Source</th>
              </tr>
            </thead>
            <tbody>
              {incomes.map((income) => (
                <tr key={income.id} className="group hover:bg-gray-50">
                  <td className="border px-4 py-2">${income.amount}</td>
                  <td className="border px-4 py-2">
                    {new Date(income.date).toLocaleDateString()}
                  </td>
                  <td className="border px-4 py-2 relative">
                    {income.source}
                    <div className="absolute right-2 top-1/2 -translate-y-1/2 opacity-0 group-hover:opacity-100 flex gap-2">
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => handleStartEdit(income)}
                        className="h-8 w-8 hover:bg-gray-100/50"
                      >
                        <Pencil className="h-4 w-4 text-blue-500" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => handleDelete(income.id)}
                        className="h-8 w-8 hover:bg-gray-100/50"
                      >
                        <Trash2 className="h-4 w-4 text-red-500" />
                      </Button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        <div>
          <IncomeSourceChart incomes={incomes} />
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
            <Button type="submit" className="w-full">
              Save Changes
            </Button>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  );
}
