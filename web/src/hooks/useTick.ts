import { useState, useEffect } from "react";

/**
 * Re-renders the caller once per interval while the page is visible, so values
 * derived from the current time stay fresh. Costs no network traffic.
 *
 * Returns the current timestamp in milliseconds. Ticking pauses while the page is
 * hidden and resumes with an immediate update, matching the polling hooks.
 */
export function useTick(intervalMs = 1000): number {
  const [now, setNow] = useState(() => Date.now());

  useEffect(() => {
    let timer: ReturnType<typeof setInterval> | undefined;

    const start = () => {
      clearInterval(timer);
      setNow(Date.now());
      timer = setInterval(() => setNow(Date.now()), intervalMs);
    };

    const handleVisibility = () => {
      if (document.visibilityState === "visible") {
        start();
      } else {
        clearInterval(timer);
      }
    };

    if (document.visibilityState === "visible") start();
    document.addEventListener("visibilitychange", handleVisibility);

    return () => {
      clearInterval(timer);
      document.removeEventListener("visibilitychange", handleVisibility);
    };
  }, [intervalMs]);

  return now;
}
