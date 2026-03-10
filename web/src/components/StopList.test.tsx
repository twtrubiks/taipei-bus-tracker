import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import StopList from "./StopList";
import type { Stop, StopETA } from "../api/types";

const stops: Stop[] = [
  { stopId: "s1", stopName: "台北車站", sequence: 1 },
  { stopId: "s2", stopName: "中山站", sequence: 2 },
];

function makeEta(sequence: number, eta: number, status: string): StopETA {
  return {
    stopId: `s${sequence}`,
    stopName: `stop-${sequence}`,
    sequence,
    eta,
    status,
    buses: [],
    source: "tdx",
  };
}

describe("StopList ETA status rendering", () => {
  it("renders '約5分' for eta=300", () => {
    const etas = [makeEta(1, 300, "約5分")];
    render(<StopList stops={stops} etas={etas} />);
    expect(screen.getByText("約5分")).toBeInTheDocument();
  });

  it("renders '進站中' for eta=60", () => {
    const etas = [makeEta(1, 60, "進站中")];
    render(<StopList stops={stops} etas={etas} />);
    expect(screen.getByText("進站中")).toBeInTheDocument();
  });

  it("renders '未發車' for eta=-1", () => {
    const etas = [makeEta(1, -1, "未發車")];
    render(<StopList stops={stops} etas={etas} />);
    expect(screen.getByText("未發車")).toBeInTheDocument();
  });

  it("renders '末班車已駛離' for eta=-2", () => {
    const etas = [makeEta(1, -2, "末班車已駛離")];
    render(<StopList stops={stops} etas={etas} />);
    expect(screen.getByText("末班車已駛離")).toBeInTheDocument();
  });

  it("renders '交管不停靠' for eta=-3", () => {
    const etas = [makeEta(1, -3, "交管不停靠")];
    render(<StopList stops={stops} etas={etas} />);
    expect(screen.getByText("交管不停靠")).toBeInTheDocument();
  });

  it("renders '未營運' for eta=-4", () => {
    const etas = [makeEta(1, -4, "未營運")];
    render(<StopList stops={stops} etas={etas} />);
    expect(screen.getByText("未營運")).toBeInTheDocument();
  });

  it("shows plate number when bus is present", () => {
    const etas: StopETA[] = [
      {
        ...makeEta(1, 60, "進站中"),
        buses: [{ plateNumb: "ABC-1234" }],
      },
    ];
    render(<StopList stops={stops} etas={etas} />);
    expect(screen.getByText("ABC-1234")).toBeInTheDocument();
  });

  it("shows '—' when no ETA data for a stop", () => {
    render(<StopList stops={stops} etas={[]} />);
    const dashes = screen.getAllByText("—");
    expect(dashes).toHaveLength(2);
  });
});
