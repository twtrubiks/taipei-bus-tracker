## MODIFIED Requirements

### Requirement: 逐行 Log 輸出
監控模式中，系統 SHALL 每次 polling 後輸出一行包含時間戳和 ETA 資訊的 log。ETA 文字 SHALL 依秒數分三級呈現，與 Web 前端一致。

#### Scenario: 正常 ETA 輸出
- **WHEN** 取得 ETA 為 480 秒（8 分鐘）
- **THEN** 輸出格式如 `14:30:01  ETA 8 分`

#### Scenario: 通知觸發時的輸出
- **WHEN** ETA 跌破閾值且觸發通知
- **THEN** 輸出格式如 `14:33:31  ETA 5 分  🔔 已通知！`

#### Scenario: 進站中
- **WHEN** ETA 介於 0 至 90 秒
- **THEN** 輸出 `進站中`

#### Scenario: 將到站
- **WHEN** ETA 介於 91 至 179 秒
- **THEN** 輸出 `將到站`

#### Scenario: 特殊狀態輸出
- **WHEN** ETA 為負數狀態碼
- **THEN** 輸出對應的狀態文字（`未發車`、`末班駛離`、`交管不停靠`、`未營運`）
