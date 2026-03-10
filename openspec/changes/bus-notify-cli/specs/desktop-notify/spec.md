## ADDED Requirements

### Requirement: notify-send 桌面通知
當 ETA 跌破閾值時，系統 SHALL 透過 `notify-send` 發送桌面通知。

#### Scenario: 發送通知
- **WHEN** ETA 跳變去重邏輯判定需要通知
- **THEN** 系統執行 `notify-send --urgency=critical` 發送通知，標題包含路線名稱，內容包含站名與 ETA

#### Scenario: notify-send 執行失敗
- **WHEN** `notify-send` 指令執行失敗
- **THEN** 系統 log 錯誤，但不中斷監控（Terminal log 仍會顯示通知觸發）

### Requirement: 啟動時檢查 notify-send
系統啟動時 SHALL 檢查 `notify-send` 指令是否存在。

#### Scenario: notify-send 存在
- **WHEN** 系統偵測到 `notify-send` 可用
- **THEN** 正常啟動

#### Scenario: notify-send 不存在
- **WHEN** 系統偵測到 `notify-send` 不可用
- **THEN** 顯示警告訊息（提示安裝 libnotify-tools），但仍允許啟動（僅 terminal log 模式）
