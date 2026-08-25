import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { renderHook, waitFor, act } from "@testing-library/react";
import { getETA } from "../api/client";
import { useFavoritesEta } from "./useFavoritesEta";
import type { Favorite } from "../api/types";

vi.mock("../api/client", () => ({
  getETA: vi.fn().mockResolvedValue({
    route: "r1",
    direction: 0,
    source: "tdx",
    updatedAt: "2026-01-01T00:00:00Z",
    stops: [
      {
        stopId: "s1",
        stopName: "台北車站",
        sequence: 1,
        eta: 120,
        buses: [],
        source: "tdx",
      },
      {
        stopId: "s2",
        stopName: "板橋車站",
        sequence: 2,
        eta: 300,
        buses: [],
        source: "tdx",
      },
    ],
  }),
}));

const favorites: Favorite[] = [
  {
    routeId: "r1",
    routeName: "299",
    direction: 0,
    stopId: "s1",
    stopName: "台北車站",
    sequence: 1,
  },
];

beforeEach(() => {
  vi.clearAllMocks();
});

afterEach(() => {
  vi.useRealTimers();
});

describe("useFavoritesEta", () => {
  it("calls onEtaFetched callback with route+direction+stops after fetch", async () => {
    const onEtaFetched = vi.fn();

    renderHook(() => useFavoritesEta(favorites, onEtaFetched));

    await waitFor(() => {
      expect(onEtaFetched).toHaveBeenCalledTimes(1);
    });

    expect(onEtaFetched).toHaveBeenCalledWith("r1", 0, expect.arrayContaining([
      expect.objectContaining({ stopId: "s1" }),
      expect.objectContaining({ stopId: "s2" }),
    ]));
  });

  it("works without onEtaFetched callback", async () => {
    const { result } = renderHook(() => useFavoritesEta(favorites));

    await waitFor(() => {
      expect(result.current.items[0].eta).toBeDefined();
    });

    expect(result.current.items[0].eta?.stopId).toBe("s1");
  });

  it("stamps fetchedAt once a poll returns data", async () => {
    const before = Date.now();
    const { result } = renderHook(() => useFavoritesEta(favorites));

    expect(result.current.fetchedAt).toBe(0);

    await waitFor(() => {
      expect(result.current.fetchedAt).toBeGreaterThanOrEqual(before);
    });
  });

  it("keeps the previous ETA and flags the failure when a poll fails", async () => {
    // Fake timers must be in place before the hook registers its polling interval
    vi.useFakeTimers({ shouldAdvanceTime: true });
    const { result } = renderHook(() => useFavoritesEta(favorites));

    await act(async () => {
      await vi.advanceTimersByTimeAsync(0);
    });
    expect(result.current.items[0].eta).toBeDefined();
    expect(result.current.failing).toBe(false);
    const fetchedAtBefore = result.current.items[0].fetchedAt;

    vi.mocked(getETA).mockRejectedValueOnce(new Error("network error"));
    await act(async () => {
      await vi.advanceTimersByTimeAsync(15_000);
    });

    expect(result.current.failing).toBe(true);
    // Stale numbers with a freshness warning beat a row of dashes
    expect(result.current.items[0].eta?.stopId).toBe("s1");
    // ...and they must still be counted down from when they were actually fetched
    expect(result.current.items[0].fetchedAt).toBe(fetchedAtBefore);
  });
});
