import { useMemo } from "react";
import type { Stop, StopETA } from "../api/types";

interface Props {
  stops: Stop[];
  etas: StopETA[];
}

function statusColor(eta: number): string {
  if (eta >= 0 && eta <= 180) return "text-green-600 font-bold";
  if (eta > 180) return "text-blue-600";
  return "text-gray-400";
}

export default function StopList({ stops, etas }: Props) {
  const etaMap = useMemo(
    () => new Map(etas.map((e) => [e.sequence, e])),
    [etas],
  );

  if (stops.length === 0) {
    return <p className="mt-4 text-gray-500">無站點資料</p>;
  }

  return (
    <ul className="mt-4 divide-y" role="list">
      {stops.map((stop) => {
        const eta = etaMap.get(stop.sequence);
        return (
          <li key={stop.sequence} className="flex items-center gap-3 py-3">
            <span className="w-6 text-center text-xs text-gray-400">
              {stop.sequence}
            </span>
            <span className="flex-1">{stop.stopName}</span>
            <div className="text-right">
              <span className={`text-sm ${statusColor(eta?.eta ?? -999)}`}>
                {eta?.status ?? "—"}
              </span>
              {eta?.buses && eta.buses.length > 0 && (
                <p className="text-xs text-gray-400">
                  {eta.buses.map((b) => b.plateNumb).join(", ")}
                </p>
              )}
            </div>
          </li>
        );
      })}
    </ul>
  );
}
