import { useState, useEffect, useRef, useMemo } from "react";
import { getETA } from "../api/client";
import type { Favorite, StopETA } from "../api/types";

const POLL_INTERVAL = 15_000;

export interface FavoriteETA {
  favorite: Favorite;
  eta?: StopETA;
}

/**
 * Fetches ETA data for a list of favorites by batching requests per route+direction.
 */
export function useFavoritesEta(favorites: Favorite[]): FavoriteETA[] {
  const [etaMap, setEtaMap] = useState<Map<string, StopETA>>(new Map());
  const timerRef = useRef<ReturnType<typeof setInterval>>();

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

    const doFetch = async () => {
      const newMap = new Map<string, StopETA>();
      const results = await Promise.allSettled(
        routeKeys.map((rk) => getETA(rk.routeId, rk.direction)),
      );
      if (cancelled) return;
      results.forEach((result, i) => {
        if (result.status === "fulfilled") {
          const rk = routeKeys[i];
          for (const stop of result.value.stops) {
            newMap.set(`${rk.routeId}:${rk.direction}:${stop.stopId}`, stop);
          }
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
    };

    doFetch();
    timerRef.current = setInterval(doFetch, POLL_INTERVAL);

    const handleVisibility = () => {
      if (document.visibilityState === "hidden") {
        clearInterval(timerRef.current);
      } else {
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
    eta: etaMap.get(`${f.routeId}:${f.direction}:${f.stopId}`),
  }));
}
