import { isStale } from "../utils/countdown";
import { useTick } from "../hooks/useTick";

function clock(ms: number): string {
  return new Date(ms).toLocaleTimeString("zh-TW", { hour12: false });
}

interface Props {
  /** Timestamp (ms) of the oldest data on screen, 0 if nothing has loaded yet. */
  fetchedAt: number;
  /** Whether the most recent poll failed. */
  failing: boolean;
}

/**
 * Tells the user how far to trust what is on screen. Stays silent while data is
 * fresh — a warning only earns its place once the numbers may no longer be true.
 */
export default function DataFreshness({ fetchedAt, failing }: Props) {
  const now = useTick();

  if (!fetchedAt) {
    return failing ? (
      <p className="mt-2 text-sm text-red-500" role="status">
        即時到站載入失敗，請檢查網路後重試
      </p>
    ) : null;
  }

  const stale = isStale(fetchedAt, now);
  if (!stale && !failing) return null;

  return (
    <p
      className={`mt-2 text-sm ${stale ? "text-red-500" : "text-amber-600"}`}
      role="status"
    >
      {stale ? "資料已延遲，倒數暫停" : "更新失敗，顯示先前的資料"}（最後更新{" "}
      {clock(fetchedAt)}）
    </p>
  );
}
