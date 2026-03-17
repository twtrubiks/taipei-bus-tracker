import { useState, useEffect, useRef } from "react";
import { getETA } from "../api/client";
import type { ETAResponse } from "../api/types";

const POLL_INTERVAL = 15_000;

export function useEta(routeId: string, direction: number) {
  const [data, setData] = useState<ETAResponse | null>(null);
  const [error, setError] = useState<Error | null>(null);
  const timerRef = useRef<ReturnType<typeof setInterval>>(undefined);

  useEffect(() => {
    if (!routeId) return;
    let cancelled = false;

    const doFetch = () =>
      getETA(routeId, direction)
        .then((res) => {
          if (!cancelled) {
            setData(res);
            setError(null);
          }
        })
        .catch((err) => {
          if (!cancelled) {
            // Keep previous data visible; only set error indicator
            setError(err instanceof Error ? err : new Error(String(err)));
          }
        });

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
  }, [routeId, direction]);

  return { data, error };
}
