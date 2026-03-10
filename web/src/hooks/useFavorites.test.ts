import { describe, it, expect, beforeEach } from "vitest";
import { renderHook, act } from "@testing-library/react";
import { useFavorites } from "./useFavorites";
import type { Favorite } from "../api/types";

const fav1: Favorite = {
  routeId: "route1",
  routeName: "299",
  direction: 0,
  stopId: "stop1",
  stopName: "台北車站",
  sequence: 3,
};

const fav2: Favorite = {
  routeId: "route1",
  routeName: "299",
  direction: 1,
  stopId: "stop2",
  stopName: "萬華",
  sequence: 1,
};

beforeEach(() => {
  localStorage.clear();
});

describe("useFavorites", () => {
  it("starts with empty favorites", () => {
    const { result } = renderHook(() => useFavorites());
    expect(result.current.favorites).toEqual([]);
  });

  it("adds a favorite and writes to localStorage", () => {
    const { result } = renderHook(() => useFavorites());

    act(() => {
      result.current.addFavorite(fav1);
    });

    expect(result.current.favorites).toEqual([fav1]);
    expect(JSON.parse(localStorage.getItem("bus-favorites")!)).toEqual([fav1]);
  });

  it("does not add duplicate favorites", () => {
    const { result } = renderHook(() => useFavorites());

    act(() => {
      result.current.addFavorite(fav1);
    });
    act(() => {
      result.current.addFavorite(fav1);
    });

    expect(result.current.favorites).toHaveLength(1);
  });

  it("removes a favorite and updates localStorage", () => {
    const { result } = renderHook(() => useFavorites());

    act(() => {
      result.current.addFavorite(fav1);
      result.current.addFavorite(fav2);
    });

    act(() => {
      result.current.removeFavorite(fav1.routeId, fav1.direction, fav1.stopId);
    });

    expect(result.current.favorites).toEqual([fav2]);
    expect(JSON.parse(localStorage.getItem("bus-favorites")!)).toEqual([fav2]);
  });

  it("restores favorites from localStorage on mount", () => {
    localStorage.setItem("bus-favorites", JSON.stringify([fav1, fav2]));

    const { result } = renderHook(() => useFavorites());

    expect(result.current.favorites).toEqual([fav1, fav2]);
  });

  it("isFavorite returns correct boolean", () => {
    const { result } = renderHook(() => useFavorites());

    act(() => {
      result.current.addFavorite(fav1);
    });

    expect(
      result.current.isFavorite(fav1.routeId, fav1.direction, fav1.stopId),
    ).toBe(true);
    expect(
      result.current.isFavorite(fav2.routeId, fav2.direction, fav2.stopId),
    ).toBe(false);
  });

  it("handles corrupted localStorage gracefully", () => {
    localStorage.setItem("bus-favorites", "not-json");

    const { result } = renderHook(() => useFavorites());

    expect(result.current.favorites).toEqual([]);
  });
});
