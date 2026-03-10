## Context

專案已有完整的公車資料抽象層（`BusDataSource` interface），包含 TDX 和 eBus 兩個 provider、fallback 機制、記憶體快取。目前這些模組被 `cmd/server/main.go` 的 HTTP server 使用，搭配 React 前端。

使用者在 KDE 桌面環境，希望用最輕量的方式監控公車到站：一個 terminal 持續跑的 CLI 工具。

## Goals / Non-Goals

**Goals:**
- 純 Go CLI 工具，無需 web server 或前端
- 互動式引導使用者選擇路線、方向、站點
- 持續 polling ETA，根據距離動態調整間隔
- 透過 `notify-send` 發送桌面通知，同一班車只通知一次
- Terminal 逐行 log 輸出，帶時間戳

**Non-Goals:**
- 原地刷新 TUI（Phase 2 再做）
- 同時監控多個站點
- 設定檔記住常用站點
- 聲音通知
- 跨平台通知（目前只支援 freedesktop / KDE）

## Decisions

### 1. 入口點：新增 `cmd/notify/main.go`

與現有 `cmd/server/` 平行，共用 `internal/` 模組。不需要修改任何既有程式碼。

**替代方案**：在 server 加 WebSocket push → 太重，違反簡化初衷。

### 2. 直接使用 FallbackService

複用 `internal/handler/fallback.go` 的 `FallbackService`，自動享有 TDX → eBus fallback 和快取。

**替代方案**：直接呼叫單一 provider → 失去 fallback 可靠性。

### 3. 通知去重：ETA 跳變偵測（方案 B）

用一個 `wasAboveThreshold bool` 狀態：
- 初始值 `true`
- ETA 從 > threshold 跌到 <= threshold → 通知，設 `false`
- ETA 回到 > threshold 或狀態變為進站中/未發車/末班駛離 → 重設 `true`

**替代方案**：用車牌去重 → eBus 不回傳車牌，TDX 有時也沒有，不可靠。

### 4. 動態 Polling 間隔

| 條件 | 間隔 |
|------|------|
| ETA > threshold × 2 | 60 秒 |
| threshold < ETA <= threshold × 2 | 30 秒 |
| ETA <= threshold | 15 秒 |
| API 錯誤 | 維持上次間隔，log 錯誤繼續 |

避免 TDX rate limit（20 次/分鐘），最快也只有 15 秒一次。

### 5. 通知方式：notify-send → kdialog fallback

啟動時依序偵測可用的通知工具：
1. `notify-send`（libnotify，freedesktop 標準）
2. `kdialog --passivepopup`（KDE 原生，不需額外安裝）
3. 都沒有 → 僅 terminal 顯示模式

**替代方案**：Go D-Bus library 直接呼叫 → 增加依賴，shell out 已經夠好。

### 6. 優雅停止：signal handling

監聽 `SIGINT`（Ctrl+C）和 `SIGTERM`，收到後印出摘要並退出。用 `context.WithCancel` 傳播取消信號。

## Risks / Trade-offs

- **[TDX rate limit]** → 動態間隔最快 15 秒，不會超過 4 次/分鐘
- **[eBus 無車牌]** → 用跳變去重而非車牌，極端情況下兩班車 ETA 都 < threshold 可能少通知一次，可接受
- **[通知工具不存在]** → 啟動時偵測 notify-send → kdialog fallback，都沒有則僅 terminal 模式
- **[API 全部失敗]** → log 錯誤繼續 polling，不中斷監控
