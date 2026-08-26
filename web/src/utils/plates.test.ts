import { describe, it, expect } from "vitest";
import { plateText } from "./plates";

describe("plateText", () => {
  it("joins plates with a comma", () => {
    expect(plateText([{ plateNumb: "A" }, { plateNumb: "B" }])).toBe("A, B");
  });

  it("returns an empty string for no buses", () => {
    expect(plateText([])).toBe("");
    expect(plateText(undefined)).toBe("");
  });
});
