## Why

eBus 的 `GetStopDyns` 每站回傳兩個車輛陣列：`bi`（車在這一站）與 `bo`（車剛離開這一站，正在往下一站的路上）。`internal/ebus/convert.go` 只讀 `bi`，`bo` 解析出來後直接丟棄，造成兩個使用者可見的問題：

1. **車號會無故消失**：實測 299 去程的一次快照，17 個有車的站中有 9 站的車只存在於 `bo`，路上 20 台車有 10 台在畫面上完全看不到。車一離開第 N 站就從 N 的 `bi` 移到 N 的 `bo`，而 N+1 站要等車真的進站才會有 `bi`；這段期間 N+1 顯示「進站中／將到站」卻沒有任何車號，正是使用者回報「將到站、進站中反而看不到車號」的情形。
2. **看不出車在哪裡**：列表只有一行狀態文字，沒有任何視覺元素指出車輛目前的位置，使用者得逐行讀文字才能推測。

`bo` 的語意正好就是「車在兩站之間」，補上它同時解掉這兩個問題。

## What Changes

- **保留 eBus 的離站車輛**：`model.StopETA` 新增 `departedBuses` 欄位，承載 eBus `bo`（剛駛離該站、正前往下一站的車）。原有的 `buses` 語意收斂為「目前停靠在該站的車」。
- **前端改為時間軸版型**：`StopList` 左側加入站點時間軸，停靠中的車以公車圖示畫在站點節點上，離站的車畫在該站與下一站之間的連線段上，車號跟著圖示走。
- **車號不再消失**：每一台上游有回報的車，都會出現在站上或站間其中一處。
- TDX provider 的 ETA endpoint 不回傳車輛位置語意，`departedBuses` 一律為空，顯示行為與現況相同。

## Capabilities

### New Capabilities

（無新 capability）

### Modified Capabilities

- **bus-data-source**：統一資料模型的 `StopETA` 新增 `departedBuses`；eBus provider 的資料轉換需求明確要求 `bi` 與 `bo` 分別對應到停靠中與離站中的車輛。
- **go-api-proxy**：ETA endpoint 回應的每個 stop 新增 `departedBuses` 陣列。
- **pwa-frontend**：站點列表改為時間軸版型，車輛以圖示呈現於站點節點或站間連線。

## Impact

- `internal/model/model.go`：`StopETA` 新增 `DepartedBuses`
- `internal/ebus/convert.go`：`convertETAs` 轉換 `BO`
- `internal/tdx/convert.go`：不變（`DepartedBuses` 留空）
- `web/src/api/types.ts`：`StopETA` 新增 `departedBuses`
- `web/src/components/StopList.tsx`：改為時間軸版型 + 公車圖示
- `web/src/components/BusMarker.tsx`（新增）：公車圖示與車號的共用呈現
- `web/src/pages/HomePage.tsx`：收藏卡片維持現有版型，僅確認新欄位不影響顯示
- 不影響 CLI（`cmd/notify` 未使用車號）與 ETA 三級狀態邏輯
