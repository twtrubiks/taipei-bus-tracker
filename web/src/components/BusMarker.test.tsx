import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import BusMarker from "./BusMarker";

describe("BusMarker", () => {
  it("labels a bus standing at the stop", () => {
    render(
      <BusMarker buses={[{ plateNumb: "EAL-5812" }]} placement="at-stop" stopName="捷運龍山寺站" />,
    );
    expect(screen.getByTestId("bus-at-stop")).toBeInTheDocument();
    expect(screen.getByText("EAL-5812 停靠於 捷運龍山寺站")).toBeInTheDocument();
  });

  it("labels a bus that already left the stop", () => {
    render(
      <BusMarker buses={[{ plateNumb: "KKB-1789" }]} placement="between" stopName="昆明街口" />,
    );
    expect(screen.getByTestId("bus-between")).toBeInTheDocument();
    expect(screen.getByText("KKB-1789 已駛離 昆明街口，前往下一站")).toBeInTheDocument();
  });

  it("lists every plate when several buses share a position", () => {
    render(
      <BusMarker
        buses={[{ plateNumb: "KKA-0161" }, { plateNumb: "KKB-1785" }]}
        placement="between"
        stopName="植物園"
      />,
    );
    expect(
      screen.getByText("KKA-0161, KKB-1785 已駛離 植物園，前往下一站"),
    ).toBeInTheDocument();
  });

  it("renders nothing without buses", () => {
    const { container } = render(<BusMarker buses={[]} placement="at-stop" stopName="植物園" />);
    expect(container).toBeEmptyDOMElement();
  });

  it("renders nothing when buses is undefined", () => {
    const { container } = render(
      <BusMarker buses={undefined} placement="between" stopName="植物園" />,
    );
    expect(container).toBeEmptyDOMElement();
  });
});
