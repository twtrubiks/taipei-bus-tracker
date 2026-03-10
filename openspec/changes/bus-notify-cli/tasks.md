## 1. 專案設定

- [ ] 1.1 建立 `cmd/notify/main.go` 入口，初始化 config、TDX/eBus provider、FallbackService
- [ ] 1.2 啟動時檢查 `notify-send` 是否可用，不存在則顯示警告

## 2. 互動式選站

- [ ] 2.1 實作路線搜尋互動：提示輸入關鍵字 → 呼叫 SearchRoutes → 列出選項（單一結果自動選擇）
- [ ] 2.2 實作方向選擇互動：顯示去程/回程（含起終站名稱）
- [ ] 2.3 實作站點選擇互動：呼叫 GetStops → 列出所有站點供選擇
- [ ] 2.4 實作閾值輸入互動：提示輸入分鐘數，預設 5，驗證輸入合法性

## 3. ETA 監控迴圈

- [ ] 3.1 實作 polling 迴圈：用 `time.Ticker` + `context.WithCancel`，定期呼叫 GetETA 取得目標站點 ETA
- [ ] 3.2 實作動態 polling 間隔：根據 ETA 與 threshold 距離調整（60s / 30s / 15s）
- [ ] 3.3 實作 ETA 跳變去重邏輯：`wasAboveThreshold` 狀態追蹤，含狀態重設條件
- [ ] 3.4 實作逐行 log 輸出：每次 poll 輸出時間戳 + ETA + 通知標記

## 4. 桌面通知

- [ ] 4.1 實作 `notify-send` 呼叫：`exec.Command` 發送 urgency=critical 通知，含路線名和站名
- [ ] 4.2 處理 notify-send 執行失敗：log 錯誤但不中斷監控

## 5. 生命週期管理

- [ ] 5.1 實作 signal handling：監聽 SIGINT/SIGTERM，觸發 context cancel
- [ ] 5.2 優雅停止時輸出摘要：總監控時間、通知次數
