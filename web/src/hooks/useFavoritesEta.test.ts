import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
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
        status: "約2分",
        buses: [],
        source: "tdx",
      },
      {
        stopId: "s2",
        stopName: "板橋車站",
        sequence: 2,
        eta: 300,
        status: "約5分",
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
      expect(result.current[0].eta).toBeDefined();
    });

    expect(result.current[0].eta?.stopId).toBe("s1");
  });
});
