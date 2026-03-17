import { useState, useEffect, useRef, useMemo } from "react";
import { getETA, searchRoutes, getStops } from "../api/client";
import type { Favorite, StopETA } from "../api/types";
import { normalizeName } from "../utils/normalize";

const POLL_INTERVAL = 15_000;

export interface FavoriteETA {
  favorite: Favorite;
  eta?: StopETA;
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
): FavoriteETA[] {
  const [etaMap, setEtaMap] = useState<Map<string, StopETA>>(new Map());
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

      results.forEach((result, i) => {
        const rk = routeKeys[i];
        if (result.status === "fulfilled" && result.value.stops.length > 0) {
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

      // Only update state if data actually changed to avoid unnecessary re-renders
      setEtaMap((prev) => {
        if (prev.size !== newMap.size) return newMap;
        for (const [key, stop] of newMap) {
          const old = prev.get(key);
          if (!old || old.eta !== stop.eta || old.status !== stop.status) return newMap;
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

  return favorites.map((f) => ({
    favorite: f,
    eta:
      etaMap.get(`${f.routeId}:${f.direction}:${f.stopId}`) ??
      etaMap.get(`${f.routeId}:${f.direction}:name:${f.stopName}`),
  }));
}
