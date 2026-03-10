import { useState, useEffect, useCallback } from "react";

const STORAGE_KEY = "bus-dark-mode";

type Mode = "light" | "dark" | "system";

function getSystemPreference(): boolean {
  return window.matchMedia("(prefers-color-scheme: dark)").matches;
}

function loadMode(): Mode {
  const stored = localStorage.getItem(STORAGE_KEY);
  if (stored === "light" || stored === "dark") return stored;
  return "system";
}

function applyDark(isDark: boolean): void {
  document.documentElement.classList.toggle("dark", isDark);
}

export function useDarkMode() {
  const [mode, setModeState] = useState<Mode>(loadMode);

  const isDark = mode === "system" ? getSystemPreference() : mode === "dark";

  useEffect(() => {
    applyDark(isDark);
  }, [isDark]);

  // Listen for system preference changes
  useEffect(() => {
    if (mode !== "system") return;
    const mq = window.matchMedia("(prefers-color-scheme: dark)");
    const handler = (e: MediaQueryListEvent) => applyDark(e.matches);
    mq.addEventListener("change", handler);
    return () => mq.removeEventListener("change", handler);
  }, [mode]);

  const setMode = useCallback((newMode: Mode) => {
    setModeState(newMode);
    if (newMode === "system") {
      localStorage.removeItem(STORAGE_KEY);
    } else {
      localStorage.setItem(STORAGE_KEY, newMode);
    }
  }, []);

  const toggle = useCallback(() => {
    setMode(isDark ? "light" : "dark");
  }, [isDark, setMode]);

  return { mode, isDark, setMode, toggle };
}
