import { ARRIVED_MAX_SEC, ARRIVING_MAX_SEC } from "./etaStatus";

/** Colors follow the same three tiers as etaStatus: arrived / arriving / on the way. */
export function statusColor(eta: number): string {
  if (eta < 0) return "text-gray-400";
  if (eta <= ARRIVED_MAX_SEC) return "text-green-600 font-bold";
  if (eta < ARRIVING_MAX_SEC) return "text-green-600";
  return "text-blue-600";
}
