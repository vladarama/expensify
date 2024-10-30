import { useMemo } from "react";
import {
  PieChart,
  Pie,
  Cell,
  ResponsiveContainer,
  Tooltip,
  Legend,
} from "recharts";
import { Income } from "@/types/income";

interface IncomeSourceChartProps {
  incomes: Income[];
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

export function IncomeSourceChart({ incomes }: IncomeSourceChartProps) {
  const sourceData = useMemo(() => {
    const totals: { [key: string]: number } = {};
    let totalIncome = 0;

    // Calculate totals for each source
    incomes.forEach((income) => {
      totals[income.source] = (totals[income.source] || 0) + income.amount;
      totalIncome += income.amount;
    });

    // Convert to array and calculate percentages
    return Object.entries(totals).map(([source, amount]) => ({
      name: source,
      value: amount,
      percentage: ((amount / totalIncome) * 100).toFixed(1),
    }));
  }, [incomes]);

  return (
    <div className="w-full h-full min-h-[400px]">
      <h2 className="text-xl font-semibold mb-4">Income by Source</h2>
      <ResponsiveContainer width="100%" height="90%">
        <PieChart>
          <Pie
            data={sourceData}
            dataKey="value"
            nameKey="name"
            cx="50%"
            cy="50%"
            outerRadius={120}
            label={({ name, percentage }) => `${name} (${percentage}%)`}
          >
            {sourceData.map((_entry, index) => (
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
