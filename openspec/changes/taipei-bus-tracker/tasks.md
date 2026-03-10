## 1. 專案初始化

- [x] 1.0 `git init` + `.gitignore`（Go binary、node_modules/、dist/、.vite/、config.yaml、.env、.DS_Store、*.log）
- [x] 1.0.1 建立 `README.md`（專案說明、環境需求、啟動方式、開發指令）
- [x] 1.1 建立 Go module（`go mod init`）、專案目錄結構（handler/, ebus/, tdx/, cache/, config/）
- [x] 1.2 建立 React 專案（Vite + TypeScript + Tailwind CSS），設定 PWA manifest 和 Service Worker，設定 Vitest + React Testing Library + MSW
- [x] 1.3 建立 config.example.yaml 和設定載入邏輯（YAML + 環境變數，環境變數優先）
- [x] 1.4 Go 品質工具設定：golangci-lint 設定檔（.golangci.yml）、go vet、go test 基礎結構
- [x] 1.5 React 品質工具設定：ESLint + Prettier 設定、TypeScript strict mode
- [x] 1.6 建立 Makefile lint/test targets：`make lint`（Go + React）、`make test`（Go test + Vitest）、`make check`（lint + test 一次跑完）
- [x] 1.7 設定 pre-commit hook：使用 husky 或 shell script，commit 前自動跑 `make lint`，失敗則阻擋 commit

## 2. Go 後端 - 資料源抽象層

- [x] 2.1 定義 BusDataSource interface 和統一資料模型（Route, Stop, StopETA, Bus struct）
- [x] 2.2 [TDD] 先寫 TDX 資料轉換測試：TDX JSON → 統一模型的欄位對應、邊界值
- [x] 2.3 實作 TDX provider：OAuth Client Credentials 認證、Token 自動更新
- [x] 2.4 實作 TDX provider：SearchRoutes、GetStops、GetETA 三個方法 + 資料轉換
- [x] 2.5 [TDD] 先寫 eBus 資料轉換測試：eBus `{sn, eta, bi, bo}` → 統一模型的欄位對應
- [x] 2.6 實作 eBus provider：CSRF token 取得與管理
- [x] 2.7 實作 eBus provider：GetStopDyns 呼叫與資料轉換為統一模型
- [x] 2.8 [TDD] 先寫 Fallback 測試：TDX 成功/失敗 × eBus 成功/失敗 × 有無快取，共 6 種情境
- [x] 2.9 實作 Fallback 機制：TDX 失敗（逾時 5 秒 / 5xx / 額度用盡）自動切換 eBus

## 2.5. eBus 補齊 SearchRoutes 與 GetStops（HTML Scraping）

- [x] 2.10 [TDD] 先寫 eBus SearchRoutes HTML 解析測試：解析搜尋結果 HTML → []Route（routeId、routeName、startStop、endStop）
- [x] 2.11 實作 eBus SearchRoutes：POST /Query/QBusRoute + HTML 解析，CSRF token 共用現有機制
- [x] 2.12 [TDD] 先寫 eBus GetStops HTML 解析測試：解析站序頁面 HTML → []Stop（stopId、stopName、sequence），含去程/回程區分
- [x] 2.13 實作 eBus GetStops：GET /Route/StopsOfRoute + HTML 解析
- [x] 2.14 整合測試：用 spike 驗證 SearchRoutes → GetStops → GetETA 完整流程

## 3. Go 後端 - API 與快取

- [x] 3.1 [TDD] 先寫 cache 測試：Set/Get、TTL 過期、並發安全
- [x] 3.2 實作 in-memory cache（sync.Map + TTL 10 秒）
- [x] 3.3 [TDD] 先寫 ETA 狀態轉換測試：eta=300→"約5分"、eta=60→"進站中"、eta=-1→"未發車"、eta=-2→"末班車已駛離"、eta=-3→"交管不停靠"、eta=-4→"未營運"
- [x] 3.4 實作 `GET /api/routes/search` handler
- [x] 3.5 實作 `GET /api/routes/{routeId}/stops` handler
- [x] 3.6 實作 `GET /api/routes/{routeId}/eta` handler，包含 ETA → status 字串轉換
- [x] 3.7 補寫 API handler 測試：驗證 status code + response JSON 格式（httptest）
- [x] 3.8 實作統一錯誤回應格式和 HTTP status code
- [x] 3.9 實作靜態檔案 serve + SPA fallback（非 /api/* 路由回傳 index.html）
- [x] 3.10 串接 main.go：載入 config → 初始化 providers → 啟動 HTTP server

## 4. React PWA 前端 - 核心頁面

- [x] 4.1 設定 React Router，建立頁面結構：首頁（收藏）、搜尋頁、路線詳情頁
- [x] 4.2 實作路線搜尋元件：搜尋框 + 即時匹配結果列表
- [x] 4.2.1 [TDD] 搜尋元件測試：輸入關鍵字顯示匹配結果、空輸入清空列表、點擊結果導航、無結果提示
- [x] 4.3 實作方向選擇元件：顯示去程/回程兩個方向供選擇
- [x] 4.4 實作站點列表元件：顯示所有站的到站資訊（站名、ETA 狀態、車牌）
- [x] 4.4.1 [TDD] ETA 狀態渲染測試：eta=300→"約5分"、eta=60→"進站中"、eta=-1→"未發車"、eta=-2→"末班車已駛離"、eta=-3→"交管不停靠"、eta=-4→"未營運"、有車牌時顯示車牌
- [x] 4.5 實作 useEta hook：每 15 秒自動輪詢，頁面不可見時停止
- [x] 4.5.1 [TDD] useEta hook 測試：掛載後立即請求、15 秒自動輪詢、頁面不可見時停止、重新可見時恢復、API 錯誤處理、unmount 清理 timer

## 5. React PWA 前端 - 收藏與通知

- [x] 5.1 實作收藏功能：新增/移除收藏到 localStorage，收藏「路線+方向+站」組合
- [x] 5.1.1 收藏功能測試：新增收藏寫入 localStorage、重新載入還原收藏、移除收藏同步更新 localStorage 與畫面
- [x] 5.2 實作首頁收藏面板：顯示所有收藏站的即時到站狀態（批次查詢）
- [x] 5.3 實作到站通知：設定 X 分鐘前提醒，ETA 到達時觸發 Browser Notification
- [x] 5.3.1 [TDD] 通知觸發測試：ETA ≤ 閾值觸發通知、ETA > 閾值不觸發、同一站不重複觸發、權限被拒顯示提示、通知內容正確（路線名+站名+到站時間）
- [x] 5.4 實作深色模式：跟隨系統 + 手動切換，選擇存 localStorage

## 6. React PWA 前端 - Responsive 與 PWA

- [ ] 6.1 Responsive 排版：桌面寬版 (>768px) + 手機窄版 (≤768px)
- [ ] 6.2 設定 PWA manifest.json（name、icons、display: standalone、theme_color）
- [ ] 6.3 設定 Service Worker（Workbox）：快取靜態資源，API 走網路

## 7. 部署

- [ ] 7.1 建立 Makefile：`make build` 一鍵完成 Go build + React build + 複製 static/
- [ ] 7.2 建立 systemd service file（自動啟動、crash 5 秒後重啟、日誌管理）
- [ ] 7.3 端到端測試：在 server 上部署，電腦 + 手機瀏覽器驗證所有功能
