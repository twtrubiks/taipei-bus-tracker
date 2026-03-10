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

  const addFavorite = useCallback((fav: Favorite) => {
    setFavorites((prev) => {
      const key = favKey(fav);
      if (prev.some((f) => favKey(f) === key)) return prev;
      return [...prev, fav];
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

  const isFavorite = useCallback(
    (routeId: string, direction: number, stopId: string): boolean => {
      return favorites.some(
        (f) => favKey(f) === favKey({ routeId, direction, stopId }),
      );
    },
    [favorites],
  );

  return { favorites, addFavorite, removeFavorite, isFavorite };
}
