## 1. 設定檔基礎

- [x] 1.1 定義 `Shortcut` struct 和 `NotifyConfig` struct（含 YAML tag），新增 `notify.yaml` 的讀寫函式（loadNotifyConfig / saveNotifyConfig）
- [x] 1.2 將 `notify.yaml` 加入 `.gitignore`

## 2. CLI 參數解析

- [x] 2.1 實作 CLI flag 解析：`--list`、`--delete <名稱>`、positional arg 作為快捷名稱
- [x] 2.2 實作 `--list` 功能：讀取 notify.yaml，格式化輸出所有快捷
- [x] 2.3 實作 `--delete` 功能：從 notify.yaml 移除指定快捷，處理不存在的情況

## 3. 快捷儲存

- [ ] 3.1 互動選站完成後，提示使用者輸入快捷名稱（Enter 跳過）
- [ ] 3.2 同名快捷覆蓋確認：顯示現有設定，詢問 y/N

## 4. 快捷啟動

- [ ] 4.1 實作快捷載入：讀取 notify.yaml 找到對應快捷，建構 Route/Stop 物件，直接呼叫 runMonitor
- [ ] 4.2 處理錯誤情況：快捷不存在時列出可用快捷、notify.yaml 不存在或格式錯誤時顯示錯誤
