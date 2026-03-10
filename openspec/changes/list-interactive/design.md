## Context

`bus-notify --list` 目前只顯示快捷列表後結束。使用者想啟動快捷仍需另外打 `bus-notify <名稱>`。改為互動式選單讓使用者一步完成。

## Goals / Non-Goals

**Goals:**
- `--list` 顯示帶編號的快捷列表
- 使用者可輸入數字選擇，直接進入監控
- Enter 或無效輸入則結束程式（向後相容）

**Non-Goals:**
- 不改變 `bus-notify`（無參數）的互動選站流程
- 不改變 `bus-notify <名稱>` 的快捷啟動流程

## Decisions

### 1. `listShortcuts` 回傳 `*Shortcut`

將 `listShortcuts()` 改為回傳 `*Shortcut`（選中）或 `nil`（取消/無快捷）。`main.go` 根據回傳值決定是否繼續進入監控。

**替代方案**：新增獨立函式 `selectShortcutFromList()` → 多一個函式但語意更清楚。不過 `listShortcuts` 本身就是 `--list` 的處理函式，直接擴充較自然。

### 2. `--list` 選中後 fall through 到資料源初始化

目前 `--list` 在 config/TDX/eBus 初始化之前就 return。選中快捷後不能 return，需要繼續初始化資料源再進入 `runMonitor`。做法是把 `--list` 選中的快捷填入與 positional arg 相同的變數，統一走後續流程。

### 3. 無快捷時不顯示選擇提示

如果 `notify.yaml` 無快捷，只顯示「尚無快捷設定」後結束，不顯示選擇提示。

## Risks / Trade-offs

- **[行為改變]** `--list` 原本是非互動的，腳本若 pipe `--list` 可能卡在等待輸入 → 此工具為個人使用，不影響。若需要可檢測 stdin 是否為 terminal
