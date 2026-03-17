## Purpose

跨 provider 收藏的雙 ID 儲存與 lazy resolve 機制。

## Requirements

### Requirement: Dual-ID storage structure
收藏/快捷 SHALL 儲存雙組 ID（tdx_route_id / tdx_stop_id 和 ebus_route_id / ebus_stop_id），取代原本的單一 route_id / stop_id。

#### Scenario: 用 TDX 建立收藏
- **WHEN** 使用者在 TDX provider 下建立收藏
- **THEN** 系統 SHALL 將 route_id 存入 tdx_route_id、stop_id 存入 tdx_stop_id，ebus 欄位留空

#### Scenario: 用 eBus 建立收藏
- **WHEN** 使用者在 eBus provider 下建立收藏
- **THEN** 系統 SHALL 將 route_id 存入 ebus_route_id、stop_id 存入 ebus_stop_id，tdx 欄位留空

### Requirement: Lazy resolve on provider switch
當前 provider 對應的 ID 欄位為空時，系統 SHALL 用名稱反查補齊 ID，並回寫儲存。

#### Scenario: 首次切換 provider
- **WHEN** 載入收藏時，當前 provider 對應的 route_id 為空
- **THEN** 系統 SHALL 用 routeName 搜尋路線，用 stopName 比對站點，取得新 ID 並回寫儲存

#### Scenario: 第二次切換同一 provider
- **WHEN** 載入收藏時，當前 provider 對應的 route_id 已有值
- **THEN** 系統 SHALL 直接使用已存的 ID，不發起額外查詢

#### Scenario: 反查比對失敗
- **WHEN** 名稱反查無法精確比對到路線或站點
- **THEN** 系統 SHALL 輸出警告訊息，保留已有的 ID 可用，不阻斷使用

### Requirement: Auto mode fallback populates second ID
auto 模式 fallback 時，系統 SHALL 順便補齊 fallback provider 的 ID。

#### Scenario: TDX 失敗 fallback 到 eBus
- **WHEN** auto 模式下 TDX 查詢失敗，fallback 到 eBus 成功
- **THEN** 系統 SHALL 將 eBus 回傳的 route_id / stop_id 存入 ebus 欄位（如果為空）

### Requirement: Source tracking on Route and Stop
model.Route 和 model.Stop SHALL 包含 Source 欄位，標記資料來自哪個 provider。

#### Scenario: TDX 回傳的路線
- **WHEN** TDX provider 回傳路線資料
- **THEN** Route.Source SHALL 為 "tdx"

#### Scenario: eBus 回傳的路線
- **WHEN** eBus provider 回傳路線資料
- **THEN** Route.Source SHALL 為 "ebus"
