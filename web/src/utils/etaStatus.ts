/** Negative ETA values are status codes, not seconds. Mirrors internal/model/model.go. */
export const ETA_NOT_DEPARTED = -1;
export const ETA_LAST_BUS_LEFT = -2;
export const ETA_NO_STOP = -3;
export const ETA_NOT_OPERATING = -4;

/** Upper bound (seconds) of the "arrived" tier — the bus is pulling in. */
export const ARRIVED_MAX_SEC = 90;
/** Upper bound (exclusive, seconds) of the "arriving" tier. */
export const ARRIVING_MAX_SEC = 180;

/**
 * Converts an ETA value to display text.
 * Keep in sync with model.ETAStatus in internal/model/eta.go — both sides
 * are covered by the same boundary cases (0, 90, 91, 179, 180, 181, -1..-4).
 */
export function etaStatus(eta: number): string {
  switch (eta) {
    case ETA_NOT_DEPARTED:
      return "未發車";
    case ETA_LAST_BUS_LEFT:
      return "末班車已駛離";
    case ETA_NO_STOP:
      return "交管不停靠";
    case ETA_NOT_OPERATING:
      return "未營運";
  }
  if (eta < 0) return "未知";
  if (eta <= ARRIVED_MAX_SEC) return "進站中";
  if (eta < ARRIVING_MAX_SEC) return "將到站";
  return `約${Math.ceil(eta / 60)}分`;
}
