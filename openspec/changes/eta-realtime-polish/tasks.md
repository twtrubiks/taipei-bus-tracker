## 1. 前端狀態文字與三級分類

- [ ] 1.1 新增 `web/src/utils/etaStatus.ts`：依秒數回傳狀態文字（0–90 進站中、91–179 將到站、≥180 約 N 分、負值對應特殊狀態）
- [ ] 1.2 `web/src/utils/statusColor.ts`：顏色同步改為三級
- [ ] 1.3 `StopList.tsx`、`HomePage.tsx` 改用 `etaStatus(eta.eta)`，不再讀 `eta.status`
- [ ] 1.4 補上 `etaStatus` 與 `statusColor` 的邊界值測試（0、90、91、179、180、181、-1~-4）

## 2. 後端移除 status 欄位

- [ ] 2.1 `internal/model/model.go`：`StopETA` 移除 `Status` 欄位
- [ ] 2.2 `internal/handler/routes.go`：`GetETA` 移除填入 status 的迴圈
- [ ] 2.3 `web/src/api/types.ts`：`StopETA` 移除 `status`
- [ ] 2.4 `useFavoritesEta.ts`：狀態比對改為只比 `eta`
- [ ] 2.5 更新 `internal/handler/routes_test.go`

## 3. CLI 狀態文字三級化

- [ ] 3.1 `internal/model/eta.go`：`ETAStatus` 改為三級
- [ ] 3.2 `cmd/notify/monitor.go`：`formatETA` 改為三級
- [ ] 3.3 更新 `internal/model/eta_test.go` 與 `cmd/notify` 測試，涵蓋相同邊界值

## 4. 本地每秒倒數

- [ ] 4.1 新增 `web/src/hooks/useTick.ts`：每秒觸發重繪，分頁不可見時停止
- [ ] 4.2 新增 `web/src/utils/countdown.ts`：`countdownEta(eta, fetchedAt, now)`，負值原樣回傳、elapsed 上限 60 秒、下限 clamp 為 0
- [ ] 4.3 `useEta.ts`、`useFavoritesEta.ts` 回傳最後一次成功取得資料的 `fetchedAt`
- [ ] 4.4 `StopList.tsx`、`HomePage.tsx` 套用倒數後的秒數計算文字與顏色
- [ ] 4.5 補上 `countdownEta` 與 `useTick` 測試

## 5. 顯示狀態語意化

- [ ] 5.1 資料延遲（超過 60 秒未成功更新）時顯示提示與最後更新時間
- [ ] 5.2 上游回空 stops 時顯示「暫無班次資訊」，與「站點載入失敗」區分
- [ ] 5.3 補上對應測試

## 6. 驗證

- [ ] 6.1 `make lint` 通過
- [ ] 6.2 `make test` 通過（Go + React）
- [ ] 6.3 手動驗證：路線詳情頁秒數每秒遞減、15 秒刷新為真值、切到背景再回來立即更新
- [ ] 6.4 手動驗證：斷開後端後畫面凍結倒數並顯示資料延遲
