## 1. Config 層擴充

- [x] 1.1 `internal/config/config.go`：Config struct 新增 `Provider string` 欄位（yaml tag `provider`）
- [x] 1.2 `internal/config/config.go`：Load() 加 `BUS_PROVIDER` 環境變數 override
- [x] 1.3 `config.yaml` 和 `config.example.yaml`：新增 `provider: auto` 欄位與註解

## 2. Provider 工廠套件

- [x] 2.1 新增 `internal/provider/provider.go`：實作 `Build(cfg *config.Config, mode string) (primary, fallback model.BusDataSource, err error)`
- [x] 2.2 Build() 處理三種模式（auto / tdx / ebus）與錯誤情境（無效 mode、tdx 缺憑證）

## 3. CLI 整合

- [x] 3.1 `cmd/notify/main.go`：新增 `--provider` flag
- [x] 3.2 `cmd/notify/main.go`：flag 覆蓋 config 的優先順序邏輯（flag 有值用 flag，否則用 cfg.Provider）
- [x] 3.3 `cmd/notify/main.go`：移除舊的 provider 硬編碼初始化，改呼叫 `provider.Build()`
- [x] 3.4 啟動時印出目前使用的 provider 模式（如「Provider 模式: tdx」）

## 4. 驗證

- [x] 4.1 `go build ./...` 編譯通過
- [x] 4.2 `go test ./...` 全部通過
- [x] 4.3 手動測試：`--provider ebus`、`--provider tdx`、`--provider auto`、不帶 flag 各跑一次確認行為正確
