## Purpose

ETA 監控模式，持續 polling 到站時間並在跌破閾值時觸發通知。

## Requirements

### Requirement: 持續 Polling ETA
進入監控模式後，系統 SHALL 定期呼叫 GetETA 取得目標站點的即時到站資訊，並持續運行直到使用者按 Ctrl+C。

#### Scenario: 正常 polling
- **WHEN** 監控模式啟動
- **THEN** 系統依照動態間隔持續呼叫 GetETA，每次結果輸出一行 log（含時間戳與 ETA）

#### Scenario: API 錯誤
- **WHEN** GetETA 呼叫失敗
- **THEN** 系統 log 錯誤訊息，維持上次 polling 間隔繼續嘗試，不中斷監控

### Requirement: 動態 Polling 間隔
系統 SHALL 根據目前 ETA 與閾值的距離動態調整 polling 間隔，離閾值越近越頻繁。

#### Scenario: ETA 遠離閾值
- **WHEN** ETA > threshold × 2
- **THEN** polling 間隔為 60 秒

#### Scenario: ETA 接近閾值
- **WHEN** threshold < ETA <= threshold × 2
- **THEN** polling 間隔為 30 秒

#### Scenario: ETA 已在閾值內
- **WHEN** ETA <= threshold
- **THEN** polling 間隔為 15 秒

### Requirement: ETA 跳變去重
系統 SHALL 使用跳變偵測邏輯，確保同一班車只觸發一次通知，下一班車進入閾值時重新觸發。

#### Scenario: 首次跌破閾值
- **WHEN** ETA 從 > threshold 變為 <= threshold
- **THEN** 系統觸發通知，並記錄已通知狀態

#### Scenario: 持續在閾值內
- **WHEN** ETA 已在 <= threshold 且已通知
- **THEN** 系統不重複通知

#### Scenario: 下一班車進入
- **WHEN** ETA 回到 > threshold（前一班已進站）後再次跌到 <= threshold
- **THEN** 系統重新觸發通知（視為下一班車）

#### Scenario: 狀態重設條件
- **WHEN** ETA 狀態變為進站中、未發車、末班車已駛離、或 ETA > threshold
- **THEN** 系統重設跳變狀態為「未通知」，準備偵測下一班

### Requirement: 逐行 Log 輸出
監控模式中，系統 SHALL 每次 polling 後輸出一行包含時間戳和 ETA 資訊的 log。

#### Scenario: 正常 ETA 輸出
- **WHEN** 取得 ETA 為 8 分鐘
- **THEN** 輸出格式如 `14:30:01  ETA 8 分`

#### Scenario: 通知觸發時的輸出
- **WHEN** ETA 跌破閾值且觸發通知
- **THEN** 輸出格式如 `14:33:31  ETA 5 分  🔔 已通知！`

#### Scenario: 特殊狀態輸出
- **WHEN** ETA 為進站中、未發車等特殊狀態
- **THEN** 輸出對應的狀態文字（如 `進站中`、`未發車`、`末班駛離`）

### Requirement: 優雅停止
使用者按 Ctrl+C 時，系統 SHALL 優雅停止 polling 並顯示摘要後退出。

#### Scenario: Ctrl+C 停止
- **WHEN** 使用者按 Ctrl+C
- **THEN** 系統停止 polling，顯示監控摘要（總監控時間、通知次數），然後退出
