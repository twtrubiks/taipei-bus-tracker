## Why

目前 provider 選擇邏輯硬編碼在各 cmd 的 main.go 中（有 TDX key 就用 TDX + eBus fallback，沒有就只用 eBus），使用者無法控制要用哪個資料來源。開發除錯時需要隔離測試單一 provider，或在某個 provider 不穩定時手動切換，目前都做不到。

## What Changes

- 新增 `--provider` CLI flag，可指定 `auto`（預設）、`tdx`、`ebus` 三種模式
- config.yaml 新增 `provider` 欄位，作為持久化預設值
- 支援 `BUS_PROVIDER` 環境變數覆蓋
- 優先順序：CLI flag > 環境變數 > config.yaml > `auto`
- 選 `tdx` 但無 TDX 憑證時，直接報錯退出
- 新增 `internal/provider` 套件，將 provider 初始化邏輯集中為工廠函式，供 CLI 和 Web server 共用
- 本次先改 CLI（`cmd/notify`），Web server 暫不改動

## Capabilities

### New Capabilities
- `provider-switch`: CLI flag / config / 環境變數控制 provider 模式（auto / tdx / ebus），含工廠函式與錯誤處理

### Modified Capabilities

（無既有 spec 需要修改）

## Impact

- `internal/config/config.go`：Config struct 新增 Provider 欄位
- `internal/provider/provider.go`（新檔案）：Build() 工廠函式
- `cmd/notify/main.go`：加 `--provider` flag，改用 `provider.Build()`
- `config.yaml` / `config.example.yaml`：新增 provider 欄位範例
- handler、ebus、tdx、model 套件不受影響
