import { useCallback, useEffect, useRef, useState } from "react";
import { Link } from "react-router-dom";
import { useFavorites } from "../hooks/useFavorites";
import { useFavoritesEta } from "../hooks/useFavoritesEta";
import { useNotificationContext } from "../hooks/NotificationContext";
import { AlertBell, AlertMenu } from "../components/AlertButton";
import { searchRoutes, getStops } from "../api/client";
import { statusColor } from "../utils/statusColor";
import { normalizeName } from "../utils/normalize";
import type { StopETA } from "../api/types";

export default function HomePage() {
  const { favorites, removeFavorite, updateFavoriteIds, resolveFavorite } = useFavorites();
  const { checkAlerts, getAlert, addAlert, removeAlert } = useNotificationContext();
  const [alertMenuKey, setAlertMenuKey] = useState<string | null>(null);

  const favoritesRef = useRef(favorites);
  useEffect(() => {
    favoritesRef.current = favorites;
  });

  const handleEtaFetched = useCallback(
    (routeId: string, direction: number, stops: StopETA[]) => {
      checkAlerts(routeId, direction, stops);
      // Fill in second provider IDs from fallback response
      if (stops.length === 0 || !stops[0].source) return;
      const source = stops[0].source;
      for (const f of favoritesRef.current) {
        if (f.routeId !== routeId || f.direction !== direction) continue;
        // Skip if provider IDs already populated
        const hasIds = source === "tdx" ? f.tdxRouteId && f.tdxStopId : f.ebusRouteId && f.ebusStopId;
        if (hasIds) continue;
        const matched = stops.find((s) => s.stopName === f.stopName);
        if (matched) {
          updateFavoriteIds(f.routeId, f.direction, f.stopId, source, routeId, matched.stopId);
        }
      }
    },
    [checkAlerts, updateFavoriteIds],
  );

  // Proactive resolve: on mount, check if favorites need ID conversion for current provider
  const resolvedRef = useRef(false);
  useEffect(() => {
    if (resolvedRef.current) return;
    const currentFavs = favoritesRef.current;
    if (currentFavs.length === 0) return;
    resolvedRef.current = true;
    let cancelled = false;

    const resolveAll = async () => {
      // Group favorites by routeName + direction
      const groups = new Map<string, typeof currentFavs>();
      for (const f of currentFavs) {
        const key = `${f.routeName}:${f.direction}`;
        if (!groups.has(key)) groups.set(key, []);
        groups.get(key)!.push(f);
      }

      const resolveGroup = async (group: typeof currentFavs) => {
        const sample = group[0];
        const routes = await searchRoutes(sample.routeName);
        if (cancelled) return;
        const matched = routes.find((r) => r.routeName === sample.routeName);
        if (!matched?.source || group.every((f) => f.routeId === matched.routeId)) return;

        const stops = await getStops(matched.routeId, sample.direction);
        if (cancelled) return;

        for (const fav of group) {
          if (fav.routeId === matched.routeId) continue;
          const normalizedFavStop = normalizeName(fav.stopName);
          const stop = stops.find(
            (s) => normalizeName(s.stopName) === normalizedFavStop,
          );
          if (stop) {
            resolveFavorite(
              fav.routeId, fav.direction, fav.stopId,
              matched.source, matched.routeId, stop.stopId,
            );
          }
        }
      };

      await Promise.allSettled(
        Array.from(groups.values()).map((group) => resolveGroup(group)),
      );
    };

    resolveAll();
    return () => { cancelled = true; };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const favoritesEta = useFavoritesEta(favorites, handleEtaFetched, resolveFavorite);

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

      {favorites.length === 0 && (
        <p className="mt-6 text-sm text-gray-500">
          尚無收藏站點，搜尋路線後可收藏常用站點，即時到站資訊會顯示在這裡。
        </p>
      )}

      {favorites.length > 0 && (
        <section className="mt-6">
          <h2 className="mb-3 text-lg font-semibold">收藏站點</h2>
          <ul className="grid grid-cols-1 gap-3 md:grid-cols-2" role="list">
            {favoritesEta.map(({ favorite: f, eta }) => {
              const favKey = `${f.routeId}:${f.direction}:${f.stopId}`;
              const alert = getAlert(f.routeId, f.direction, f.stopId);
              const showMenu = alertMenuKey === favKey;

              return (
                <li
                  key={favKey}
                  className="rounded-lg border border-gray-200 px-4 py-3 dark:border-gray-700"
                >
                  <div className="flex items-center gap-3">
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
                    <AlertBell
                      alert={alert}
                      onClick={() =>
                        alert
                          ? removeAlert(f.routeId, f.direction, f.stopId)
                          : setAlertMenuKey(showMenu ? null : favKey)
                      }
                    />
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
                  </div>
                  {showMenu && (
                    <AlertMenu
                      className="ml-2"
                      onSelect={(min) => {
                        addAlert({
                          routeId: f.routeId,
                          routeName: f.routeName,
                          direction: f.direction,
                          stopId: f.stopId,
                          stopName: f.stopName,
                          thresholdMinutes: min,
                        });
                        setAlertMenuKey(null);
                      }}
                    />
                  )}
                </li>
              );
            })}
          </ul>
        </section>
      )}
    </div>
  );
}
