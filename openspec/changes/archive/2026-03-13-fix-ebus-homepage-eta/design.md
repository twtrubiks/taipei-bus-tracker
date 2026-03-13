## Context

eBus GetStopDyns API 回傳的 JSON 結構：
```json
[{"sn": 0, "eta": 5, "bi": [...], "bo": [...]}, ...]
```
只有站序 `sn`（0-based），沒有 stopId 和 stopName。目前 `convertETAs` 直接設為空字串。

前端 etaMap 用 `${routeId}:${direction}:${stopId}` 和 `${routeId}:${direction}:name:${stopName}` 做 lookup，兩者都是空就查不到。

路線詳情頁不受影響，因為 `StopList.tsx` 已實作 StopID 優先、Sequence 次之的雙重匹配策略。但首頁收藏面板只用 stopId/stopName 查找，沒有 sequence fallback。

## Goals / Non-Goals

**Goals:**
- eBus ETA 回傳的每個站都有正確的 stopId 和 stopName
- 首頁收藏和路線詳情頁都能正常顯示

**Non-Goals:**
- 不改前端 etaMap 的 lookup 邏輯
- 不改 TDX 的流程（TDX 已正確回傳 stopId/stopName）

## Decisions

### 在 handler 層補上站點資訊

**選擇**: 在 `GetETA` handler 中，當 source 為 eBus 時，額外呼叫 GetStops 取得站點列表，用 Sequence 匹配補上 stopId 和 stopName。

**替代方案**:
- 在 eBus GetETA 內部呼叫 GetStops：耦合 GetETA 和 GetStops 兩個 API 呼叫
- 在前端補 sequence fallback：首頁收藏面板的資料結構不方便做 sequence 匹配

**理由**: handler 層是最自然的整合點，已經有 provider 資訊，且 GetStops 結果可透過 cache 加速。

## Risks / Trade-offs

- eBus ETA 每次請求多一次 GetStops 呼叫 → cache 可緩解，GetStops 結果較穩定
