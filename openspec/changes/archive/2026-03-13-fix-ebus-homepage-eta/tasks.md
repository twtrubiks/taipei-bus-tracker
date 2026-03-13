## 1. 後端修復

- [x] 1.1 `internal/handler/routes.go`：GetETA handler 偵測 eBus source 時，用 GetStops 取得站點列表，以 Sequence 匹配補上 stopId 和 stopName
- [x] 1.2 確保 GetStops 結果走 cache（避免每次 ETA 都打一次 GetStops API）

## 2. 測試

- [x] 2.1 補寫或更新 handler 測試：驗證 eBus ETA 回傳包含正確的 stopId 和 stopName
- [x] 2.2 `go test ./...` 全部通過

## 3. 驗證

- [x] 3.1 手動測試：eBus 資料源下，首頁收藏站點顯示正確的到站時間
- [x] 3.2 手動測試：路線詳情頁在 eBus 資料源下仍正常
