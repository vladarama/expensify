import { Link, useLocation } from "react-router-dom";
import { cn } from "@/lib/utils";

export function Navbar() {
  const location = useLocation();
  const pathname = location.pathname;

  return (
    <nav className="border-b">
      <div className="max-w-7xl mx-auto px-4 sm:px-6">
        <div className="flex justify-between h-16">
          <div className="flex items-center space-x-2 mr-8">
            <img
              src="/logo.png"
              alt="Expensify Logo"
              className="w-8 h-8 rounded-sm"
            />
            <span className="text-xl font-bold">Expensify</span>
          </div>

          <div className="flex space-x-8">
            <Link
              to="/income"
              className={cn(
                "inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium",
                pathname === "/income"
                  ? "border-emerald-500 text-emerald-600"
                  : "border-transparent text-gray-500 hover:border-emerald-300 hover:text-emerald-600",
                "transition-colors duration-200"
              )}
            >
              Income
            </Link>

            <Link
              to="/expenses"
              className={cn(
                "inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium",
                pathname === "/expenses"
                  ? "border-red-500 text-red-600"
                  : "border-transparent text-gray-500 hover:border-red-300 hover:text-red-600",
                "transition-colors duration-200"
              )}
            >
              Expenses
            </Link>

            <Link
              to="/categories"
              className={cn(
                "inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium",
                pathname === "/categories"
                  ? "border-orange-500 text-orange-600"
                  : "border-transparent text-gray-500 hover:border-orange-300 hover:text-orange-600",
                "transition-colors duration-200"
              )}
            >
              Categories
            </Link>

            <Link
              to="/budgets"
              className={cn(
                "inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium",
                pathname === "/budgets"
                  ? "border-purple-500 text-purple-600"
                  : "border-transparent text-gray-500 hover:border-purple-300 hover:text-purple-600",
                "transition-colors duration-200"
              )}
            >
              Budgets
            </Link>
          </div>
        </div>
      </div>
    </nav>
  );
}
