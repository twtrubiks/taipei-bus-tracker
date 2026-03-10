## Context

`bus-notify` CLI 已完成互動式選站和 ETA 監控功能。使用者希望能把常用的站點組合儲存起來，下次一行指令直接啟動。

## Goals / Non-Goals

**Goals:**
- 互動選站後可選擇儲存為命名快捷
- 指定快捷名稱直接啟動監控，跳過互動流程
- 列出、刪除已存快捷
- 同名覆蓋前提醒確認

**Non-Goals:**
- 同時監控多個快捷
- 快捷排程（cron 整合）
- 快捷匯出/匯入

## Decisions

### 1. 設定檔位置：專案根目錄 `notify.yaml`

與 `config.yaml` 同層，管理直覺。

**替代方案**：`~/.config/bus-notify/notify.yaml`（XDG 標準） → 較正式但對此專案過重，使用者希望集中管理。

### 2. 設定檔格式

```yaml
shortcuts:
  - name: "上班"
    route_id: "Taipei-299"
    route_name: "299"
    start_stop: "台北車站"
    end_stop: "永和"
    direction: 0
    stop_id: "TPE-15234"
    stop_name: "捷運忠孝復興站"
    stop_sequence: 12
    threshold: 5
```

儲存完整 ID（route_id、stop_id、stop_sequence）確保精確匹配，同時存人類可讀欄位（route_name、stop_name）方便 `--list` 顯示和手動編輯。

### 3. 快捷啟動時不驗證 ID

直接用儲存的 ID 呼叫 GetETA，不額外打 GetStops 驗證。公車站 ID 極少變動，且 GetETA 失敗已有錯誤處理（log 錯誤繼續 polling）。

**替代方案**：啟動時 GetStops 驗證 → 多一次 API call，啟動慢幾秒，價值不大。

### 4. CLI 參數解析

```
bus-notify              → 互動式（現有流程）
bus-notify 上班         → 載入快捷，直接監控
bus-notify --list       → 列出所有快捷
bus-notify --delete 上班 → 刪除指定快捷
```

用 Go 標準 `flag` 套件，不引入第三方 CLI 框架。positional arg 作為快捷名稱。

### 5. 儲存時機：互動完成後、監控開始前

互動選站完成後提示「儲存為快捷？輸入名稱（Enter 跳過）」。同名時顯示現有設定並要求確認覆蓋。

### 6. notify.yaml 加入 .gitignore

快捷設定包含個人偏好，不應進版控。

## Risks / Trade-offs

- **[ID 過期]** → 公車站 ID 極少變動，過期時 ETA 顯示「未發車」，使用者可重新互動選站覆蓋快捷
- **[檔案損壞]** → YAML 解析失敗時 log 錯誤，fallback 到互動模式
- **[notify.yaml 不存在]** → 首次使用或 --list/--delete 時提示「尚無快捷」
