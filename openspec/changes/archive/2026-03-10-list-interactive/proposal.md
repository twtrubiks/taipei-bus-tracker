## Summary

將 `bus-notify --list` 從純顯示改為互動式選單。使用者可用數字選擇快捷直接啟動監控，省去手動輸入長名稱的麻煩。

## Motivation

目前啟動已儲存的快捷需要打完整名稱：`bus-notify 看診台大醫院`。名稱長時不方便，且容易打錯。`--list` 已經列出所有快捷，差一步就能直接選擇啟動。

## Proposed Change

- `--list` 輸出加上編號 `[1]`、`[2]`、...
- 列出後新增提示「選擇（Enter 取消）:」
- 使用者輸入數字 → 載入對應快捷，進入監控
- 使用者按 Enter 或無效輸入 → 結束（與現有行為一致）

## Scope

- 僅修改 `cmd/notify/shortcut.go`（`listShortcuts` 函式）
- 僅修改 `cmd/notify/main.go`（`--list` 分支邏輯）
- 不影響 `bus-notify <名稱>`、`--delete`、互動選站等現有功能

## Example

```
$ bus-notify --list
快捷列表：
  [1] 上班     299 捷運忠孝復興站（去程）5 分鐘
  [2] 回家     299 永和（回程）5 分鐘
  [3] 看醫生   紅32 馬偕醫院（去程）3 分鐘

選擇（Enter 取消）: 2
載入快捷「回家」: 299 永和（回程），5 分鐘前通知
→ 進入監控模式...
```
