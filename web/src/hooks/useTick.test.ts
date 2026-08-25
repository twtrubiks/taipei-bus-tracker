import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { renderHook, act } from "@testing-library/react";
import { useTick } from "./useTick";

beforeEach(() => {
  vi.useFakeTimers();
  Object.defineProperty(document, "visibilityState", {
    writable: true,
    value: "visible",
  });
});

afterEach(() => {
  vi.useRealTimers();
});

function setVisibility(value: "visible" | "hidden") {
  act(() => {
    Object.defineProperty(document, "visibilityState", { writable: true, value });
    document.dispatchEvent(new Event("visibilitychange"));
  });
}

describe("useTick", () => {
  it("advances once per second", () => {
    const { result } = renderHook(() => useTick());
    const start = result.current;

    act(() => {
      vi.advanceTimersByTime(1_000);
    });
    expect(result.current - start).toBe(1_000);

    act(() => {
      vi.advanceTimersByTime(3_000);
    });
    expect(result.current - start).toBe(4_000);
  });

  it("stops ticking while the page is hidden and resumes on return", () => {
    const { result } = renderHook(() => useTick());

    setVisibility("hidden");
    const frozen = result.current;

    act(() => {
      vi.advanceTimersByTime(10_000);
    });
    expect(result.current).toBe(frozen);

    setVisibility("visible");
    expect(result.current).toBe(frozen + 10_000);
  });

  it("stops the timer on unmount", () => {
    const { result, unmount } = renderHook(() => useTick());
    const last = result.current;
    unmount();

    act(() => {
      vi.advanceTimersByTime(5_000);
    });
    expect(result.current).toBe(last);
  });
});
