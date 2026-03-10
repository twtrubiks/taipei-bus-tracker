import { Link } from "react-router-dom";
import { useFavorites } from "../hooks/useFavorites";
import { useFavoritesEta } from "../hooks/useFavoritesEta";
import { useNotificationContext } from "../hooks/NotificationContext";
import { statusColor } from "../utils/statusColor";

export default function HomePage() {
  const { favorites, removeFavorite } = useFavorites();
  const { checkAlerts } = useNotificationContext();
  const favoritesEta = useFavoritesEta(favorites, checkAlerts);

  return (
    <div className="mx-auto max-w-lg p-4 md:max-w-2xl">
      <h1 className="mb-6 text-2xl font-bold">公車到站查詢</h1>
      <Link
        to="/search"
        className="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-3 text-white shadow hover:bg-blue-700"
      >
        <span className="text-xl">&#128269;</span>
        搜尋路線
      </Link>

      {favorites.length > 0 && (
        <section className="mt-6">
          <h2 className="mb-3 text-lg font-semibold">收藏站點</h2>
          <ul className="grid grid-cols-1 gap-3 md:grid-cols-2" role="list">
            {favoritesEta.map(({ favorite: f, eta }) => (
              <li
                key={`${f.routeId}:${f.direction}:${f.stopId}`}
                className="flex items-center gap-3 rounded-lg border border-gray-200 px-4 py-3 dark:border-gray-700"
              >
                <Link
                  to={`/route/${f.routeId}?name=${encodeURIComponent(f.routeName)}&dir=${f.direction}`}
                  className="flex flex-1 items-center gap-2"
                >
                  <span className="rounded bg-blue-100 px-2 py-0.5 text-sm font-semibold text-blue-700 dark:bg-blue-900 dark:text-blue-300">
                    {f.routeName}
                  </span>
                  <span className="text-sm text-gray-500">
                    {f.direction === 0 ? "去" : "回"}
                  </span>
                  <span>{f.stopName}</span>
                </Link>
                <span
                  className={`text-sm ${statusColor(eta?.eta ?? -999)}`}
                >
                  {eta?.status ?? "—"}
                </span>
                {eta?.buses && eta.buses.length > 0 && (
                  <span className="text-xs text-gray-400">
                    {eta.buses.map((b) => b.plateNumb).join(", ")}
                  </span>
                )}
                <button
                  type="button"
                  aria-label="移除收藏"
                  className="text-gray-400 hover:text-red-500"
                  onClick={() =>
                    removeFavorite(f.routeId, f.direction, f.stopId)
                  }
                >
                  &#10005;
                </button>
              </li>
            ))}
          </ul>
        </section>
      )}
    </div>
  );
}
