import { describe, it, expect, vi, afterEach } from "vitest";
import { render, screen, act } from "@testing-library/react";
import DataFreshness from "./DataFreshness";
import { STALE_AFTER_MS } from "../utils/countdown";

afterEach(() => {
  vi.useRealTimers();
});

describe("DataFreshness", () => {
  it("stays silent while the data is fresh", () => {
    vi.useFakeTimers();
    const { container } = render(
      <DataFreshness fetchedAt={Date.now()} failing={false} />,
    );
    expect(container).toBeEmptyDOMElement();
  });

  it("stays silent before anything has been requested", () => {
    const { container } = render(<DataFreshness fetchedAt={0} failing={false} />);
    expect(container).toBeEmptyDOMElement();
  });

  it("reports a failed first load", () => {
    render(<DataFreshness fetchedAt={0} failing />);
    expect(screen.getByRole("status")).toHaveTextContent("即時到站載入失敗");
  });

  it("says the shown data is from before when a refresh fails", () => {
    vi.useFakeTimers();
    render(<DataFreshness fetchedAt={Date.now()} failing />);
    expect(screen.getByRole("status")).toHaveTextContent("更新失敗");
  });

  it("warns that the countdown is frozen once the data goes stale", () => {
    vi.useFakeTimers();
    render(<DataFreshness fetchedAt={Date.now()} failing={false} />);

    act(() => {
      vi.advanceTimersByTime(STALE_AFTER_MS);
    });
    expect(screen.getByRole("status")).toHaveTextContent("資料已延遲");
  });

  it("clears the warning once a refresh succeeds", () => {
    vi.useFakeTimers();
    const { rerender, container } = render(
      <DataFreshness fetchedAt={Date.now()} failing />,
    );
    expect(screen.getByRole("status")).toBeInTheDocument();

    rerender(<DataFreshness fetchedAt={Date.now()} failing={false} />);
    expect(container).toBeEmptyDOMElement();
  });
});
