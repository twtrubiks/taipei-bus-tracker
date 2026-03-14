import { useState, useCallback, useRef, useEffect } from "react";
import type { Favorite } from "../api/types";

const STORAGE_KEY = "bus-favorites";

function loadFavorites(): Favorite[] {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) return [];
    return JSON.parse(raw) as Favorite[];
  } catch {
    return [];
  }
}

function saveFavorites(favs: Favorite[]): void {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(favs));
}

function favKey(f: Pick<Favorite, "routeId" | "direction" | "stopId">): string {
  return `${f.routeId}:${f.direction}:${f.stopId}`;
}

export function useFavorites() {
  const [favorites, setFavorites] = useState<Favorite[]>(loadFavorites);
  const isInitial = useRef(true);

  useEffect(() => {
    if (isInitial.current) {
      isInitial.current = false;
      return;
    }
    saveFavorites(favorites);
  }, [favorites]);

  const addFavorite = useCallback((fav: Favorite, source?: string) => {
    setFavorites((prev) => {
      const key = favKey(fav);
      if (prev.some((f) => favKey(f) === key)) return prev;
      const enriched = { ...fav };
      if (source === "tdx") {
        enriched.tdxRouteId = fav.routeId;
        enriched.tdxStopId = fav.stopId;
      } else if (source === "ebus") {
        enriched.ebusRouteId = fav.routeId;
        enriched.ebusStopId = fav.stopId;
      }
      return [...prev, enriched];
    });
  }, []);

  const removeFavorite = useCallback(
    (routeId: string, direction: number, stopId: string) => {
      setFavorites((prev) =>
        prev.filter((f) => favKey(f) !== favKey({ routeId, direction, stopId })),
      );
    },
    [],
  );

  const updateFavoriteIds = useCallback(
    (routeId: string, direction: number, stopId: string, source: string, newRouteId: string, newStopId: string) => {
      setFavorites((prev) => {
        const key = favKey({ routeId, direction, stopId });
        const idx = prev.findIndex((f) => favKey(f) === key);
        if (idx === -1) return prev;
        const f = prev[idx];
        const field = source === "tdx" ? "tdxRouteId" : "ebusRouteId";
        const stopField = source === "tdx" ? "tdxStopId" : "ebusStopId";
        if (f[field] && f[stopField]) return prev; // already populated
        const updated = [...prev];
        updated[idx] = { ...f, [field]: newRouteId, [stopField]: newStopId };
        return updated;
      });
    },
    [],
  );

  const resolveFavorite = useCallback(
    (oldRouteId: string, direction: number, oldStopId: string, source: string, newRouteId: string, newStopId: string) => {
      setFavorites((prev) => {
        const key = favKey({ routeId: oldRouteId, direction, stopId: oldStopId });
        const idx = prev.findIndex((f) => favKey(f) === key);
        if (idx === -1) return prev;
        const f = prev[idx];
        const updated = [...prev];
        // Preserve old provider IDs (old IDs belong to the opposite provider)
        const oldFields = source === "ebus"
          ? { tdxRouteId: f.tdxRouteId || oldRouteId, tdxStopId: f.tdxStopId || oldStopId }
          : { ebusRouteId: f.ebusRouteId || oldRouteId, ebusStopId: f.ebusStopId || oldStopId };
        updated[idx] = {
          ...f,
          routeId: newRouteId,
          stopId: newStopId,
          ...oldFields,
          ...(source === "ebus"
            ? { ebusRouteId: newRouteId, ebusStopId: newStopId }
            : { tdxRouteId: newRouteId, tdxStopId: newStopId }),
        };
        return updated;
      });
    },
    [],
  );

  const isFavorite = useCallback(
    (routeId: string, direction: number, stopId: string): boolean => {
      return favorites.some(
        (f) => favKey(f) === favKey({ routeId, direction, stopId }),
      );
    },
    [favorites],
  );

  return { favorites, addFavorite, removeFavorite, isFavorite, updateFavoriteIds, resolveFavorite };
}
