import { useState, useEffect, useRef } from "react";
import { getETA } from "../api/client";
import type { ETAResponse } from "../api/types";

const POLL_INTERVAL = 15_000;

export function useEta(routeId: string, direction: number) {
  const [data, setData] = useState<ETAResponse | null>(null);
  const [error, setError] = useState<Error | null>(null);
  const timerRef = useRef<ReturnType<typeof setInterval>>();

  useEffect(() => {
    if (!routeId) return;

    const doFetch = () =>
      getETA(routeId, direction)
        .then((res) => {
          setData(res);
          setError(null);
        })
        .catch((err) => {
          setError(err instanceof Error ? err : new Error(String(err)));
        });

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
      clearInterval(timerRef.current);
      document.removeEventListener("visibilitychange", handleVisibility);
    };
  }, [routeId, direction]);

  return { data, error };
}
