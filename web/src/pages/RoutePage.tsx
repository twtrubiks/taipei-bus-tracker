import { useState, useEffect } from "react";
import { useParams, useSearchParams, Link } from "react-router-dom";
import { getStops } from "../api/client";
import type { Stop } from "../api/types";
import DirectionSelector from "../components/DirectionSelector";
import StopList from "../components/StopList";
import { useEta } from "../hooks/useEta";
import { useFavorites } from "../hooks/useFavorites";
import { useNotification } from "../hooks/useNotification";

export default function RoutePage() {
  const { routeId } = useParams<{ routeId: string }>();
  const [searchParams] = useSearchParams();
  const routeName = searchParams.get("name") ?? routeId ?? "";
  const [direction, setDirection] = useState(0);
  const [stops, setStops] = useState<Stop[]>([]);
  const [loadedKey, setLoadedKey] = useState("");

  const stopsKey = `${routeId}-${direction}`;
  const loadingStops = stopsKey !== loadedKey;

  const { data: etaData, error: etaError } = useEta(routeId ?? "", direction);
  const { addFavorite, removeFavorite, isFavorite } = useFavorites();
  const { addAlert, removeAlert, getAlert, checkAlerts, permissionDenied } =
    useNotification();

  // Check notification alerts whenever ETA data updates
  useEffect(() => {
    if (routeId && etaData?.stops) {
      checkAlerts(routeId, direction, etaData.stops);
    }
  }, [routeId, direction, etaData, checkAlerts]);

  useEffect(() => {
    if (!routeId) return;
    let cancelled = false;
    getStops(routeId, direction)
      .then((data) => {
        if (!cancelled) setStops(data);
      })
      .catch(() => {
        if (!cancelled) setStops([]);
      })
      .finally(() => {
        if (!cancelled) setLoadedKey(`${routeId}-${direction}`);
      });
    return () => {
      cancelled = true;
    };
  }, [routeId, direction]);

  if (!routeId) return null;

  const handleToggleFavorite = (stop: Stop) => {
    if (isFavorite(routeId, direction, stop.stopId)) {
      removeFavorite(routeId, direction, stop.stopId);
    } else {
      addFavorite({
        routeId,
        routeName,
        direction,
        stopId: stop.stopId,
        stopName: stop.stopName,
        sequence: stop.sequence,
      });
    }
  };

  const handleSetAlert = (stop: Stop, minutes: number) => {
    addAlert({
      routeId,
      routeName,
      direction,
      stopId: stop.stopId,
      stopName: stop.stopName,
      thresholdMinutes: minutes,
    });
  };

  const handleRemoveAlert = (stop: Stop) => {
    removeAlert(routeId, direction, stop.stopId);
  };

  return (
    <div className="mx-auto max-w-lg p-4">
      <div className="mb-4 flex items-center gap-2">
        <Link to="/search" className="text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300">
          &larr; 搜尋
        </Link>
        <h1 className="text-xl font-bold">路線 {routeName}</h1>
      </div>

      <DirectionSelector direction={direction} onChange={setDirection} />

      {permissionDenied && (
        <p className="mt-2 text-sm text-amber-600">
          通知權限被拒絕，請在瀏覽器設定中允許通知
        </p>
      )}

      {etaData?.source && (
        <p className="mt-2 text-xs text-gray-400">
          資料來源：{etaData.source === "tdx" ? "TDX" : "eBus"}
        </p>
      )}

      {etaError && (
        <p className="mt-2 text-sm text-red-500">載入失敗，稍後重試</p>
      )}

      {loadingStops ? (
        <p className="mt-4 text-gray-500">載入站點中...</p>
      ) : (
        <StopList
          stops={stops}
          etas={etaData?.stops ?? []}
          routeId={routeId}
          direction={direction}
          isFavorite={isFavorite}
          onToggleFavorite={handleToggleFavorite}
          getAlert={
            routeId
              ? (stopId: string) => getAlert(routeId, direction, stopId)
              : undefined
          }
          onSetAlert={handleSetAlert}
          onRemoveAlert={handleRemoveAlert}
        />
      )}
    </div>
  );
}
