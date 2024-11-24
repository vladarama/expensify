import React from "react";
import { Budget } from "@/types/budget";
import { Bar } from "react-chartjs-2";

// Register Chart.js components
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
} from "chart.js";

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend);

interface Props {
  budgets: Budget[];
  selectedCategoryId: number;
  found_categories: { id: number; name: string }[];
}

export const CategoryWiseBudgetChart: React.FC<Props> = ({
  budgets,
  selectedCategoryId,
//   found_categories,
}) => {
  // Filter budgets by selected category
  const filteredBudgets = budgets.filter(
    (budget) => budget.category_id === selectedCategoryId
  );

  // Format the labels as MM/YY
  const data = {
    labels: filteredBudgets.map((budget) => {
    //   const categoryName =
    //     found_categories.find((category) => category.id === budget.category_id)?.name || "Unknown";
      const budgetDate = new Date(budget.start_date);
      const formattedDate = `${(budgetDate.getMonth() + 1)
        .toString()
        .padStart(2, "0")}/${budgetDate.getFullYear().toString().slice(-2)}`;
      return `${formattedDate}`;
    }),
    datasets: [
      {
        label: "Amount",
        data: filteredBudgets.map((budget) => budget.amount),
        backgroundColor: "rgba(75, 192, 192, 0.6)",
      },
      {
        label: "Spent",
        data: filteredBudgets.map((budget) => budget.spent),
        backgroundColor: "rgba(255, 99, 132, 0.6)",
      },
    ],
  };

  return (
    <div>
      {/* <h2 className="text-xl font-bold mb-4">Category-wise Budget Chart</h2> */}
      <Bar data={data} options={{ responsive: true }} />
    </div>
  );
};
