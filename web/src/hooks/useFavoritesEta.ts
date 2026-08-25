import { useState, useEffect, useRef, useMemo } from "react";
import { getETA, searchRoutes, getStops } from "../api/client";
import type { Favorite, StopETA } from "../api/types";
import { normalizeName } from "../utils/normalize";

const POLL_INTERVAL = 15_000;

export interface FavoriteETA {
  favorite: Favorite;
  eta?: StopETA;
  /** When this favorite's route last refreshed, 0 if it never has. */
  fetchedAt: number;
}

export interface FavoritesEtaResult {
  items: FavoriteETA[];
  /**
   * Timestamp (ms) of the oldest data currently on screen — the one that decides
   * whether anything shown has gone stale. 0 before the first successful poll.
   */
  fetchedAt: number;
  /** Whether the most recent poll failed for at least one route. */
  failing: boolean;
}

/**
 * Fetches ETA data for a list of favorites by batching requests per route+direction.
 * Optional onEtaFetched callback is invoked per route+direction with fresh ETA data.
 * Optional resolveFavorite callback is called when ETA fails and lazy resolve finds new IDs.
 */
export function useFavoritesEta(
  favorites: Favorite[],
  onEtaFetched?: (routeId: string, direction: number, stops: StopETA[]) => void,
  resolveFavorite?: (
    oldRouteId: string, direction: number, oldStopId: string,
    source: string, newRouteId: string, newStopId: string,
  ) => void,
): FavoritesEtaResult {
  const [etaMap, setEtaMap] = useState<Map<string, StopETA>>(new Map());
  // Per route+direction, so a route that stops refreshing does not get counted
  // down from another route's fresh timestamp.
  const [fetchedAtMap, setFetchedAtMap] = useState<Map<string, number>>(new Map());
  const [failing, setFailing] = useState(false);
  const timerRef = useRef<ReturnType<typeof setInterval>>(undefined);
  const onEtaFetchedRef = useRef(onEtaFetched);
  const resolveFavoriteRef = useRef(resolveFavorite);
  const favoritesRef = useRef(favorites);
  useEffect(() => {
    onEtaFetchedRef.current = onEtaFetched;
  });
  useEffect(() => {
    resolveFavoriteRef.current = resolveFavorite;
  });
  useEffect(() => {
    favoritesRef.current = favorites;
  });

  // Derive unique route+direction combos
  const routeKeys = useMemo(() => {
    const seen = new Set<string>();
    const keys: { routeId: string; direction: number }[] = [];
    for (const f of favorites) {
      const k = `${f.routeId}:${f.direction}`;
      if (!seen.has(k)) {
        seen.add(k);
        keys.push({ routeId: f.routeId, direction: f.direction });
      }
    }
    return keys;
  }, [favorites]);

  useEffect(() => {
    if (routeKeys.length === 0) return;
    let cancelled = false;
    const resolved = new Set<string>();
    const inFlight = new Set<string>();

    const tryResolve = async (rk: { routeId: string; direction: number }) => {
      const resolveKey = `${rk.routeId}:${rk.direction}`;
      if (resolved.has(resolveKey) || inFlight.has(resolveKey) || !resolveFavoriteRef.current) return;
      inFlight.add(resolveKey);

      const affected = favoritesRef.current.filter(
        (f) => f.routeId === rk.routeId && f.direction === rk.direction,
      );
      if (affected.length === 0) return;

      try {
        const routes = await searchRoutes(affected[0].routeName);
        if (cancelled) return;
        const matched = routes.find((r) => r.routeName === affected[0].routeName);
        if (!matched) return;

        const stops = await getStops(matched.routeId, rk.direction);
        if (cancelled) return;

        for (const fav of affected) {
          const normalizedFavStop = normalizeName(fav.stopName);
          const stop = stops.find((s) => normalizeName(s.stopName) === normalizedFavStop);
          if (stop && resolveFavoriteRef.current && matched.source) {
            resolveFavoriteRef.current(
              fav.routeId, fav.direction, fav.stopId,
              matched.source, matched.routeId, stop.stopId,
            );
          }
        }
        resolved.add(resolveKey);
      } catch {
        // Resolve failed — allow retry on next poll
      } finally {
        inFlight.delete(resolveKey);
      }
    };

    const doFetch = async () => {
      const newMap = new Map<string, StopETA>();
      const results = await Promise.allSettled(
        routeKeys.map((rk) => getETA(rk.routeId, rk.direction)),
      );
      if (cancelled) return;

      const failedKeys: { routeId: string; direction: number }[] = [];
      const refreshed: string[] = [];

      results.forEach((result, i) => {
        const rk = routeKeys[i];
        if (result.status === "fulfilled" && result.value.stops.length > 0) {
          refreshed.push(`${rk.routeId}:${rk.direction}`);
          for (const stop of result.value.stops) {
            newMap.set(`${rk.routeId}:${rk.direction}:${stop.stopId}`, stop);
            // Secondary key by stopName for cross-provider fallback matching
            newMap.set(`${rk.routeId}:${rk.direction}:name:${stop.stopName}`, stop);
          }
          onEtaFetchedRef.current?.(rk.routeId, rk.direction, result.value.stops);
        } else {
          failedKeys.push(rk);
        }
      });

      setFailing(failedKeys.length > 0);

      if (refreshed.length > 0) {
        const at = Date.now();
        setFetchedAtMap((prev) => {
          const next = new Map(prev);
          for (const key of refreshed) next.set(key, at);
          return next;
        });
      }

      // Routes that failed this round keep their previous values — stale numbers
      // carrying a freshness warning beat a row of dashes.
      setEtaMap((prev) => {
        const merged = new Map<string, StopETA>();
        for (const [key, stop] of prev) {
          if (!refreshed.some((r) => key.startsWith(`${r}:`))) merged.set(key, stop);
        }
        for (const [key, stop] of newMap) merged.set(key, stop);

        // Only swap in a new map if something actually changed, to avoid re-renders
        if (prev.size !== merged.size) return merged;
        for (const [key, stop] of merged) {
          const old = prev.get(key);
          if (!old || old.eta !== stop.eta) return merged;
        }
        return prev;
      });

      // Lazy resolve for failed routes (async, non-blocking)
      for (const rk of failedKeys) {
        if (cancelled) return;
        tryResolve(rk);
      }
    };

    doFetch();
    timerRef.current = setInterval(doFetch, POLL_INTERVAL);

    const handleVisibility = () => {
      clearInterval(timerRef.current);
      if (document.visibilityState === "visible") {
        doFetch();
        timerRef.current = setInterval(doFetch, POLL_INTERVAL);
      }
    };

    document.addEventListener("visibilitychange", handleVisibility);

    return () => {
      cancelled = true;
      clearInterval(timerRef.current);
      document.removeEventListener("visibilitychange", handleVisibility);
    };
  }, [routeKeys]);

  const items = favorites.map((f) => ({
    favorite: f,
    eta:
      etaMap.get(`${f.routeId}:${f.direction}:${f.stopId}`) ??
      etaMap.get(`${f.routeId}:${f.direction}:name:${f.stopName}`),
    fetchedAt: fetchedAtMap.get(`${f.routeId}:${f.direction}`) ?? 0,
  }));

  const loaded = items.map((i) => i.fetchedAt).filter((at) => at > 0);
  const fetchedAt = loaded.length > 0 ? Math.min(...loaded) : 0;

  return { items, fetchedAt, failing };
}
