import { useMemo, useState } from "react";
import type { Stop, StopETA } from "../api/types";
import type { NotificationAlert } from "../hooks/useNotification";
import { AlertBell, AlertMenu } from "./AlertButton";
import { statusColor } from "../utils/statusColor";

interface Props {
  stops: Stop[];
  etas: StopETA[];
  routeId?: string;
  direction?: number;
  isFavorite?: (routeId: string, direction: number, stopId: string) => boolean;
  onToggleFavorite?: (stop: Stop) => void;
  getAlert?: (stopId: string) => NotificationAlert | undefined;
  onSetAlert?: (stop: Stop, minutes: number) => void;
  onRemoveAlert?: (stop: Stop) => void;
}

export default function StopList({
  stops,
  etas,
  routeId,
  direction,
  isFavorite,
  onToggleFavorite,
  getAlert,
  onSetAlert,
  onRemoveAlert,
}: Props) {
  // Match ETA to stops: prefer stopId (TDX), fall back to sequence (eBus)
  const etaByStopId = useMemo(
    () => new Map((etas ?? []).filter((e) => e.stopId).map((e) => [e.stopId, e])),
    [etas],
  );
  const etaBySeq = useMemo(
    () => new Map((etas ?? []).filter((e) => e.sequence > 0).map((e) => [e.sequence, e])),
    [etas],
  );
  const [alertMenuStop, setAlertMenuStop] = useState<string | null>(null);

  if (!stops || stops.length === 0) {
    return <p className="mt-4 text-gray-500">無站點資料</p>;
  }

  return (
    <ul className="mt-4 divide-y" role="list">
      {stops.map((stop) => {
        const eta = etaByStopId.get(stop.stopId) ?? etaBySeq.get(stop.sequence);
        const fav =
          routeId !== undefined &&
          direction !== undefined &&
          isFavorite?.(routeId, direction, stop.stopId);
        const alert = getAlert?.(stop.stopId);
        const showMenu = alertMenuStop === stop.stopId;

        return (
          <li key={stop.sequence} className="py-3">
            <div className="flex items-center gap-3">
              <span className="w-6 shrink-0 text-center text-xs text-gray-400">
                {stop.sequence}
              </span>
              <span className="min-w-0 flex-1 truncate" title={stop.stopName}>{stop.stopName}</span>
              <div className="shrink-0 text-right">
                <span className={`text-sm font-medium ${statusColor(eta?.eta ?? -999)}`}>
                  {eta?.status ?? "—"}
                </span>
                {eta?.buses && eta.buses.length > 0 && (
                  <p className="text-xs text-gray-400">
                    {eta.buses.map((b) => b.plateNumb).join(", ")}
                  </p>
                )}
              </div>
              {onSetAlert && (
                <AlertBell
                  alert={alert}
                  onClick={() =>
                    alert
                      ? onRemoveAlert?.(stop)
                      : setAlertMenuStop(showMenu ? null : stop.stopId)
                  }
                />
              )}
              {onToggleFavorite && (
                <button
                  type="button"
                  aria-label={fav ? "取消收藏" : "加入收藏"}
                  className="text-lg"
                  onClick={() => onToggleFavorite(stop)}
                >
                  {fav ? "\u2605" : "\u2606"}
                </button>
              )}
            </div>
            {showMenu && onSetAlert && (
              <AlertMenu
                className="ml-9"
                onSelect={(min) => {
                  onSetAlert(stop, min);
                  setAlertMenuStop(null);
                }}
              />
            )}
          </li>
        );
      })}
    </ul>
  );
}
