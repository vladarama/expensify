import { useState, useEffect, useCallback } from "react";
import { Category } from "../types/category";
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
import { CategoryExpenseChart } from "@/components/charts/category-expense-chart";
import { Expense } from "@/types/expense";

export function CategoriesPage() {
  return (
    <ErrorBoundary>
      <Categories />
    </ErrorBoundary>
  );
}

export function Categories() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isOpen, setIsOpen] = useState(false);
  const [newCategory, setNewCategory] = useState({
    name: "",
    description: "",
  });
  const [editingCategory, setEditingCategory] = useState<Category | null>(null);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [expenses, setExpenses] = useState<Expense[]>([]);

  const getCategoryName = useCallback(
    (categoryId: number) => {
      const category = categories.find((cat) => cat.id === categoryId);
      return category?.name || "Unknown Category";
    },
    [categories]
  );

  useEffect(() => {
    let isMounted = true;

    async function fetchCategories() {
      try {
        setIsLoading(true);
        const response = await fetch("http://localhost:8080/categories");
        if (!response.ok)
          throw new Error(`HTTP error! status: ${response.status}`);
        const data = await response.json();
        if (isMounted) {
          setCategories(data || []);
        }
      } catch (error) {
        console.error("Failed to fetch categories:", error);
        if (isMounted) {
          setCategories([]);
        }
      } finally {
        if (isMounted) {
          setIsLoading(false);
        }
      }
    }

    fetchCategories();

    return () => {
      isMounted = false;
    };
  }, []);

  useEffect(() => {
    let isMounted = true;

    async function fetchExpenses() {
      try {
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
      }
    }

    fetchExpenses();
    return () => {
      isMounted = false;
    };
  }, []);

  if (isLoading) {
    return <div className="p-4">Loading categories...</div>;
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    try {
      const response = await fetch("http://localhost:8080/categories", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(newCategory),
      });

      if (response.ok) {
        const data = await response.json();
        setCategories([...categories, data]);
        setIsOpen(false);
        setNewCategory({ name: "", description: "" });
      }
    } catch (error) {
      console.error("Failed to add category:", error);
    }
  }

  async function handleDelete(id: number) {
    try {
      const response = await fetch(`http://localhost:8080/categories/${id}`, {
        method: "DELETE",
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || "Failed to delete category");
      }

      setCategories(categories.filter((category) => category.id !== id));
    } catch (error) {
      console.error("Failed to delete category:", error);
      alert("Failed to delete category. Please try again.");
    }
  }

  async function handleEdit(e: React.FormEvent) {
    e.preventDefault();
    if (!editingCategory) return;

    try {
      const response = await fetch(
        `http://localhost:8080/categories/${editingCategory.id}`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            name: editingCategory.name,
            description: editingCategory.description,
          }),
        }
      );

      if (response.ok) {
        const updatedCategory = await response.json();
        setCategories(
          categories.map((category) =>
            category.id === editingCategory.id ? updatedCategory : category
          )
        );
        setIsEditDialogOpen(false);
        setEditingCategory(null);
      } else {
        console.error("Failed to update category:", await response.text());
      }
    } catch (error) {
      console.error("Failed to update category:", error);
    }
  }

  function handleStartEdit(category: Category) {
    setEditingCategory(category);
    setIsEditDialogOpen(true);
  }

  return (
    <div>
      <h1 className="text-2xl font-bold mb-4">Categories</h1>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <div className="overflow-auto">
          <Dialog open={isOpen} onOpenChange={setIsOpen}>
            <DialogTrigger asChild>
              <Button className="mb-4">+ Add Category</Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Add New Category</DialogTitle>
              </DialogHeader>
              <form onSubmit={handleSubmit} className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="name">Name</Label>
                  <Input
                    id="name"
                    required
                    value={newCategory.name}
                    onChange={(e) =>
                      setNewCategory({ ...newCategory, name: e.target.value })
                    }
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="description">Description</Label>
                  <Input
                    id="description"
                    required
                    value={newCategory.description}
                    onChange={(e) =>
                      setNewCategory({
                        ...newCategory,
                        description: e.target.value,
                      })
                    }
                  />
                </div>
                <Button type="submit" className="w-full">
                  Add Category
                </Button>
              </form>
            </DialogContent>
          </Dialog>

          <table className="min-w-full bg-white">
            <thead>
              <tr>
                <th className="border px-4 py-2">Name</th>
                <th className="border px-4 py-2">Description</th>
              </tr>
            </thead>
            <tbody>
              {categories?.length > 0 ? (
                categories.map((category) => (
                  <tr key={category.id} className="group hover:bg-gray-50">
                    <td className="border px-4 py-2">{category.name}</td>
                    <td className="border px-4 py-2 relative">
                      {category.description}
                      <div className="absolute right-2 top-1/2 -translate-y-1/2 opacity-0 group-hover:opacity-100 flex gap-2">
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => handleStartEdit(category)}
                          className="h-8 w-8 hover:bg-gray-100/50"
                        >
                          <Pencil className="h-4 w-4 text-blue-500" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => handleDelete(category.id)}
                          className="h-8 w-8 hover:bg-gray-100/50"
                        >
                          <Trash2 className="h-4 w-4 text-red-500" />
                        </Button>
                      </div>
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td
                    colSpan={2}
                    className="border px-4 py-2 text-center text-gray-500"
                  >
                    No categories found
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>

        <div>
          <CategoryExpenseChart
            expenses={expenses}
            categories={categories}
            getCategoryName={getCategoryName}
          />
        </div>
      </div>

      <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Edit Category</DialogTitle>
          </DialogHeader>
          <form onSubmit={handleEdit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="edit-name">Name</Label>
              <Input
                id="edit-name"
                required
                value={editingCategory?.name || ""}
                onChange={(e) =>
                  setEditingCategory((prev) =>
                    prev ? { ...prev, name: e.target.value } : null
                  )
                }
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="edit-description">Description</Label>
              <Input
                id="edit-description"
                required
                value={editingCategory?.description || ""}
                onChange={(e) =>
                  setEditingCategory((prev) =>
                    prev ? { ...prev, description: e.target.value } : null
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
