import type { NotificationAlert } from "../hooks/useNotification";

const ALERT_OPTIONS = [1, 3, 5];

interface AlertBellProps {
  alert: NotificationAlert | undefined;
  onClick: () => void;
}

export function AlertBell({ alert, onClick }: AlertBellProps) {
  return (
    <button
      type="button"
      aria-label={alert ? `已設定 ${alert.thresholdMinutes} 分鐘提醒` : "設定到站提醒"}
      className={`text-sm ${alert ? "text-amber-500" : "text-gray-400 hover:text-amber-500"}`}
      onClick={onClick}
    >
      {alert ? "\u{1F514}" : "\u{1F515}"}
    </button>
  );
}

interface AlertMenuProps {
  onSelect: (minutes: number) => void;
  className?: string;
}

export function AlertMenu({ onSelect, className }: AlertMenuProps) {
  return (
    <div className={`mt-2 flex gap-2 ${className ?? ""}`}>
      {ALERT_OPTIONS.map((min) => (
        <button
          key={min}
          type="button"
          className="rounded bg-amber-100 px-2 py-1 text-xs text-amber-700 hover:bg-amber-200"
          onClick={() => onSelect(min)}
        >
          {min} 分鐘前
        </button>
      ))}
    </div>
  );
}
