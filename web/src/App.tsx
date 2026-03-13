import { BrowserRouter, Routes, Route, Link, useLocation } from "react-router-dom";
import HomePage from "./pages/HomePage";
import SearchPage from "./pages/SearchPage";
import RoutePage from "./pages/RoutePage";
import { useDarkMode } from "./hooks/useDarkMode";
import { NotificationProvider } from "./hooks/NotificationContext";

function BottomNav() {
  const { pathname } = useLocation();
  const isHome = pathname === "/";
  const isSearch = pathname === "/search" || pathname.startsWith("/route/");

  return (
    <nav className="fixed bottom-0 left-0 right-0 z-50 border-t border-gray-200 bg-white pb-[env(safe-area-inset-bottom)] dark:border-gray-700 dark:bg-gray-900">
      <div className="mx-auto flex max-w-lg items-center justify-around md:max-w-2xl">
        <Link
          to="/"
          className={`flex flex-1 flex-col items-center py-2 text-xs ${isHome ? "text-blue-600 dark:text-blue-400" : "text-gray-500 dark:text-gray-400"}`}
        >
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" className="mb-0.5 h-6 w-6">
            <path d="M11.47 3.841a.75.75 0 0 1 1.06 0l8.69 8.69a.75.75 0 1 0 1.06-1.061l-8.689-8.69a2.25 2.25 0 0 0-3.182 0l-8.69 8.69a.75.75 0 1 0 1.061 1.06l8.69-8.689Z" />
            <path d="m12 5.432 8.159 8.159c.03.03.06.058.091.086v6.198c0 1.035-.84 1.875-1.875 1.875H15a.75.75 0 0 1-.75-.75v-4.5a.75.75 0 0 0-.75-.75h-3a.75.75 0 0 0-.75.75V21a.75.75 0 0 1-.75.75H5.625a1.875 1.875 0 0 1-1.875-1.875v-6.198a.75.75 0 0 1 .091-.086L12 5.432Z" />
          </svg>
          首頁
        </Link>
        <Link
          to="/search"
          className={`flex flex-1 flex-col items-center py-2 text-xs ${isSearch ? "text-blue-600 dark:text-blue-400" : "text-gray-500 dark:text-gray-400"}`}
        >
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" className="mb-0.5 h-6 w-6">
            <path fillRule="evenodd" d="M10.5 3.75a6.75 6.75 0 1 0 0 13.5 6.75 6.75 0 0 0 0-13.5ZM2.25 10.5a8.25 8.25 0 1 1 14.59 5.28l4.69 4.69a.75.75 0 1 1-1.06 1.06l-4.69-4.69A8.25 8.25 0 0 1 2.25 10.5Z" clipRule="evenodd" />
          </svg>
          搜尋
        </Link>
      </div>
    </nav>
  );
}

export default function App() {
  const { isDark, toggle } = useDarkMode();

  return (
    <BrowserRouter>
      <NotificationProvider>
        <div className="min-h-screen bg-white pb-16 text-gray-900 dark:bg-gray-900 dark:text-gray-100">
          <header className="mx-auto flex max-w-lg items-center justify-between px-4 py-2 md:max-w-2xl">
            <Link to="/" className="text-sm font-medium text-gray-400 hover:text-gray-600 dark:text-gray-500 dark:hover:text-gray-300">
              Taipei Bus Tracker
            </Link>
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
          <BottomNav />
        </div>
      </NotificationProvider>
    </BrowserRouter>
  );
}
