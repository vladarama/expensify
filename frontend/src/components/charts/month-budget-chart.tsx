import React from "react"; //, { useState } from "react";
import { Budget } from "@/types/budget";
// import { Category } from "@/types/category";
// import DatePicker from "react-datepicker";
import "react-datepicker/dist/react-datepicker.css";
import { Bar } from "react-chartjs-2";

// Add this to both month-budget-chart.tsx and category-budget-chart.tsx
import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    BarElement,
    Title,
    Tooltip,
    Legend,
  } from 'chart.js';
  
ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend);
  

interface Props {
  budgets: Budget[];
  selectedMonth: Date;
  found_categories: { id: number; name: string }[];
}

export const MonthWiseBudgetChart: React.FC<Props> = ({ budgets, selectedMonth, found_categories }) => {
    const filteredBudgets = budgets.filter((budget) => {
      const budgetDate = new Date(budget.start_date);
      return (
        budgetDate.getFullYear() === selectedMonth.getFullYear() &&
        budgetDate.getMonth() === selectedMonth.getMonth()
      );
    });
  
    const data = {
      labels: filteredBudgets.map(
        (budget) =>
          found_categories.find((category) => category.id === budget.category_id)?.name || "Unknown"
      ),
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
        {/* <h2 className="text-xl font-bold mb-4">Month-wise Budget Chart</h2> */}
        <Bar data={data} options={{ responsive: true }} />
      </div>
    );
  };
  
  

// export const MonthWiseBudgetChart: React.FC<Props> = ({ budgets }) => {
//   const [selectedMonth, setSelectedMonth] = useState(new Date());

//   const filteredBudgets = budgets.filter((budget) => {
//     const budgetDate = new Date(budget.start_date);
//     return (
//       budgetDate.getFullYear() === selectedMonth.getFullYear() &&
//       budgetDate.getMonth() === selectedMonth.getMonth()
//     );
//   });

//   const data = {
//     labels: filteredBudgets.map(
//       (budget) => `Category: ${budget.category_id}` // Replace with category names if needed
//     ),
//     datasets: [
//       {
//         label: "Amount",
//         data: filteredBudgets.map((budget) => budget.amount),
//         backgroundColor: "rgba(75, 192, 192, 0.6)",
//       },
//       {
//         label: "Spent",
//         data: filteredBudgets.map((budget) => budget.spent),
//         backgroundColor: "rgba(255, 99, 132, 0.6)",
//       },
//     ],
//   };

//   return (
//     <div>
//       <h2 className="text-xl font-bold mb-4">Month-wise Budget Chart</h2>
//       <DatePicker
//         selected={selectedMonth}
//         onChange={(date) => setSelectedMonth(date!)}
//         dateFormat="MM/yyyy"
//         showMonthYearPicker
//         className="mb-4 border rounded p-2"
//       />
//       <Bar data={data} options={{ responsive: true }} />
//     </div>
//   );
// };
