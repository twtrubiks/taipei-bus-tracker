## Why

現有架構需要啟動 Go HTTP server + React 前端才能監控公車到站通知，對於「訂閱一個站點、快到就提醒我」這個核心需求過於笨重。需要一個純 Go CLI 工具，在 terminal 持續運行，直接透過 KDE 桌面通知（`notify-send`）提醒使用者。

## What Changes

- 新增 `cmd/notify/main.go` — CLI 入口點，互動式選站 + 監控迴圈
- 新增互動式選站流程：搜尋路線 → 選方向 → 選站點 → 設定閾值
- 新增 ETA polling 迴圈，根據距離閾值動態調整 polling 間隔
- 新增 `notify-send` 桌面通知整合（freedesktop 標準，KDE 原生支援）
- 新增 ETA 跳變去重邏輯，同一班車只通知一次，下一班車自動重新觸發
- Terminal 逐行 log 輸出（帶時間戳）

## Capabilities

### New Capabilities
- `interactive-selection`: 互動式選站流程（路線搜尋、方向選擇、站點選擇、閾值設定）
- `eta-monitor`: ETA polling 監控迴圈，含動態間隔調整與跳變去重邏輯
- `desktop-notify`: 透過 notify-send 發送 KDE 桌面通知

### Modified Capabilities

（無既有 spec 需修改）

## Impact

- **新增程式碼**：`cmd/notify/` 目錄
- **複用既有模組**：`internal/model`、`internal/tdx`、`internal/ebus`、`internal/config`、`internal/cache`
- **外部依賴**：`notify-send` 指令（KDE 預裝 libnotify-tools）
- **無破壞性變更**：不影響現有 web server 或前端
