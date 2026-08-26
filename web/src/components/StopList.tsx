import { useMemo, useState } from "react";
import type { Stop, StopETA } from "../api/types";
import type { NotificationAlert } from "../hooks/useNotification";
import { AlertBell, AlertMenu } from "./AlertButton";
import BusMarker from "./BusMarker";
import { plateText } from "../utils/plates";
import { statusColor } from "../utils/statusColor";
import { etaStatus } from "../utils/etaStatus";
import { countdownEta } from "../utils/countdown";
import { useTick } from "../hooks/useTick";

interface Props {
  stops: Stop[];
  etas: StopETA[];
  /** Timestamp (ms) the etas were fetched; drives the between-poll countdown. */
  fetchedAt?: number;
  routeId?: string;
  direction?: number;
  isFavorite?: (routeId: string, direction: number, stopId: string) => boolean;
  onToggleFavorite?: (stop: Stop) => void;
  getAlert?: (stopId: string) => NotificationAlert | undefined;
  onSetAlert?: (stop: Stop, minutes: number) => void;
  onRemoveAlert?: (stop: Stop) => void;
}

const RAIL = "absolute left-1/2 w-px -translate-x-1/2 bg-gray-200 dark:bg-gray-700";

/** Vertical extent of the rail within a row, given which ends it must reach. */
function railExtent(above: boolean, below: boolean): string {
  if (above && below) return "inset-y-0";
  if (below) return "top-1/2 bottom-0";
  if (above) return "top-0 h-1/2";
  return "hidden";
}

export default function StopList({
  stops,
  etas,
  fetchedAt = 0,
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
  const now = useTick();

  if (!stops || stops.length === 0) {
    return <p className="mt-4 text-gray-500">無站點資料</p>;
  }

  return (
    <ul className="mt-4" role="list">
      {stops.map((stop, idx) => {
        const eta = etaByStopId.get(stop.stopId) ?? etaBySeq.get(stop.sequence);
        const displayEta = eta ? countdownEta(eta.eta, fetchedAt, now) : undefined;
        const fav =
          routeId !== undefined &&
          direction !== undefined &&
          isFavorite?.(routeId, direction, stop.stopId);
        const alert = getAlert?.(stop.stopId);
        const showMenu = alertMenuStop === stop.stopId;

        const atStop = eta?.buses ?? [];
        // Buses that left this stop; they belong on the segment towards the next one.
        const departed = eta?.departedBuses ?? [];
        const isFirst = idx === 0;
        const isLast = idx === stops.length - 1;

        return (
          <li key={stop.sequence}>
            <div className="flex gap-3">
              <div className="relative w-8 shrink-0">
                <span
                  aria-hidden="true"
                  className={`${RAIL} ${railExtent(!isFirst, !isLast || departed.length > 0)}`}
                />
                <span className="absolute left-1/2 top-1/2 z-10 -translate-x-1/2 -translate-y-1/2">
                  {atStop.length > 0 ? (
                    <BusMarker buses={atStop} placement="at-stop" stopName={stop.stopName} />
                  ) : (
                    <span
                      aria-hidden="true"
                      className="block h-2.5 w-2.5 rounded-full bg-white ring-2 ring-gray-300 dark:bg-gray-900 dark:ring-gray-600"
                    />
                  )}
                </span>
              </div>

              <div
                className={`flex flex-1 items-center gap-3 py-3 ${
                  // No trailing rule under the final row of the list.
                  isLast && departed.length === 0
                    ? ""
                    : "border-b border-gray-100 dark:border-gray-800"
                }`}
              >
                <span className="w-6 shrink-0 text-center text-xs text-gray-400">
                  {stop.sequence}
                </span>
                <span className="min-w-0 flex-1 truncate" title={stop.stopName}>
                  {stop.stopName}
                </span>
                <div className="shrink-0 text-right">
                  <span className={`text-sm font-medium ${statusColor(displayEta ?? -999)}`}>
                    {displayEta === undefined ? "—" : etaStatus(displayEta)}
                  </span>
                  {atStop.length > 0 && (
                    <p aria-hidden="true" className="text-xs text-gray-500 dark:text-gray-400">
                      {plateText(atStop)}
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
                    {fav ? "★" : "☆"}
                  </button>
                )}
              </div>
            </div>

            {departed.length > 0 && (
              <div className="flex gap-3">
                <div className="relative w-8 shrink-0">
                  <span aria-hidden="true" className={`${RAIL} ${railExtent(true, !isLast)}`} />
                  <span className="absolute left-1/2 top-1/2 z-10 -translate-x-1/2 -translate-y-1/2">
                    <BusMarker buses={departed} placement="between" stopName={stop.stopName} />
                  </span>
                </div>
                <p aria-hidden="true" className="flex-1 py-2 text-xs text-gray-500 dark:text-gray-400">
                  {plateText(departed)}
                </p>
              </div>
            )}

            {showMenu && onSetAlert && (
              <AlertMenu
                className="ml-11"
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
