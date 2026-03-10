import { useState, useEffect } from "react";
import { useParams, Link } from "react-router-dom";
import { getStops } from "../api/client";
import type { Stop } from "../api/types";
import DirectionSelector from "../components/DirectionSelector";
import StopList from "../components/StopList";
import { useEta } from "../hooks/useEta";

export default function RoutePage() {
  const { routeId } = useParams<{ routeId: string }>();
  const [direction, setDirection] = useState(0);
  const [stops, setStops] = useState<Stop[]>([]);
  const [loadedKey, setLoadedKey] = useState("");

  const stopsKey = `${routeId}-${direction}`;
  const loadingStops = stopsKey !== loadedKey;

  const { data: etaData, error: etaError } = useEta(routeId ?? "", direction);

  useEffect(() => {
    if (!routeId) return;
    getStops(routeId, direction)
      .then((data) => {
        setStops(data);
        setLoadedKey(`${routeId}-${direction}`);
      })
      .catch(() => {
        setStops([]);
        setLoadedKey(`${routeId}-${direction}`);
      });
  }, [routeId, direction]);

  if (!routeId) return null;

  return (
    <div className="mx-auto max-w-lg p-4">
      <div className="mb-4 flex items-center gap-2">
        <Link to="/search" className="text-blue-600 hover:text-blue-800">
          &larr; 搜尋
        </Link>
        <h1 className="text-xl font-bold">路線 {routeId}</h1>
      </div>

      <DirectionSelector direction={direction} onChange={setDirection} />

      {etaError && (
        <p className="mt-2 text-sm text-red-500">載入失敗，稍後重試</p>
      )}

      {loadingStops ? (
        <p className="mt-4 text-gray-500">載入站點中...</p>
      ) : (
        <StopList stops={stops} etas={etaData?.stops ?? []} />
      )}
    </div>
  );
}
