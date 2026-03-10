## Why

目前 `bus-notify` CLI 每次啟動都需要經過完整的互動式選站流程（搜尋路線 → 選方向 → 選站 → 設閾值），對於每天搭同一班車的使用者太繁瑣。需要「快捷」功能，讓使用者互動選站一次後自動儲存，之後指定名稱即可直接啟動監控。

## What Changes

- 新增 `notify.yaml` 設定檔（專案根目錄），儲存快捷設定
- 互動式選站流程結束後，詢問是否儲存為快捷（可跳過）
- 支援 `bus-notify <快捷名稱>` 直接啟動監控
- 支援 `bus-notify --list` 列出所有快捷
- 支援 `bus-notify --delete <名稱>` 刪除快捷
- 同名快捷覆蓋前提醒使用者確認

## Capabilities

### New Capabilities
- `notify-shortcut`: 快捷設定的儲存、載入、列出、刪除功能

### Modified Capabilities
- `interactive-selection`: 互動式選站流程結束後新增「儲存為快捷」步驟

## Impact

- **新增檔案**：`notify.yaml`（專案根目錄，需加入 `.gitignore`）
- **修改檔案**：`cmd/notify/main.go`（CLI 參數解析、快捷載入/儲存邏輯）
- **無外部依賴變更**
- **無破壞性變更**：無參數時行為完全不變
