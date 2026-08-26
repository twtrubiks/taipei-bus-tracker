import type { Bus } from "../api/types";
import { plateText } from "../utils/plates";

/**
 * Where a bus sits relative to the stop it is attached to.
 * - "at-stop": standing at the stop (eBus `bi`, or a TDX plate reported on the stop)
 * - "between": already left the stop, on the way to the next one (eBus `bo`)
 */
export type BusPlacement = "at-stop" | "between";

/** Inline SVG so the icon inherits currentColor and needs no icon package. */
function BusIcon({ className }: { className?: string }) {
  return (
    <svg viewBox="0 0 24 24" aria-hidden="true" focusable="false" className={className} fill="currentColor">
      <path d="M6 2h12a3 3 0 0 1 3 3v10a3 3 0 0 1-1 2.2V20a1 1 0 0 1-1 1h-1.5a1 1 0 0 1-1-1v-1H7.5v1a1 1 0 0 1-1 1H5a1 1 0 0 1-1-1v-2.8A3 3 0 0 1 3 15V5a3 3 0 0 1 3-3Zm-1 4v5h14V6H5Zm2.5 7.5a1.25 1.25 0 1 0 0 2.5 1.25 1.25 0 0 0 0-2.5Zm9 0a1.25 1.25 0 1 0 0 2.5 1.25 1.25 0 0 0 0-2.5Z" />
    </svg>
  );
}

interface Props {
  buses: Bus[] | undefined;
  placement: BusPlacement;
  /** Stop the buses are attached to; makes the accessible label locatable. */
  stopName: string;
}

/**
 * The bus icon drawn on the timeline — on a stop node for "at-stop", on the
 * segment towards the next stop for "between". Carries the accessible label
 * (plate numbers + position); the visible plate text next to it is decorative.
 * Returns null when there is no bus, so callers need no guard.
 */
export default function BusMarker({ buses, placement, stopName }: Props) {
  if (!buses || buses.length === 0) return null;

  const plates = plateText(buses);
  const label =
    placement === "at-stop"
      ? `${plates} 停靠於 ${stopName}`
      : `${plates} 已駛離 ${stopName}，前往下一站`;

  return (
    <span
      data-testid={`bus-${placement}`}
      title={label}
      className={
        placement === "at-stop"
          ? "block rounded-full bg-white p-0.5 text-green-600 dark:bg-gray-900 dark:text-green-400"
          : "block rounded-full bg-white p-0.5 text-gray-500 dark:bg-gray-900 dark:text-gray-400"
      }
    >
      <BusIcon className="h-4 w-4" />
      <span className="sr-only">{label}</span>
    </span>
  );
}
