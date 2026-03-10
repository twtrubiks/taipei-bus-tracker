import { BrowserRouter, Routes, Route } from "react-router-dom";
import HomePage from "./pages/HomePage";
import SearchPage from "./pages/SearchPage";
import RoutePage from "./pages/RoutePage";
import { useDarkMode } from "./hooks/useDarkMode";

export default function App() {
  const { isDark, toggle } = useDarkMode();

  return (
    <BrowserRouter>
      <div className="min-h-screen bg-white text-gray-900 dark:bg-gray-900 dark:text-gray-100">
        <header className="flex items-center justify-end px-4 py-2">
          <button
            type="button"
            onClick={toggle}
            aria-label={isDark ? "切換淺色模式" : "切換深色模式"}
            className="rounded-lg px-3 py-1 text-sm text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-800"
          >
            {isDark ? "\u2600\uFE0F 淺色" : "\u{1F319} 深色"}
          </button>
        </header>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/search" element={<SearchPage />} />
          <Route path="/route/:routeId" element={<RoutePage />} />
        </Routes>
      </div>
    </BrowserRouter>
  );
}
