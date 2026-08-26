## Purpose

Go 後端 API，提供路線搜尋、站序查詢、即時到站時間等 RESTful endpoints。

## Requirements

### Requirement: API 路線搜尋
系統 SHALL 提供 `GET /api/routes/search?q={keyword}&city={city}` endpoint，根據關鍵字搜尋公車路線。city 預設為 Taipei。回傳路線列表包含 routeId、routeName、起站、迄站。

#### Scenario: 搜尋路線成功
- **WHEN** client 發送 `GET /api/routes/search?q=1`
- **THEN** 系統回傳 JSON array，包含所有路線名稱含 "1" 的路線資訊，HTTP 200

#### Scenario: 搜尋無結果
- **WHEN** client 發送 `GET /api/routes/search?q=不存在的路線`
- **THEN** 系統回傳空 JSON array `[]`，HTTP 200

### Requirement: API 取得站序
系統 SHALL 提供 `GET /api/routes/{routeId}/stops?gb={direction}` endpoint，回傳指定路線某方向的所有站點及站序。direction 為 0（去程）或 1（回程）。

#### Scenario: 取得站序成功
- **WHEN** client 發送 `GET /api/routes/0100000100/stops?gb=0`
- **THEN** 系統回傳該路線去程的所有站點 JSON array，每筆包含 stopId、stopName、sequence，按站序排列，HTTP 200

### Requirement: API 取得即時到站時間
系統 SHALL 提供 `GET /api/routes/{routeId}/eta?gb={direction}` endpoint，回傳指定路線某方向所有站的即時到站預估時間。回應中的到站資訊 SHALL 只包含秒數（`eta`），不包含人可讀的狀態文字——狀態文字由消費端（Web 前端、CLI）依秒數自行計算。每個站 SHALL 分別回傳停靠中車輛（`buses`）與離站中車輛（`departedBuses`）。

#### Scenario: 取得到站時間成功
- **WHEN** client 發送 `GET /api/routes/0100000100/eta?gb=0`
- **THEN** 系統回傳 JSON，包含 route、direction、source（"tdx" 或 "ebus"）、updatedAt、stops array。每個 stop 包含 stopId、stopName、sequence、eta（秒）、buses array、departedBuses array

#### Scenario: 回應不含 status 欄位
- **WHEN** client 取得任一筆 ETA 回應
- **THEN** stops 中的每個物件 SHALL NOT 包含 `status` 欄位

#### Scenario: 到站秒數的特殊值
- **WHEN** 上游回報非即時到站的狀態
- **THEN** `eta` SHALL 以負數表示：-1 未發車、-2 末班車已駛離、-3 交管不停靠、-4 未營運。正數與 0 一律代表實際到站秒數

#### Scenario: 離站車輛歸屬於來源站
- **WHEN** 上游回報某車輛已駛離第 N 站
- **THEN** 該車號 SHALL 出現在 stops 中 sequence 為 N 的物件的 `departedBuses`，SHALL NOT 被移動到下一站

### Requirement: API 回應格式統一
所有 API endpoint SHALL 回傳統一的 JSON 格式。錯誤時 SHALL 回傳 `{"error": "message"}` 搭配適當的 HTTP status code。

#### Scenario: 上游 API 全部失敗
- **WHEN** TDX 和 eBus 都無法回應
- **THEN** 系統回傳 HTTP 502，body 為 `{"error": "upstream unavailable"}`

### Requirement: API 快取
系統 SHALL 對上游 API 回應做 in-memory 快取，TTL 為 10 秒。相同參數的請求在 TTL 內 SHALL 直接回傳快取資料，不重複請求上游。

#### Scenario: 快取命中
- **WHEN** 同一路線同一方向的 ETA 請求在 10 秒內重複發送
- **THEN** 第二次請求 SHALL 直接回傳快取資料，不呼叫上游 API

### Requirement: 靜態檔案 serve
Go server SHALL 同時 serve `/api/*` 路由和 PWA 靜態檔案。非 `/api/*` 的請求 SHALL fallback 到 `index.html`（SPA routing）。

#### Scenario: 存取前端頁面
- **WHEN** client 發送 `GET /`
- **THEN** 系統回傳 PWA 的 index.html

#### Scenario: SPA 路由 fallback
- **WHEN** client 發送 `GET /route/0100000100`（前端路由）
- **THEN** 系統回傳 index.html，由前端 JavaScript 處理路由

### Requirement: 設定管理
系統 SHALL 支援 YAML config file 和環境變數兩種設定方式，環境變數優先覆蓋 config file。必要設定項：TDX Client ID、TDX Client Secret、server port、static files path。

#### Scenario: 環境變數覆蓋 config
- **WHEN** config.yaml 設定 port 為 8080，但環境變數 `BUS_PORT=9090`
- **THEN** server SHALL 啟動在 port 9090
