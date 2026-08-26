import { describe, it, expect, vi, afterEach } from "vitest";
import { render, screen, act } from "@testing-library/react";
import StopList from "./StopList";
import type { Stop, StopETA } from "../api/types";

const stops: Stop[] = [
  { stopId: "s1", stopName: "台北車站", sequence: 1 },
  { stopId: "s2", stopName: "中山站", sequence: 2 },
];

function makeEta(sequence: number, eta: number): StopETA {
  return {
    stopId: `s${sequence}`,
    stopName: `stop-${sequence}`,
    sequence,
    eta,
    buses: [],
    source: "tdx",
  };
}

describe("StopList ETA status rendering", () => {
  it("renders '約5分' for eta=300", () => {
    const etas = [makeEta(1, 300)];
    render(<StopList stops={stops} etas={etas} />);
    expect(screen.getByText("約5分")).toBeInTheDocument();
  });

  it("renders '進站中' for eta=60", () => {
    const etas = [makeEta(1, 60)];
    render(<StopList stops={stops} etas={etas} />);
    expect(screen.getByText("進站中")).toBeInTheDocument();
  });

  it("renders '將到站' for eta=120", () => {
    const etas = [makeEta(1, 120)];
    render(<StopList stops={stops} etas={etas} />);
    expect(screen.getByText("將到站")).toBeInTheDocument();
  });

  it("renders '未發車' for eta=-1", () => {
    const etas = [makeEta(1, -1)];
    render(<StopList stops={stops} etas={etas} />);
    expect(screen.getByText("未發車")).toBeInTheDocument();
  });

  it("renders '末班車已駛離' for eta=-2", () => {
    const etas = [makeEta(1, -2)];
    render(<StopList stops={stops} etas={etas} />);
    expect(screen.getByText("末班車已駛離")).toBeInTheDocument();
  });

  it("renders '交管不停靠' for eta=-3", () => {
    const etas = [makeEta(1, -3)];
    render(<StopList stops={stops} etas={etas} />);
    expect(screen.getByText("交管不停靠")).toBeInTheDocument();
  });

  it("renders '未營運' for eta=-4", () => {
    const etas = [makeEta(1, -4)];
    render(<StopList stops={stops} etas={etas} />);
    expect(screen.getByText("未營運")).toBeInTheDocument();
  });

  it("shows plate number when bus is present", () => {
    const etas: StopETA[] = [
      {
        ...makeEta(1, 60),
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

  it("matches ETA by stopId when sequence is 0 (TDX mode)", () => {
    const etas: StopETA[] = [
      { stopId: "s2", stopName: "中山站", sequence: 0, eta: 120, buses: [], source: "tdx" },
    ];
    render(<StopList stops={stops} etas={etas} />);
    // s2 = 中山站 should show 將到站, 台北車站 should show —
    expect(screen.getByText("將到站")).toBeInTheDocument();
    expect(screen.getAllByText("—")).toHaveLength(1);
  });

  it("renders empty state when stops is null", () => {
    render(<StopList stops={null as unknown as Stop[]} etas={[]} />);
    expect(screen.getByText("無站點資料")).toBeInTheDocument();
  });
});

describe("StopList countdown", () => {
  afterEach(() => {
    vi.useRealTimers();
  });

  it("counts down between polls and crosses into the next tier", () => {
    vi.useFakeTimers();
    const fetchedAt = Date.now();

    render(<StopList stops={stops} etas={[makeEta(1, 100)]} fetchedAt={fetchedAt} />);
    expect(screen.getByText("將到站")).toBeInTheDocument();

    act(() => {
      vi.advanceTimersByTime(15_000);
    });
    expect(screen.getByText("進站中")).toBeInTheDocument();
  });

  it("freezes the countdown once the data goes stale", () => {
    vi.useFakeTimers();
    const fetchedAt = Date.now();

    render(<StopList stops={stops} etas={[makeEta(1, 300)]} fetchedAt={fetchedAt} />);
    expect(screen.getByText("約5分")).toBeInTheDocument();

    // 60s elapsed = the staleness cap: 300 - 60 = 240s
    act(() => {
      vi.advanceTimersByTime(60_000);
    });
    expect(screen.getByText("約4分")).toBeInTheDocument();

    // Beyond the cap the value must not keep dropping
    act(() => {
      vi.advanceTimersByTime(120_000);
    });
    expect(screen.getByText("約4分")).toBeInTheDocument();
  });

  it("does not count down status codes", () => {
    vi.useFakeTimers();
    const fetchedAt = Date.now();

    render(<StopList stops={stops} etas={[makeEta(1, -1)]} fetchedAt={fetchedAt} />);
    expect(screen.getByText("未發車")).toBeInTheDocument();

    act(() => {
      vi.advanceTimersByTime(30_000);
    });
    expect(screen.getByText("未發車")).toBeInTheDocument();
  });
});

describe("StopList bus markers", () => {
  const threeStops: Stop[] = [
    { stopId: "s1", stopName: "台北車站", sequence: 1 },
    { stopId: "s2", stopName: "中山站", sequence: 2 },
    { stopId: "s3", stopName: "民權西路", sequence: 3 },
  ];

  it("shows a departed bus between its stop and the next one", () => {
    // The bug this guards: a bus that left stop 1 used to vanish entirely,
    // leaving stop 2 showing 「將到站」 with no plate anywhere.
    const etas: StopETA[] = [
      { ...makeEta(1, 300), departedBuses: [{ plateNumb: "KKB-1789" }] },
      makeEta(2, 120),
    ];
    render(<StopList stops={threeStops} etas={etas} />);

    expect(screen.getByTestId("bus-between")).toBeInTheDocument();
    expect(screen.queryByTestId("bus-at-stop")).not.toBeInTheDocument();
    expect(screen.getByText("KKB-1789")).toBeInTheDocument();
    expect(
      screen.getByText("KKB-1789 已駛離 台北車站，前往下一站"),
    ).toBeInTheDocument();
  });

  it("shows both markers when a stop has an arriving and a departed bus", () => {
    const etas: StopETA[] = [
      {
        ...makeEta(1, 30),
        buses: [{ plateNumb: "EAL-5812" }],
        departedBuses: [{ plateNumb: "KKA-0161" }, { plateNumb: "KKB-1785" }],
      },
    ];
    render(<StopList stops={threeStops} etas={etas} />);

    expect(screen.getByTestId("bus-at-stop")).toBeInTheDocument();
    expect(screen.getByTestId("bus-between")).toBeInTheDocument();
    expect(screen.getByText("EAL-5812")).toBeInTheDocument();
    expect(screen.getByText("KKA-0161, KKB-1785")).toBeInTheDocument();
  });

  it("still shows a departed bus on the final stop", () => {
    const etas: StopETA[] = [
      { ...makeEta(3, 60), departedBuses: [{ plateNumb: "EAL-5290" }] },
    ];
    render(<StopList stops={threeStops} etas={etas} />);

    expect(screen.getByTestId("bus-between")).toBeInTheDocument();
    expect(screen.getByText("EAL-5290")).toBeInTheDocument();
  });

  it("draws no bus icon when a stop has neither arriving nor departed buses", () => {
    render(<StopList stops={threeStops} etas={[makeEta(1, 300)]} />);

    expect(screen.queryByTestId("bus-at-stop")).not.toBeInTheDocument();
    expect(screen.queryByTestId("bus-between")).not.toBeInTheDocument();
  });

  it("handles a TDX response with no departedBuses field at all", () => {
    const etas: StopETA[] = [
      { ...makeEta(1, 60), buses: [{ plateNumb: "ABC-1234" }] },
    ];
    render(<StopList stops={threeStops} etas={etas} />);

    expect(screen.getByTestId("bus-at-stop")).toBeInTheDocument();
    expect(screen.queryByTestId("bus-between")).not.toBeInTheDocument();
  });
});
