## Context

全新專案，沒有既有 codebase。目標是建立一個自架的台北公車即時到站查詢工具，部署在使用者自己的 Linux server 上。資料來源為 TDX 官方 API（主要）和 ebus.gov.taipei（備援）。前端為 React PWA，電腦和手機瀏覽器都能使用。

## Goals / Non-Goals

**Goals:**
- 即時查詢台北公車到站時間，每 15 秒自動更新
- 雙資料源 fallback：TDX 失敗時自動切換 eBus
- 電腦 + 手機瀏覽器都有良好體驗（PWA）
- 收藏常用路線/站，快速存取
- 單一 Go binary 部署，簡單易維護

**Non-Goals:**
- 不做全台灣公車（先做大台北，架構預留擴展）
- 不做地圖顯示公車位置（Phase 1 不做）
- 不做使用者帳號系統（本地收藏用 localStorage）
- 不上架 App Store / Google Play
- 不做多人共用（自己用，不考慮高併發）

## Decisions

### 1. 後端語言：Go

**選擇**: Go standard library + 輕量 router（chi 或 standard mux）

**理由**: 單一 binary 部署、內建 HTTP client/server、記憶體佔用低（~10-20MB）、對這種 proxy 型服務是最佳選擇。

**替代方案**: Rust（過度工程）、Python（runtime 依賴重、記憶體高）。

### 2. 資料源架構：Interface + Provider pattern

**選擇**: 定義 `BusDataSource` interface，TDX 和 eBus 各實作一個 provider。

```
BusDataSource interface
├── TDXProvider   (主要)
└── EBusProvider  (備援)
```

**理由**: 資料源可獨立測試、替換、擴展。Fallback 邏輯在上層處理，provider 只負責單一資料源。

### 3. 快取策略：In-memory cache with TTL

**選擇**: 簡單的 sync.Map + TTL（10 秒），不用 Redis。

**理由**: 自己用，單一 process，不需要分散式快取。10 秒 TTL 可以減少對上游的請求，又不會太過時。

### 4. 前端框架：React + TypeScript + Vite

**選擇**: React 19 + TypeScript + Vite + Tailwind CSS。

**理由**: 生態成熟、PWA 支援好、Vite 開發體驗快。Tailwind 快速做出好看的 responsive UI。

**替代方案**: Vue（也可以，但 React 生態更大）、Flutter Web（初次載入 2MB，且偏離網頁原生體驗）。

### 5. PWA 策略：Service Worker + manifest

**選擇**: Workbox 產生 Service Worker，快取靜態資源。API 資料不做離線快取。

**理由**: 靜態資源離線可用，但公車到站資料離線沒意義。加到主畫面後全螢幕開啟，像 App 一樣。

### 6. 前後端部署：Go serve 靜態檔案

**選擇**: Go backend 同時 serve API 和 PWA 靜態檔案，單一 port。

**理由**: 一個 binary、一個 port、沒有 CORS 問題。前端 build 產物直接放在 Go binary 旁邊的 `static/` 目錄。

### 8. eBus HTML Scraping：補齊 SearchRoutes 與 GetStops

**選擇**: 透過 scraping ebus.gov.taipei 的 HTML 頁面，為 eBus provider 補上 SearchRoutes 和 GetStops。

**端點**:
- `POST /Query/QBusRoute`：搜尋路線，傳入 `QueryModel.QueryString`，回傳 HTML 列表
- `GET /Route/StopsOfRoute?routeId=xxx`：取得路線站序頁面，去程在 `#GoDirectionRoute`、回程在 `#BackDirectionRoute`

**解析方式**: 正則表達式解析 HTML（不引入 HTML parser 依賴）

**CSRF token**: 從 `/Query/BusRoute` 頁面取得，與現有 GetETA 的 token 共用同一管理機制

**理由**: TDX 憑證未核發前，eBus 是唯一資料源。原本 eBus provider 只有 GetETA，缺少 SearchRoutes 和 GetStops 導致前端無法使用。透過 scraping 補齊三個方法，讓 eBus 可以獨立運作。

**風險**: HTML 結構變動會導致解析失敗，但 eBus 本身就是非正式 API，此風險已在預期內。

### 7. 設定檔：YAML 或環境變數

**選擇**: 支援環境變數 + YAML config file，環境變數優先。

**理由**: TDX Client ID/Secret 用環境變數（安全）、server port 等用 config file。

## Risks / Trade-offs

- **[TDX API 額度]** 免費 20,000 次/天，每 15 秒查一次路線 = 5,760 次/天/路線 → 如果同時監看 3 條以上路線可能超額 → **Mitigation**: 快取 + 只在前端有人看時才輪詢，沒人看就停止
- **[eBus CSRF token 過期]** token 可能有時效 → **Mitigation**: token 管理器定期重新取得，失敗時重試一次
- **[eBus 改版]** 非正式 API，官方改版就壞 → **Mitigation**: eBus 只是備援，壞了還有 TDX
- **[TDX API 認證方式變更]** TDX 已從舊的 HMAC 改為 OAuth Client Credentials → **Mitigation**: 實作時確認最新認證文件
