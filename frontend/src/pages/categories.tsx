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
import { cn } from "@/lib/utils";
import { API_ENDPOINTS } from "@/config/api";

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
        const response = await fetch(API_ENDPOINTS.categories.getAll);
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
      const response = await fetch(API_ENDPOINTS.categories.create, {
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
      const response = await fetch(API_ENDPOINTS.categories.delete(id), {
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
        API_ENDPOINTS.categories.update(editingCategory.id),
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
    <div className="min-h-screen p-4 sm:p-6">
      <div className="max-w-7xl mx-auto space-y-6 sm:space-y-8">
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 sm:gap-0">
          <h1 className="text-2xl sm:text-3xl font-bold bg-gradient-to-r from-orange-600 to-amber-500 bg-clip-text text-transparent">
            Categories Dashboard
          </h1>
          <Dialog open={isOpen} onOpenChange={setIsOpen}>
            <DialogTrigger asChild>
              <Button
                className={cn(
                  "w-full sm:w-auto",
                  "bg-gradient-to-r from-orange-600 to-amber-500",
                  "hover:from-orange-700 hover:to-amber-600",
                  "text-white shadow-md hover:shadow-lg",
                  "transition-all duration-200"
                )}
              >
                + Add Category
              </Button>
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
                <Button
                  type="submit"
                  className={cn(
                    "w-full",
                    "bg-gradient-to-r from-orange-600 to-amber-500",
                    "hover:from-orange-700 hover:to-amber-600",
                    "text-white shadow-md hover:shadow-lg",
                    "transition-all duration-200"
                  )}
                >
                  Add Category
                </Button>
              </form>
            </DialogContent>
          </Dialog>
        </div>

        <div className="grid grid-cols-1 xl:grid-cols-2 gap-4 sm:gap-8">
          <div className="bg-white rounded-xl shadow-lg p-4 sm:p-6 overflow-hidden flex flex-col min-h-[500px] xl:h-[600px]">
            <h2 className="text-lg sm:text-xl font-semibold mb-4 text-gray-800">
              Category List
            </h2>
            <div className="space-y-4">
              {categories?.length > 0 ? (
                categories.map((category) => (
                  <div
                    key={category.id}
                    className="group bg-white p-4 rounded-lg shadow-sm hover:bg-orange-50/50 transition-colors lg:hidden"
                  >
                    <div className="flex justify-between items-start mb-2">
                      <div>
                        <div className="font-medium text-gray-900">
                          {category.name}
                        </div>
                        <div className="text-gray-600 mt-1">
                          {category.description}
                        </div>
                      </div>
                      <div className="flex gap-2">
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => handleStartEdit(category)}
                          className="h-8 w-8 hover:bg-orange-100 rounded-full"
                        >
                          <Pencil className="h-4 w-4 text-orange-600" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => handleDelete(category.id)}
                          className="h-8 w-8 hover:bg-red-100 rounded-full"
                        >
                          <Trash2 className="h-4 w-4 text-red-500" />
                        </Button>
                      </div>
                    </div>
                  </div>
                ))
              ) : (
                <div className="text-center text-gray-500 py-4">
                  No categories found
                </div>
              )}

              {/* Traditional table for larger screens */}
              <div className="hidden lg:block">
                <table className="w-full">
                  <thead>
                    <tr className="border-b border-gray-200">
                      <th className="text-left py-3 px-4 font-semibold text-gray-600 w-1/4">
                        Name
                      </th>
                      <th className="text-left py-3 px-4 font-semibold text-gray-600 w-2/4">
                        Description
                      </th>
                    </tr>
                  </thead>
                  <tbody>
                    {categories?.length > 0 ? (
                      categories.map((category) => (
                        <tr
                          key={category.id}
                          className="group border-b border-gray-100 last:border-none hover:bg-orange-50/50 transition-colors"
                        >
                          <td className="py-3 px-4">
                            <span className="font-medium text-gray-900">
                              {category.name}
                            </span>
                          </td>
                          <td className="py-3 px-4">
                            <span className="text-gray-600">
                              {category.description}
                            </span>
                          </td>
                          <td className="py-3 px-4">
                            <div className="flex gap-1 sm:gap-2 transition-opacity">
                              <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => handleStartEdit(category)}
                                className="h-7 w-7 sm:h-8 sm:w-8 hover:bg-orange-100 rounded-full"
                              >
                                <Pencil className="h-3 w-3 sm:h-4 sm:w-4 text-orange-600" />
                              </Button>
                              <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => handleDelete(category.id)}
                                className="h-7 w-7 sm:h-8 sm:w-8 hover:bg-red-100 rounded-full"
                              >
                                <Trash2 className="h-3 w-3 sm:h-4 sm:w-4 text-red-500" />
                              </Button>
                            </div>
                          </td>
                        </tr>
                      ))
                    ) : (
                      <tr>
                        <td
                          colSpan={3}
                          className="py-3 px-4 text-center text-gray-500"
                        >
                          No categories found
                        </td>
                      </tr>
                    )}
                  </tbody>
                </table>
              </div>
            </div>
          </div>

          <div className="bg-white rounded-xl shadow-lg p-4 sm:p-6 flex flex-col min-h-[500px] xl:h-[600px]">
            <div className="flex-1">
              <CategoryExpenseChart
                expenses={expenses}
                categories={categories}
                getCategoryName={getCategoryName}
              />
            </div>
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
              <Button
                type="submit"
                className={cn(
                  "w-full",
                  "bg-gradient-to-r from-orange-600 to-amber-500",
                  "hover:from-orange-700 hover:to-amber-600",
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
