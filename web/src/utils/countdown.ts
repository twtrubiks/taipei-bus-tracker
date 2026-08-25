/** Data older than this stops counting down — a frozen value beats an invented one. */
export const STALE_AFTER_MS = 60_000;

/**
 * Counts an ETA down by the time elapsed since it was fetched, so the display keeps
 * moving between polls without any extra requests.
 *
 * Negative values are status codes rather than seconds, and pass through untouched.
 * Elapsed time is capped at STALE_AFTER_MS: once the data is that old the countdown
 * freezes rather than inventing an arrival the upstream never reported.
 */
export function countdownEta(eta: number, fetchedAt: number, now: number): number {
  if (eta < 0 || !fetchedAt) return eta;
  const elapsedMs = Math.min(Math.max(now - fetchedAt, 0), STALE_AFTER_MS);
  return Math.max(0, eta - Math.floor(elapsedMs / 1000));
}
