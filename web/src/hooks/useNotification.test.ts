import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, act } from "@testing-library/react";
import { useNotification } from "./useNotification";
import type { StopETA } from "../api/types";

const mockEta = (stopId: string, eta: number): StopETA => ({
  stopId,
  stopName: "台北車站",
  sequence: 1,
  eta,
  status: eta > 0 ? `約${Math.ceil(eta / 60)}分` : "未發車",
  buses: [],
  source: "tdx",
});

let notifConstructorSpy: ReturnType<typeof vi.fn>;

beforeEach(() => {
  localStorage.clear();
  notifConstructorSpy = vi.fn();

  // Mock Notification API using a class so `new Notification(...)` works
  class MockNotification {
    static permission = "granted";
    static requestPermission = vi.fn().mockResolvedValue("granted");
    constructor(...args: unknown[]) {
      notifConstructorSpy(...args);
    }
  }

  Object.defineProperty(window, "Notification", {
    writable: true,
    configurable: true,
    value: MockNotification,
  });
});

describe("useNotification", () => {
  it("fires notification when ETA <= threshold", async () => {
    const { result } = renderHook(() => useNotification());

    await act(async () => {
      await result.current.addAlert({
        routeId: "r1",
        routeName: "299",
        direction: 0,
        stopId: "s1",
        stopName: "台北車站",
        thresholdMinutes: 3,
      });
    });

    // ETA = 120 seconds (2 min) <= 3 min threshold → should fire
    act(() => {
      result.current.checkAlerts("r1", 0, [mockEta("s1", 120)]);
    });

    expect(notifConstructorSpy).toHaveBeenCalledWith("299 - 台北車站", {
      body: "約 2 分鐘後到站",
      tag: "r1:0:s1",
    });
  });

  it("does not fire when ETA > threshold", async () => {
    const { result } = renderHook(() => useNotification());

    await act(async () => {
      await result.current.addAlert({
        routeId: "r1",
        routeName: "299",
        direction: 0,
        stopId: "s1",
        stopName: "台北車站",
        thresholdMinutes: 3,
      });
    });

    // ETA = 300 seconds (5 min) > 3 min threshold → should NOT fire
    act(() => {
      result.current.checkAlerts("r1", 0, [mockEta("s1", 300)]);
    });

    expect(notifConstructorSpy).not.toHaveBeenCalled();
  });

  it("does not fire twice for the same alert", async () => {
    const { result } = renderHook(() => useNotification());

    await act(async () => {
      await result.current.addAlert({
        routeId: "r1",
        routeName: "299",
        direction: 0,
        stopId: "s1",
        stopName: "台北車站",
        thresholdMinutes: 3,
      });
    });

    act(() => {
      result.current.checkAlerts("r1", 0, [mockEta("s1", 120)]);
    });
    act(() => {
      result.current.checkAlerts("r1", 0, [mockEta("s1", 60)]);
    });

    expect(notifConstructorSpy).toHaveBeenCalledTimes(1);
  });

  it("shows permissionDenied when permission is denied", async () => {
    class DeniedNotification {
      static permission = "default";
      static requestPermission = vi.fn().mockResolvedValue("denied");
      constructor() {}
    }
    Object.defineProperty(window, "Notification", {
      writable: true,
      configurable: true,
      value: DeniedNotification,
    });

    const { result } = renderHook(() => useNotification());

    await act(async () => {
      await result.current.addAlert({
        routeId: "r1",
        routeName: "299",
        direction: 0,
        stopId: "s1",
        stopName: "台北車站",
        thresholdMinutes: 3,
      });
    });

    expect(result.current.permissionDenied).toBe(true);
  });

  it("notification content includes correct route and stop name", async () => {
    const { result } = renderHook(() => useNotification());

    await act(async () => {
      await result.current.addAlert({
        routeId: "r2",
        routeName: "307",
        direction: 1,
        stopId: "s5",
        stopName: "板橋車站",
        thresholdMinutes: 5,
      });
    });

    act(() => {
      result.current.checkAlerts("r2", 1, [mockEta("s5", 180)]);
    });

    expect(notifConstructorSpy).toHaveBeenCalledWith("307 - 板橋車站", {
      body: "約 3 分鐘後到站",
      tag: "r2:1:s5",
    });
  });

  it("does not fire for negative ETA (not running)", async () => {
    const { result } = renderHook(() => useNotification());

    await act(async () => {
      await result.current.addAlert({
        routeId: "r1",
        routeName: "299",
        direction: 0,
        stopId: "s1",
        stopName: "台北車站",
        thresholdMinutes: 3,
      });
    });

    act(() => {
      result.current.checkAlerts("r1", 0, [mockEta("s1", -1)]);
    });

    expect(notifConstructorSpy).not.toHaveBeenCalled();
  });
});
