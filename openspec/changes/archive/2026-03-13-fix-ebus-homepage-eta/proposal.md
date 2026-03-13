## Why

eBus 資料源的首頁收藏站點不顯示到站時間（顯示「—」）。原因是 eBus API 的 ETA 回傳不包含 stopId 和 stopName（只有站序 sn），導致前端 etaMap 用 stopId 或 stopName 查不到對應的 ETA 資料。TDX 不受影響，因為 TDX ETA 回傳有完整的 StopID 和 StopName。

## What Changes

- 後端在回傳 eBus ETA 時，用站序（Sequence）比對 GetStops 的結果，補上 stopId 和 stopName
- 確保前端收藏面板和路線詳情頁在 eBus 資料源下都能正確顯示到站狀態

## Capabilities

### New Capabilities

（無）

### Modified Capabilities

（無 spec 層級變更，純 bug fix — eBus convertETAs 未填入 stopId/stopName 是實作遺漏）

## Impact

- `internal/ebus/ebus.go`：GetETA 需要額外呼叫 GetStops 或在 handler 層補上站點資訊
- `internal/handler/routes.go`：GetETA handler 可能需要在 eBus 回傳後補上站點資訊
- 前端不需改動（etaMap 的 stopId/stopName 查找邏輯已正確）
