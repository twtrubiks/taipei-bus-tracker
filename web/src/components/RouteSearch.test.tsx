import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter } from "react-router-dom";
import RouteSearch from "./RouteSearch";
import * as client from "../api/client";

const mockNavigate = vi.fn();
vi.mock("react-router-dom", async () => {
  const actual = await vi.importActual("react-router-dom");
  return { ...actual, useNavigate: () => mockNavigate };
});

vi.mock("../api/client");

const mockRoutes = [
  {
    routeId: "0100000100",
    routeName: "1",
    startStop: "萬華",
    endStop: "松仁路",
  },
  {
    routeId: "0100001200",
    routeName: "12",
    startStop: "松德",
    endStop: "圓環",
  },
];

beforeEach(() => {
  vi.clearAllMocks();
});

function renderSearch() {
  return render(
    <MemoryRouter>
      <RouteSearch />
    </MemoryRouter>,
  );
}

describe("RouteSearch", () => {
  it("shows matching results when keyword is entered", async () => {
    vi.mocked(client.searchRoutes).mockResolvedValue(mockRoutes);
    const user = userEvent.setup();
    renderSearch();

    await user.type(screen.getByRole("textbox"), "1");

    await waitFor(() => {
      expect(screen.getByText("1")).toBeInTheDocument();
      expect(screen.getByText("12")).toBeInTheDocument();
    });
  });

  it("clears results when input is empty", async () => {
    vi.mocked(client.searchRoutes).mockResolvedValue(mockRoutes);
    const user = userEvent.setup();
    renderSearch();

    await user.type(screen.getByRole("textbox"), "1");
    await waitFor(() => expect(screen.getByText("1")).toBeInTheDocument());

    await user.clear(screen.getByRole("textbox"));
    await waitFor(() =>
      expect(screen.queryByRole("list")).not.toBeInTheDocument(),
    );
  });

  it("navigates to route page on result click", async () => {
    vi.mocked(client.searchRoutes).mockResolvedValue([mockRoutes[0]]);
    const user = userEvent.setup();
    renderSearch();

    await user.type(screen.getByRole("textbox"), "1");
    await waitFor(() => expect(screen.getByText("1")).toBeInTheDocument());

    await user.click(screen.getByText("1"));
    expect(mockNavigate).toHaveBeenCalledWith("/route/0100000100?name=1");
  });

  it("shows no results message when search returns empty", async () => {
    vi.mocked(client.searchRoutes).mockResolvedValue([]);
    const user = userEvent.setup();
    renderSearch();

    await user.type(screen.getByRole("textbox"), "zzz");
    await waitFor(() =>
      expect(screen.getByText("找不到符合的路線")).toBeInTheDocument(),
    );
  });
});
