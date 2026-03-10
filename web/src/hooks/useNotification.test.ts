import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, act } from "@testing-library/react";
import { useNotification } from "./useNotification";
import { setupMockNotification } from "../test/mockNotification";
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

const defaultAlert = {
  routeId: "r1",
  routeName: "299",
  direction: 0,
  stopId: "s1",
  stopName: "台北車站",
  thresholdMinutes: 3,
};

let notificationSpy: ReturnType<typeof vi.fn>;

beforeEach(() => {
  localStorage.clear();
  notificationSpy = setupMockNotification();
});

describe("useNotification", () => {
  it("fires notification when ETA <= threshold", async () => {
    const { result } = renderHook(() => useNotification());

    await act(async () => {
      await result.current.addAlert(defaultAlert);
    });

    // ETA = 120 seconds (2 min) <= 3 min threshold → should fire
    act(() => {
      result.current.checkAlerts("r1", 0, [mockEta("s1", 120)]);
    });

    expect(notificationSpy).toHaveBeenCalledWith("299 - 台北車站", {
      body: "約 2 分鐘後到站",
      tag: "r1:0:s1",
    });
  });

  it("does not fire when ETA > threshold", async () => {
    const { result } = renderHook(() => useNotification());

    await act(async () => {
      await result.current.addAlert(defaultAlert);
    });

    // ETA = 300 seconds (5 min) > 3 min threshold → should NOT fire
    act(() => {
      result.current.checkAlerts("r1", 0, [mockEta("s1", 300)]);
    });

    expect(notificationSpy).not.toHaveBeenCalled();
  });

  it("does not fire twice for the same alert", async () => {
    const { result } = renderHook(() => useNotification());

    await act(async () => {
      await result.current.addAlert(defaultAlert);
    });

    act(() => {
      result.current.checkAlerts("r1", 0, [mockEta("s1", 120)]);
    });
    act(() => {
      result.current.checkAlerts("r1", 0, [mockEta("s1", 60)]);
    });

    expect(notificationSpy).toHaveBeenCalledTimes(1);
  });

  it("shows permissionDenied when permission is denied", async () => {
    setupMockNotification("default", "denied");

    const { result } = renderHook(() => useNotification());

    await act(async () => {
      await result.current.addAlert(defaultAlert);
    });

    expect(result.current.permissionDenied).toBe(true);
  });

  it("notification content includes correct route and stop name", async () => {
    const { result } = renderHook(() => useNotification());

    await act(async () => {
      await result.current.addAlert({
        ...defaultAlert,
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

    expect(notificationSpy).toHaveBeenCalledWith("307 - 板橋車站", {
      body: "約 3 分鐘後到站",
      tag: "r2:1:s5",
    });
  });

  it("does not fire for negative ETA (not running)", async () => {
    const { result } = renderHook(() => useNotification());

    await act(async () => {
      await result.current.addAlert(defaultAlert);
    });

    act(() => {
      result.current.checkAlerts("r1", 0, [mockEta("s1", -1)]);
    });

    expect(notificationSpy).not.toHaveBeenCalled();
  });

  it("fires again after removeAlert then re-addAlert", async () => {
    const { result } = renderHook(() => useNotification());

    await act(async () => {
      await result.current.addAlert(defaultAlert);
    });

    act(() => {
      result.current.checkAlerts("r1", 0, [mockEta("s1", 120)]);
    });
    expect(notificationSpy).toHaveBeenCalledTimes(1);

    // Remove and re-add → should be able to fire again
    act(() => {
      result.current.removeAlert("r1", 0, "s1");
    });
    await act(async () => {
      await result.current.addAlert(defaultAlert);
    });

    act(() => {
      result.current.checkAlerts("r1", 0, [mockEta("s1", 60)]);
    });
    expect(notificationSpy).toHaveBeenCalledTimes(2);
  });

  it("does not fire for mismatched routeId or direction", async () => {
    const { result } = renderHook(() => useNotification());

    await act(async () => {
      await result.current.addAlert(defaultAlert);
    });

    // Wrong routeId
    act(() => {
      result.current.checkAlerts("r999", 0, [mockEta("s1", 120)]);
    });
    expect(notificationSpy).not.toHaveBeenCalled();

    // Wrong direction
    act(() => {
      result.current.checkAlerts("r1", 1, [mockEta("s1", 120)]);
    });
    expect(notificationSpy).not.toHaveBeenCalled();
  });

  it("persists alerts to localStorage", async () => {
    const { result, unmount } = renderHook(() => useNotification());

    await act(async () => {
      await result.current.addAlert(defaultAlert);
    });

    unmount();

    // Re-mount → alerts should be loaded from localStorage
    const { result: result2 } = renderHook(() => useNotification());

    act(() => {
      result2.current.checkAlerts("r1", 0, [mockEta("s1", 120)]);
    });
    expect(notificationSpy).toHaveBeenCalledTimes(1);
  });

  it("does not fire when Notification API is unavailable", () => {
    // Seed localStorage with an alert, then remove Notification API
    localStorage.setItem("bus-notifications", JSON.stringify([defaultAlert]));

    // @ts-expect-error -- intentionally removing Notification for test
    delete window.Notification;

    const { result } = renderHook(() => useNotification());

    act(() => {
      result.current.checkAlerts("r1", 0, [mockEta("s1", 120)]);
    });

    expect(notificationSpy).not.toHaveBeenCalled();
  });
});
