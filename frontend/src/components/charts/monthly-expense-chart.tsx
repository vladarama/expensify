import { useMemo } from "react";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import { Expense } from "@/types/expense";

interface MonthlyExpenseChartProps {
  expenses: Expense[];
}

export function MonthlyExpenseChart({ expenses }: MonthlyExpenseChartProps) {
  const monthlyData = useMemo(() => {
    const lastYear = new Date();
    lastYear.setFullYear(lastYear.getFullYear() - 1);

    // Create an array of the last 12 months
    const months: { [key: string]: number } = {};
    for (let i = 0; i < 12; i++) {
      const date = new Date();
      date.setMonth(date.getMonth() - i);
      const monthKey = date.toLocaleString("default", {
        month: "short",
        year: "2-digit",
      });
      months[monthKey] = 0;
    }

    // Sum expenses for each month
    expenses.forEach((expense) => {
      const expenseDate = new Date(expense.date);
      if (expenseDate >= lastYear) {
        const monthKey = expenseDate.toLocaleString("default", {
          month: "short",
          year: "2-digit",
        });
        if (monthKey in months) {
          months[monthKey] += expense.amount;
        }
      }
    });

    // Convert to array and reverse to show oldest to newest
    return Object.entries(months)
      .map(([month, amount]) => ({
        month,
        amount: Number(amount.toFixed(2)),
      }))
      .reverse();
  }, [expenses]);

  return (
    <div className="w-full h-full min-h-[500px]">
      <ResponsiveContainer width="100%" height="90%">
        <BarChart data={monthlyData}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="month" />
          <YAxis />
          <Tooltip
            formatter={(value: number) => [`$${value.toFixed(2)}`, "Amount"]}
          />
          <Bar
            dataKey="amount"
            fill="url(#redGradient)"
            radius={[4, 4, 0, 0]}
          />
          <defs>
            <linearGradient id="redGradient" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stopColor="#dc2626" stopOpacity={1} />
              <stop offset="100%" stopColor="#f43f5e" stopOpacity={1} />
            </linearGradient>
          </defs>
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}
