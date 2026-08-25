import { describe, it, expect } from "vitest";
import { countdownEta, STALE_AFTER_MS } from "./countdown";

const T0 = 1_700_000_000_000;

describe("countdownEta", () => {
  it("subtracts elapsed whole seconds", () => {
    expect(countdownEta(300, T0, T0)).toBe(300);
    expect(countdownEta(300, T0, T0 + 1_000)).toBe(299);
    expect(countdownEta(300, T0, T0 + 10_000)).toBe(290);
  });

  it("ignores sub-second elapsed time", () => {
    expect(countdownEta(300, T0, T0 + 999)).toBe(300);
    expect(countdownEta(300, T0, T0 + 1_999)).toBe(299);
  });

  it("floors at zero instead of crossing into status-code territory", () => {
    expect(countdownEta(5, T0, T0 + 5_000)).toBe(0);
    expect(countdownEta(5, T0, T0 + 30_000)).toBe(0);
  });

  it("leaves status codes untouched", () => {
    for (const code of [-1, -2, -3, -4]) {
      expect(countdownEta(code, T0, T0 + 30_000)).toBe(code);
    }
  });

  it("freezes once the data is stale", () => {
    const atLimit = countdownEta(300, T0, T0 + STALE_AFTER_MS);
    expect(atLimit).toBe(300 - STALE_AFTER_MS / 1000);
    expect(countdownEta(300, T0, T0 + STALE_AFTER_MS * 5)).toBe(atLimit);
  });

  it("returns the raw value before the first fetch", () => {
    expect(countdownEta(300, 0, T0)).toBe(300);
  });

  it("does not count up when the clock goes backwards", () => {
    expect(countdownEta(300, T0, T0 - 10_000)).toBe(300);
  });
});
