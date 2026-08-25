import { describe, it, expect } from "vitest";
import { etaStatus } from "./etaStatus";
import { statusColor } from "./statusColor";

describe("etaStatus", () => {
  it.each([
    [0, "進站中"],
    [60, "進站中"],
    [90, "進站中"],
    [91, "將到站"],
    [179, "將到站"],
    [180, "約3分"],
    [181, "約4分"],
    [300, "約5分"],
    [600, "約10分"],
  ])("maps %d seconds to %s", (eta, expected) => {
    expect(etaStatus(eta)).toBe(expected);
  });

  it.each([
    [-1, "未發車"],
    [-2, "末班車已駛離"],
    [-3, "交管不停靠"],
    [-4, "未營運"],
    [-99, "未知"],
  ])("maps status code %d to %s", (eta, expected) => {
    expect(etaStatus(eta)).toBe(expected);
  });
});

describe("statusColor", () => {
  it("uses the boldest tier while the bus is pulling in", () => {
    expect(statusColor(0)).toBe("text-green-600 font-bold");
    expect(statusColor(90)).toBe("text-green-600 font-bold");
  });

  it("separates arriving from arrived", () => {
    expect(statusColor(91)).toBe("text-green-600");
    expect(statusColor(179)).toBe("text-green-600");
  });

  it("uses the on-the-way tier from 180 seconds up", () => {
    expect(statusColor(180)).toBe("text-blue-600");
    expect(statusColor(600)).toBe("text-blue-600");
  });

  it("greys out status codes", () => {
    expect(statusColor(-1)).toBe("text-gray-400");
    expect(statusColor(-999)).toBe("text-gray-400");
  });
});
