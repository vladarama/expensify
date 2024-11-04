import { useMemo } from "react";
import {
  PieChart,
  Pie,
  Cell,
  ResponsiveContainer,
  Tooltip,
  Legend,
} from "recharts";
import { Expense } from "@/types/expense";
import { Category } from "@/types/category";

interface CategoryExpenseChartProps {
  expenses: Expense[];
  categories: Category[];
  getCategoryName: (id: number) => string;
}

const COLORS = [
  "#8884d8",
  "#82ca9d",
  "#ffc658",
  "#ff8042",
  "#0088FE",
  "#00C49F",
  "#FFBB28",
  "#FF8042",
  "#a4de6c",
  "#d0ed57",
];

export function CategoryExpenseChart({
  expenses,
  getCategoryName,
}: CategoryExpenseChartProps) {
  const categoryData = useMemo(() => {
    const totals: { [key: number]: number } = {};
    let totalExpenses = 0;

    // Calculate totals for each category
    expenses.forEach((expense) => {
      totals[expense.category_id] =
        (totals[expense.category_id] || 0) + expense.amount;
      totalExpenses += expense.amount;
    });

    // Convert to array and calculate percentages
    return Object.entries(totals).map(([categoryId, amount]) => ({
      name: getCategoryName(Number(categoryId)),
      value: amount,
      percentage: ((amount / totalExpenses) * 100).toFixed(1),
    }));
  }, [expenses, getCategoryName]);

  return (
    <div className="w-full h-full min-h-[400px]">
      <h2 className="text-xl font-semibold mb-4">Expenses by Category</h2>
      <ResponsiveContainer width="100%" height="90%">
        <PieChart>
          <Pie
            data={categoryData}
            dataKey="value"
            nameKey="name"
            cx="50%"
            cy="50%"
            outerRadius={120}
            label={({ name, percentage }) => `${name} (${percentage}%)`}
          >
            {categoryData.map((_entry, index) => (
              <Cell
                key={`cell-${index}`}
                fill={COLORS[index % COLORS.length]}
              />
            ))}
          </Pie>
          <Tooltip formatter={(value: number) => `$${value.toFixed(2)}`} />
          <Legend />
        </PieChart>
      </ResponsiveContainer>
    </div>
  );
}
