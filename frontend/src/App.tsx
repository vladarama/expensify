import { Routes, Route } from "react-router-dom";
import { Navbar } from "./components/navbar";
import { Dashboard } from "./pages/dashboard";
import { Income } from "./pages/income";
import { Expenses } from "./pages/expenses";
import { Categories } from "./pages/categories";
import { Budgets } from "./pages/budgets";

function App() {
  return (
    <div className="min-h-screen bg-background">
      <Navbar />
      <div className="container mx-auto px-4 py-8">
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/income" element={<Income />} />
          <Route path="/expenses" element={<Expenses />} />
          <Route path="/categories" element={<Categories />} />
          <Route path="/budgets" element={<Budgets />} />
        </Routes>
      </div>
    </div>
  );
}

export default App;
