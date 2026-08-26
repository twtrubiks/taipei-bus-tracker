export interface Route {
  routeId: string;
  routeName: string;
  startStop: string;
  endStop: string;
  source?: string;
}

export interface Stop {
  stopId: string;
  stopName: string;
  sequence: number;
  source?: string;
}

export interface Bus {
  plateNumb: string;
}

export interface StopETA {
  stopId: string;
  stopName: string;
  sequence: number;
  eta: number;
  /** Buses currently stopped at this stop. */
  buses: Bus[];
  /**
   * Buses that already left this stop and are on their way to the next one.
   * Empty for providers whose upstream does not distinguish the two (TDX).
   */
  departedBuses?: Bus[];
  source: string;
}

export interface ETAResponse {
  route: string;
  direction: number;
  source: string;
  updatedAt: string;
  stops: StopETA[];
}

export interface Favorite {
  routeId: string;
  routeName: string;
  direction: number;
  stopId: string;
  stopName: string;
  sequence: number;
  tdxRouteId?: string;
  ebusRouteId?: string;
  tdxStopId?: string;
  ebusStopId?: string;
}
