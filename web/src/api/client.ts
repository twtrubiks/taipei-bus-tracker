import type { Route, Stop, ETAResponse } from "./types";

const BASE = "/api";

export async function searchRoutes(keyword: string): Promise<Route[]> {
  const res = await fetch(
    `${BASE}/routes/search?q=${encodeURIComponent(keyword)}`,
  );
  if (!res.ok) throw new Error(`search failed: ${res.status}`);
  return res.json() as Promise<Route[]>;
}

export async function getStops(
  routeId: string,
  direction: number,
): Promise<Stop[]> {
  const res = await fetch(`${BASE}/routes/${routeId}/stops?gb=${direction}`);
  if (!res.ok) throw new Error(`getStops failed: ${res.status}`);
  return res.json() as Promise<Stop[]>;
}

export async function getETA(
  routeId: string,
  direction: number,
): Promise<ETAResponse> {
  const res = await fetch(`${BASE}/routes/${routeId}/eta?gb=${direction}`);
  if (!res.ok) throw new Error(`getETA failed: ${res.status}`);
  return res.json() as Promise<ETAResponse>;
}
