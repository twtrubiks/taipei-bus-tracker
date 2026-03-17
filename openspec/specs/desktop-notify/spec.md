## Purpose

桌面通知功能，透過 notify-send 或 kdialog 發送到站提醒。

## Requirements

### Requirement: 桌面通知（notify-send / kdialog fallback）
當 ETA 跌破閾值時，系統 SHALL 透過偵測到的通知工具發送桌面通知。

#### Scenario: 使用 notify-send 發送通知
- **WHEN** 系統偵測到 `notify-send` 可用且 ETA 跳變去重邏輯判定需要通知
- **THEN** 系統執行 `notify-send --urgency=critical` 發送通知，標題包含路線名稱，內容包含站名與 ETA

#### Scenario: 使用 kdialog 發送通知
- **WHEN** 系統未偵測到 `notify-send` 但偵測到 `kdialog` 可用且需要通知
- **THEN** 系統執行 `kdialog --passivepopup` 發送通知，包含路線名稱、站名與 ETA

#### Scenario: 通知指令執行失敗
- **WHEN** 通知指令執行失敗
- **THEN** 系統 log 錯誤，但不中斷監控（Terminal log 仍會顯示通知觸發）

### Requirement: 啟動時偵測通知工具
系統啟動時 SHALL 依序偵測可用的通知工具：`notify-send` → `kdialog` → 無。

#### Scenario: notify-send 存在
- **WHEN** 系統偵測到 `notify-send` 可用
- **THEN** 使用 `notify-send` 作為通知工具，顯示「通知工具: notify-send」

#### Scenario: 僅 kdialog 存在
- **WHEN** 系統未偵測到 `notify-send` 但偵測到 `kdialog` 可用
- **THEN** 使用 `kdialog` 作為通知工具，顯示「通知工具: kdialog (KDE)」

#### Scenario: 兩者皆不存在
- **WHEN** 系統未偵測到 `notify-send` 也未偵測到 `kdialog`
- **THEN** 顯示警告訊息，但仍允許啟動（僅 terminal log 模式）
