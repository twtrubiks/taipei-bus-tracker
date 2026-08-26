## 1. 後端保留 eBus 離站車輛

- [x] 1.1 `internal/model/model.go`：`StopETA` 新增 `DepartedBuses []Bus`（`json:"departedBuses"`）
- [x] 1.2 `internal/ebus/convert.go`：`convertETAs` 將 `s.BO` 轉為 `DepartedBuses`，`s.BI` 維持 `Buses`
- [x] 1.3 `internal/ebus/convert_test.go`：補上只有 `bo`、只有 `bi`、兩者皆有、兩者皆空四種情境，以及離站車輛歸屬來源站、空車號略過
- [x] 1.4 確認 `internal/tdx/convert.go` 的 `DepartedBuses` 為空且測試涵蓋

## 2. API 回應與前端型別

- [x] 2.1 `web/src/api/types.ts`：`StopETA` 新增 `departedBuses?: Bus[]`
- [x] 2.2 `internal/handler/routes_test.go`：驗證 ETA 回應含 `departedBuses` 欄位且不被搬到下一站
- [x] 2.3 用 `cmd/test-ebus` 實際打 API 驗證 `departedBuses` 有值（同時擴充該工具輸出站上／站間統計）

## 3. 前端公車圖示與時間軸版型

- [x] 3.1 新增 `web/src/components/BusMarker.tsx`：inline SVG 公車圖示，支援 `placement="at-stop" | "between"`
- [x] 3.2 新增 `web/src/utils/plates.ts`：車號字串格式化，供 StopList 與 HomePage 共用
- [x] 3.3 `StopList.tsx` 改為時間軸版型：左側節點與連線，站序保留在內容欄
- [x] 3.4 停靠中車輛（`buses`）繪於站點節點
- [x] 3.5 離站中車輛（`departedBuses`）繪於該站與下一站之間的連線段
- [x] 3.6 末站的 `departedBuses` 繪於末站之後，不隱藏
- [x] 3.7 深色模式下圖示與連線顏色檢查；車號對比由 `gray-400` 提高為 `gray-500 dark:gray-400`
- [x] 3.8 `HomePage.tsx` 收藏卡片維持單站版型，改用共用的 `plateText`

## 4. 前端測試

- [x] 4.1 `StopList.test.tsx`：只有 `bo` 時車號顯示於站間
- [x] 4.2 `StopList.test.tsx`：同站同時有 `bi` 與 `bo` 時顯示兩個圖示
- [x] 4.3 `StopList.test.tsx`：末站 `bo` 仍顯示
- [x] 4.4 `StopList.test.tsx`：兩者皆空時不顯示公車圖示；TDX 無 `departedBuses` 欄位時正常
- [x] 4.5 `BusMarker.test.tsx` 與 `plates.test.ts`：兩種 placement 的渲染與無障礙標籤

## 5. 驗證

- [x] 5.1 `make lint` 通過
- [x] 5.2 `make test` 通過（Go 全數 ok + React 95 tests）
- [x] 5.3 手動驗證：eBus 299 去程，畫面車輛數與 API `bi`+`bo` 總數一致（站上 14 台、站間 7 台）
- [x] 5.4 手動驗證：原本無車號的「將到站／進站中」站點，前方連線段已出現車輛圖示與車號
- [x] 5.5 手動驗證：手機寬度（430px）與深色／淺色模式版型正常
- [x] 5.6 手動驗證 TDX 來源路線僅顯示站上圖示 — **已於 2026-08-26 補驗**：`config.yaml` 填入 TDX 憑證後，以 `BUS_PORT=8099 BUS_PROVIDER=tdx` 啟動 server，打 `/api/routes/TPE11411/eta?gb=0`（299 去程）取得 59 站，`departedBuses` 全為 `null`（非空站數 0），確認 TDX 不會產生站間圖示；`internal/tdx/convert.go` 的 `convertETAs` 亦未賦值 `DepartedBuses`。整合測試 `go test -tags=integration ./internal/tdx/` 三項全數通過
  - **補充發現**：TDX 路線實際上連站上圖示也不會出現。`convertETAs` 僅在 `PlateNumb != ""` 時填入 `Buses`，而 TDX 的 ETA API 不回傳 `PlateNumb`（59 站 `buses` 全為 `null`，`sequence` 亦全為 0）。故 TDX 來源的正確描述是「只有倒數秒數、無任何公車圖示」，本項標題的「僅顯示站上圖示」措辭偏樂觀
