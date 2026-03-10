import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { renderHook, act } from "@testing-library/react";
import { useEta } from "./useEta";
import * as client from "../api/client";
import type { ETAResponse } from "../api/types";

vi.mock("../api/client");

const mockResponse: ETAResponse = {
  route: "0100000100",
  direction: 0,
  source: "tdx",
  updatedAt: "2026-01-01T00:00:00Z",
  stops: [
    {
      stopId: "s1",
      stopName: "台北車站",
      sequence: 1,
      eta: 300,
      status: "約5分",
      buses: [],
      source: "tdx",
    },
  ],
};

beforeEach(() => {
  vi.useFakeTimers({ shouldAdvanceTime: true });
  vi.clearAllMocks();
  Object.defineProperty(document, "visibilityState", {
    writable: true,
    value: "visible",
  });
});

afterEach(() => {
  vi.useRealTimers();
});

describe("useEta", () => {
  it("fetches ETA immediately on mount", async () => {
    vi.mocked(client.getETA).mockResolvedValue(mockResponse);

    const { result } = renderHook(() => useEta("0100000100", 0));

    // Flush the initial async fetch
    await act(async () => {
      await vi.advanceTimersByTimeAsync(0);
    });

    expect(client.getETA).toHaveBeenCalledWith("0100000100", 0);
    expect(result.current.data).toEqual(mockResponse);
  });

  it("polls every 15 seconds", async () => {
    vi.mocked(client.getETA).mockResolvedValue(mockResponse);

    renderHook(() => useEta("0100000100", 0));

    await act(async () => {
      await vi.advanceTimersByTimeAsync(0);
    });
    expect(client.getETA).toHaveBeenCalledTimes(1);

    await act(async () => {
      await vi.advanceTimersByTimeAsync(15_000);
    });
    expect(client.getETA).toHaveBeenCalledTimes(2);

    await act(async () => {
      await vi.advanceTimersByTimeAsync(15_000);
    });
    expect(client.getETA).toHaveBeenCalledTimes(3);
  });

  it("stops polling when page becomes hidden", async () => {
    vi.mocked(client.getETA).mockResolvedValue(mockResponse);

    renderHook(() => useEta("0100000100", 0));

    await act(async () => {
      await vi.advanceTimersByTimeAsync(0);
    });
    expect(client.getETA).toHaveBeenCalledTimes(1);

    act(() => {
      Object.defineProperty(document, "visibilityState", {
        writable: true,
        value: "hidden",
      });
      document.dispatchEvent(new Event("visibilitychange"));
    });

    await act(async () => {
      await vi.advanceTimersByTimeAsync(15_000);
    });
    expect(client.getETA).toHaveBeenCalledTimes(1);
  });

  it("resumes polling when page becomes visible again", async () => {
    vi.mocked(client.getETA).mockResolvedValue(mockResponse);

    renderHook(() => useEta("0100000100", 0));

    await act(async () => {
      await vi.advanceTimersByTimeAsync(0);
    });
    expect(client.getETA).toHaveBeenCalledTimes(1);

    // Hide
    act(() => {
      Object.defineProperty(document, "visibilityState", {
        writable: true,
        value: "hidden",
      });
      document.dispatchEvent(new Event("visibilitychange"));
    });

    // Show again
    act(() => {
      Object.defineProperty(document, "visibilityState", {
        writable: true,
        value: "visible",
      });
      document.dispatchEvent(new Event("visibilitychange"));
    });

    await act(async () => {
      await vi.advanceTimersByTimeAsync(0);
    });
    expect(client.getETA).toHaveBeenCalledTimes(2);
  });

  it("sets error on API failure", async () => {
    vi.mocked(client.getETA).mockRejectedValue(new Error("network error"));

    const { result } = renderHook(() => useEta("0100000100", 0));

    await act(async () => {
      await vi.advanceTimersByTimeAsync(0);
    });

    expect(result.current.error).toBeTruthy();
    expect(result.current.error?.message).toBe("network error");
  });

  it("cleans up timer on unmount", async () => {
    vi.mocked(client.getETA).mockResolvedValue(mockResponse);

    const { unmount } = renderHook(() => useEta("0100000100", 0));

    await act(async () => {
      await vi.advanceTimersByTimeAsync(0);
    });
    expect(client.getETA).toHaveBeenCalledTimes(1);

    unmount();

    await act(async () => {
      await vi.advanceTimersByTimeAsync(15_000);
    });
    expect(client.getETA).toHaveBeenCalledTimes(1);
  });
});
