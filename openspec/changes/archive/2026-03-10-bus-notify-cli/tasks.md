## 1. 專案設定

- [x] 1.1 建立 `cmd/notify/main.go` 入口，初始化 config、TDX/eBus provider、FallbackService
- [x] 1.2 啟動時偵測通知工具（notify-send → kdialog fallback），不存在則顯示警告

## 2. 互動式選站

- [x] 2.1 實作路線搜尋互動：提示輸入關鍵字 → 呼叫 SearchRoutes → 列出選項（單一結果自動選擇）
- [x] 2.2 實作方向選擇互動：顯示去程/回程（含起終站名稱）
- [x] 2.3 實作站點選擇互動：呼叫 GetStops → 列出所有站點供選擇
- [x] 2.4 實作閾值輸入互動：提示輸入分鐘數，預設 5，驗證輸入合法性

## 3. ETA 監控迴圈

- [x] 3.2a 撰寫 `calcInterval` 純函式的 table-driven 測試（TDD：先寫測試）
- [x] 3.2b 實作 `calcInterval(eta, threshold) → time.Duration` 純函式，通過測試
- [x] 3.3a 撰寫 ETA 跳變去重狀態機的 table-driven 測試（TDD：先寫測試），覆蓋：正常到站、兩班車依序、不重複通知、ETA 波動不誤觸、未發車/末班駛離重設
- [x] 3.3b 實作跳變去重邏輯（`wasAboveThreshold` 狀態追蹤），通過測試
- [x] 3.1 實作 polling 迴圈：用 `time.Ticker` + `context.WithCancel`，定期呼叫 GetETA 取得目標站點 ETA，整合 3.2b + 3.3b
- [x] 3.4 實作逐行 log 輸出：每次 poll 輸出時間戳 + ETA + 通知標記

## 4. 桌面通知

- [x] 4.1 實作通知發送：根據 `detectNotifyTool` 結果，使用 `notify-send`（urgency=critical）或 `kdialog`（passivepopup）發送通知，含路線名和站名
- [x] 4.2 處理通知指令執行失敗：log 錯誤但不中斷監控

## 5. 生命週期管理

- [x] 5.1 實作 signal handling：監聽 SIGINT/SIGTERM，觸發 context cancel
- [x] 5.2 優雅停止時輸出摘要：總監控時間、通知次數
