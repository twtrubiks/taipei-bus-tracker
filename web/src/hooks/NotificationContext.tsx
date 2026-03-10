import { createContext, useContext, type ReactNode } from "react";
import { useNotification } from "./useNotification";

type NotificationContextValue = ReturnType<typeof useNotification>;

const NotificationContext = createContext<NotificationContextValue | null>(null);

export function NotificationProvider({ children }: { children: ReactNode }) {
  const notification = useNotification();
  return (
    <NotificationContext.Provider value={notification}>
      {children}
    </NotificationContext.Provider>
  );
}

// eslint-disable-next-line react-refresh/only-export-components
export function useNotificationContext(): NotificationContextValue {
  const ctx = useContext(NotificationContext);
  if (!ctx) {
    throw new Error("useNotificationContext must be used within NotificationProvider");
  }
  return ctx;
}
