export interface Route {
  routeId: string;
  routeName: string;
  startStop: string;
  endStop: string;
}

export interface Stop {
  stopId: string;
  stopName: string;
  sequence: number;
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
}
