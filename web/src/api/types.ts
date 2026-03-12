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
  status: string;
  buses: Bus[];
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
