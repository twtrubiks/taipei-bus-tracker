import { useState, useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import { searchRoutes } from "../api/client";
import type { Route } from "../api/types";

export default function RouteSearch() {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState<Route[]>([]);
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const timerRef = useRef<ReturnType<typeof setTimeout>>();

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const val = e.target.value;
    setQuery(val);
    if (!val.trim()) {
      setResults([]);
    }
  };

  useEffect(() => {
    const trimmed = query.trim();
    if (!trimmed) return;

    clearTimeout(timerRef.current);
    timerRef.current = setTimeout(() => {
      setLoading(true);
      searchRoutes(trimmed)
        .then(setResults)
        .catch(() => setResults([]))
        .finally(() => setLoading(false));
    }, 300);

    return () => clearTimeout(timerRef.current);
  }, [query]);

  return (
    <div>
      <input
        type="text"
        placeholder="輸入路線名稱，例如 299"
        value={query}
        onChange={handleChange}
        className="w-full rounded-lg border border-gray-300 px-4 py-2 focus:border-blue-500 focus:outline-none"
        aria-label="搜尋路線"
      />

      {loading && <p className="mt-2 text-sm text-gray-500">搜尋中...</p>}

      {!loading && query.trim() && results.length === 0 && (
        <p className="mt-2 text-sm text-gray-500">找不到符合的路線</p>
      )}

      {results.length > 0 && (
        <ul className="mt-2 divide-y rounded-lg border" role="list">
          {results.map((route) => (
            <li key={route.routeId}>
              <button
                type="button"
                className="w-full px-4 py-3 text-left hover:bg-gray-50"
                onClick={() => navigate(`/route/${route.routeId}`)}
              >
                <span className="font-semibold">{route.routeName}</span>
                <span className="ml-2 text-sm text-gray-500">
                  {route.startStop} &harr; {route.endStop}
                </span>
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
