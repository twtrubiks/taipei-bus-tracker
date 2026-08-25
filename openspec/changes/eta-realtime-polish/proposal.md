## Why

到站資訊的呈現目前有三個體感問題：

1. **畫面每 15 秒才動一次**：前端輪詢間隔 15 秒，中間畫面完全靜止，使用者無法判斷資料是新的還是卡住了。
2. **「進站中」涵蓋範圍過寬**：`ETAStatus` 把 0~180 秒全部顯示為「進站中」，但 3 分鐘足以讓人白站在站牌前等，語意不精準。
3. **拿不到資料與沒有班次混在一起**：fetch 失敗時只有一行紅字，站點狀態仍顯示上一輪的舊值卻沒有任何標示；上游回空陣列時也只顯示「無站點資料」，使用者分不出是後端掛了還是這條路線沒車。

同時，狀態文字目前由後端 `handler.GetETA` 算好塞進 `status` 欄位。一旦前端要做每秒倒數，秒數會逐秒改變而後端算好的字串不會，兩者必然不一致。

## What Changes

- **狀態文字下放前端**：API 移除 `status` 欄位，只回傳 `eta` 秒數；前端依秒數自行計算顯示文字，成為單一真值來源。
- **到站狀態細分三級**：0–90 秒「進站中」、91–179 秒「將到站」、≥180 秒「約 N 分」，顏色同步分級。
- **本地每秒倒數（Tick）**：記錄每次成功取得資料的時間戳，每秒僅觸發重繪、由秒數與時間戳推導顯示值，不增加任何 API 請求。資料超過 60 秒未更新時凍結倒數並標示資料延遲。
- **CLI 狀態文字同步三級**：`model.ETAStatus` 與 `cmd/notify` 的 `formatETA` 一併調整，讓 Web 與 CLI 對同一筆資料講一樣的話。
- **顯示狀態語意化**：區分「載入中」「連線異常（保留舊資料並標示更新時間）」「無班次資訊」三種情況。

## Capabilities

### New Capabilities

（無新 capability）

### Modified Capabilities

- **go-api-proxy**：ETA endpoint 回應格式移除 `status` 欄位。
- **pwa-frontend**：到站狀態顯示改為三級、新增每秒本地倒數、錯誤與空資料狀態語意化。
- **eta-monitor**：CLI log 的 ETA 文字改為三級狀態。

## Impact

- `internal/model/model.go`：`StopETA` 移除 `Status` 欄位
- `internal/model/eta.go`：`ETAStatus` 改為三級（仍供 CLI 與 dev tool 使用）
- `internal/handler/routes.go`：`GetETA` 移除填入 status 的迴圈
- `cmd/notify/monitor.go`：`formatETA` 改為三級
- `web/src/utils/etaStatus.ts`（新增）、`statusColor.ts`、`countdown.ts`（新增）
- `web/src/hooks/useTick.ts`（新增）、`useEta.ts`、`useFavoritesEta.ts`
- `web/src/components/StopList.tsx`、`web/src/pages/HomePage.tsx`、`RoutePage.tsx`
- `web/src/api/types.ts`：`StopETA` 移除 `status`
