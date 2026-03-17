import { useState, useEffect, useRef, useCallback } from "react";
import type { StopETA } from "../api/types";

export interface NotificationAlert {
  routeId: string;
  routeName: string;
  direction: number;
  stopId: string;
  stopName: string;
  thresholdMinutes: number;
}

const STORAGE_KEY = "bus-notifications";

async function showNotification(title: string, options?: NotificationOptions): Promise<void> {
  try {
    const reg = await navigator.serviceWorker?.getRegistration();
    if (reg) {
      await reg.showNotification(title, options);
    } else {
      new Notification(title, options);
    }
  } catch (err) {
    console.warn("[Notification]", err instanceof Error ? err.message : err);
  }
}

function loadAlerts(): NotificationAlert[] {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) return [];
    return JSON.parse(raw) as NotificationAlert[];
  } catch {
    return [];
  }
}

function saveAlerts(alerts: NotificationAlert[]): void {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(alerts));
  } catch {
    // QuotaExceededError or restricted storage — silently ignore
  }
}

function alertKey(
  a: Pick<NotificationAlert, "routeId" | "direction" | "stopId">,
): string {
  return `${a.routeId}:${a.direction}:${a.stopId}`;
}

export function useNotification() {
  const [alerts, setAlerts] = useState<NotificationAlert[]>(loadAlerts);
  const [permissionDenied, setPermissionDenied] = useState(false);
  const firedRef = useRef<Set<string>>(new Set());
  const isInitial = useRef(true);

  useEffect(() => {
    if (isInitial.current) {
      isInitial.current = false;
      return;
    }
    saveAlerts(alerts);
  }, [alerts]);

  const requestPermission = useCallback(async (): Promise<boolean> => {
    if (!("Notification" in window)) {
      setPermissionDenied(true);
      return false;
    }
    if (Notification.permission === "granted") return true;
    if (Notification.permission === "denied") {
      setPermissionDenied(true);
      return false;
    }
    const result = await Notification.requestPermission();
    if (result === "denied") {
      setPermissionDenied(true);
      return false;
    }
    return result === "granted";
  }, []);

  const addAlert = useCallback(
    async (alert: NotificationAlert) => {
      const granted = await requestPermission();
      if (!granted) return;

      setAlerts((prev) => {
        const key = alertKey(alert);
        // Remove existing alert for same stop, replace with new threshold
        const filtered = prev.filter((a) => alertKey(a) !== key);
        return [...filtered, alert];
      });
      // Reset fired state for this alert
      firedRef.current.delete(alertKey(alert));
    },
    [requestPermission],
  );

  const removeAlert = useCallback(
    (routeId: string, direction: number, stopId: string) => {
      const key = alertKey({ routeId, direction, stopId });
      setAlerts((prev) => prev.filter((a) => alertKey(a) !== key));
      firedRef.current.delete(key);
    },
    [],
  );

  const getAlert = useCallback(
    (
      routeId: string,
      direction: number,
      stopId: string,
    ): NotificationAlert | undefined => {
      const key = alertKey({ routeId, direction, stopId });
      return alerts.find((a) => alertKey(a) === key);
    },
    [alerts],
  );

  /**
   * Call this with fresh ETA data to check if any alerts should fire.
   */
  const checkAlerts = useCallback(
    (routeId: string, direction: number, etas: StopETA[]) => {
      if (
        !("Notification" in window) ||
        Notification.permission !== "granted"
      ) {
        return;
      }

      const etaMap = new Map(etas.map((e) => [e.stopId, e]));

      for (const alert of alerts) {
        if (alert.routeId !== routeId || alert.direction !== direction) continue;
        const key = alertKey(alert);
        if (firedRef.current.has(key)) continue;

        const eta = etaMap.get(alert.stopId);
        if (!eta || eta.eta < 0) continue;

        const thresholdSeconds = alert.thresholdMinutes * 60;
        if (eta.eta <= thresholdSeconds) {
          firedRef.current.add(key);
          const minutes = Math.ceil(eta.eta / 60);
          showNotification(`${alert.routeName} - ${alert.stopName}`, {
            body: `約 ${minutes} 分鐘後到站`,
            tag: key,
          });
        }
      }
    },
    [alerts],
  );

  return {
    alerts,
    permissionDenied,
    addAlert,
    removeAlert,
    getAlert,
    checkAlerts,
  };
}
