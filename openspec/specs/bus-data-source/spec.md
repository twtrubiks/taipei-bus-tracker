## Purpose

公車資料源抽象層，定義統一介面與 TDX/eBus provider 實作及 fallback 機制。

## Requirements

### Requirement: BusDataSource 統一介面
系統 SHALL 定義 BusDataSource interface，包含以下方法：SearchRoutes(city, keyword)、GetStops(city, routeId, direction)、GetETA(city, routeId, direction)。所有 provider MUST 實作此 interface。

#### Scenario: Provider 可替換
- **WHEN** 系統初始化時配置不同的 provider
- **THEN** 上層業務邏輯不需要修改，因為所有 provider 都實作相同的 interface

### Requirement: 統一資料模型
所有 provider SHALL 將上游資料轉換為統一的資料模型：Route（routeId, name, startStop, endStop）、Stop（stopId, name, sequence）、StopETA（stopId, stopName, sequence, eta, buses, departedBuses, source）。`buses` SHALL 只包含目前停靠在該站的車輛，`departedBuses` SHALL 只包含已駛離該站、正前往下一站的車輛。上游若未提供車輛位置語意，`departedBuses` SHALL 為空。

#### Scenario: TDX 和 eBus 回傳相同格式
- **WHEN** 分別透過 TDX 和 eBus provider 查詢同一路線的 ETA
- **THEN** 回傳的 StopETA 結構相同，僅 source 欄位不同（"tdx" vs "ebus"）

#### Scenario: 上游無離站車輛資訊
- **WHEN** provider 的上游 API 不區分車輛停靠中與離站中（如 TDX 的到站預估 endpoint）
- **THEN** 該 provider 回傳的每筆 StopETA 的 `departedBuses` SHALL 為空，`buses` 行為不變

### Requirement: TDX Provider
TDX provider SHALL 使用 TDX API 的 OAuth Client Credentials 流程取得 Access Token，並用於後續 API 呼叫。Token SHALL 在過期前自動更新。

#### Scenario: TDX 認證成功
- **WHEN** 系統啟動且 TDX Client ID/Secret 有效
- **THEN** provider 成功取得 Access Token 並能查詢 API

#### Scenario: TDX Token 自動更新
- **WHEN** Access Token 即將過期（剩餘 5 分鐘內）
- **THEN** provider SHALL 自動重新取得新的 Token，不中斷服務

### Requirement: eBus Provider
eBus provider SHALL 透過 GET ebus.gov.taipei 頁面取得 CSRF token，再用 POST /EBus/GetStopDyns 取得到站資料。回應中每站的 `bi` 與 `bo` SHALL 分別對應到停靠中與離站中的車輛，兩者皆 SHALL NOT 被丟棄。

#### Scenario: eBus CSRF token 取得
- **WHEN** provider 首次請求或 token 過期
- **THEN** provider SHALL GET 頁面 HTML，從 hidden input 中解析 __RequestVerificationToken

#### Scenario: eBus 資料轉換
- **WHEN** eBus 回傳 `{sn, eta, bi, bo}` 格式
- **THEN** provider SHALL 轉換為統一的 StopETA 模型，eta 值語義保持一致

#### Scenario: bi 轉為停靠中車輛
- **WHEN** 某站的 `bi` 含有車輛
- **THEN** 這些車輛的車號 SHALL 出現在該站 StopETA 的 `buses` 中

#### Scenario: bo 轉為離站中車輛
- **WHEN** 某站的 `bo` 含有車輛
- **THEN** 這些車輛的車號 SHALL 出現在**該站**（而非下一站）StopETA 的 `departedBuses` 中

#### Scenario: 車輛不因位於站間而遺失
- **WHEN** 某站的 `bi` 為 null 但 `bo` 含有車輛
- **THEN** 該站的 StopETA SHALL 仍帶有這些車輛的車號於 `departedBuses`，不得回傳無任何車輛資訊

### Requirement: Fallback 機制
系統 SHALL 以 TDX 為主要資料源。當 TDX 請求失敗（逾時 5 秒、HTTP 5xx、或額度用盡）時，SHALL 自動切換到 eBus provider。

#### Scenario: TDX 正常時使用 TDX
- **WHEN** TDX API 正常回應
- **THEN** 系統使用 TDX 資料，response 中 source = "tdx"

#### Scenario: TDX 失敗時 fallback 到 eBus
- **WHEN** TDX API 回應逾時或 HTTP 5xx
- **THEN** 系統自動使用 eBus provider 取得資料，response 中 source = "ebus"

#### Scenario: 雙源都失敗時回傳快取
- **WHEN** TDX 和 eBus 都失敗，但快取中有該路線的資料
- **THEN** 系統回傳快取資料，並附上 stale = true 標記
