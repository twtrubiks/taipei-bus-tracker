import type { Bus } from "../api/types";

/** Formats a group of buses as a single comma-separated plate-number label. */
export function plateText(buses: Bus[] | undefined): string {
  return (buses ?? []).map((b) => b.plateNumb).join(", ");
}
