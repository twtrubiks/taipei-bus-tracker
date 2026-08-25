## MODIFIED Requirements

### Requirement: API 取得即時到站時間
系統 SHALL 提供 `GET /api/routes/{routeId}/eta?gb={direction}` endpoint，回傳指定路線某方向所有站的即時到站預估時間。回應中的到站資訊 SHALL 只包含秒數（`eta`），不包含人可讀的狀態文字——狀態文字由消費端（Web 前端、CLI）依秒數自行計算。

#### Scenario: 取得到站時間成功
- **WHEN** client 發送 `GET /api/routes/0100000100/eta?gb=0`
- **THEN** 系統回傳 JSON，包含 route、direction、source（"tdx" 或 "ebus"）、updatedAt、stops array。每個 stop 包含 stopId、stopName、sequence、eta（秒）、buses array

#### Scenario: 回應不含 status 欄位
- **WHEN** client 取得任一筆 ETA 回應
- **THEN** stops 中的每個物件 SHALL NOT 包含 `status` 欄位

#### Scenario: 到站秒數的特殊值
- **WHEN** 上游回報非即時到站的狀態
- **THEN** `eta` SHALL 以負數表示：-1 未發車、-2 末班車已駛離、-3 交管不停靠、-4 未營運。正數與 0 一律代表實際到站秒數
