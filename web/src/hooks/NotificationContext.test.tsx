import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, act } from "@testing-library/react";
import { NotificationProvider, useNotificationContext } from "./NotificationContext";
import { setupMockNotification } from "../test/mockNotification";

let notificationSpy: ReturnType<typeof vi.fn>;

beforeEach(() => {
  localStorage.clear();
  notificationSpy = setupMockNotification();
});

function TestConsumer() {
  const { addAlert, checkAlerts, alerts, permissionDenied } =
    useNotificationContext();
  return (
    <div>
      <span data-testid="count">{alerts.length}</span>
      <span data-testid="denied">{String(permissionDenied)}</span>
      <button
        data-testid="add"
        onClick={() =>
          addAlert({
            routeId: "r1",
            routeName: "299",
            direction: 0,
            stopId: "s1",
            stopName: "台北車站",
            thresholdMinutes: 3,
          })
        }
      >
        add
      </button>
      <button
        data-testid="check"
        onClick={() =>
          checkAlerts("r1", 0, [
            {
              stopId: "s1",
              stopName: "台北車站",
              sequence: 1,
              eta: 120,
              status: "約2分",
              buses: [],
              source: "tdx",
            },
          ])
        }
      >
        check
      </button>
    </div>
  );
}

describe("NotificationContext", () => {
  it("shares notification state across consumers via provider", async () => {
    render(
      <NotificationProvider>
        <TestConsumer />
      </NotificationProvider>,
    );

    expect(screen.getByTestId("count").textContent).toBe("0");

    await act(async () => {
      screen.getByTestId("add").click();
    });
    expect(screen.getByTestId("count").textContent).toBe("1");

    act(() => {
      screen.getByTestId("check").click();
    });
    expect(notificationSpy).toHaveBeenCalledWith("299 - 台北車站", {
      body: "約 2 分鐘後到站",
      tag: "r1:0:s1",
    });
  });

  it("does not fire same alert twice (shared firedRef)", async () => {
    render(
      <NotificationProvider>
        <TestConsumer />
      </NotificationProvider>,
    );

    await act(async () => {
      screen.getByTestId("add").click();
    });

    act(() => {
      screen.getByTestId("check").click();
    });
    act(() => {
      screen.getByTestId("check").click();
    });

    expect(notificationSpy).toHaveBeenCalledTimes(1);
  });

  it("throws when used outside provider", () => {
    const spy = vi.spyOn(console, "error").mockImplementation(() => {});
    expect(() => render(<TestConsumer />)).toThrow(
      "useNotificationContext must be used within NotificationProvider",
    );
    spy.mockRestore();
  });
});
