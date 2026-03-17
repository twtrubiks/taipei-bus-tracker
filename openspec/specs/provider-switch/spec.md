## Purpose

CLI provider 模式切換，支援 auto、tdx、ebus 三種模式。

## Requirements

### Requirement: Provider mode selection via CLI flag
CLI SHALL accept `--provider <mode>` flag，支援 `auto`、`tdx`、`ebus` 三個值。未提供時使用 config/env/default 決定。

#### Scenario: Explicit tdx mode
- **WHEN** 使用者執行 `bus-notify --provider tdx` 且 config 有 TDX 憑證
- **THEN** 系統 SHALL 僅使用 TDX provider，無 fallback

#### Scenario: Explicit ebus mode
- **WHEN** 使用者執行 `bus-notify --provider ebus`
- **THEN** 系統 SHALL 僅使用 eBus provider，無 fallback，不需要 TDX 憑證

#### Scenario: Auto mode with TDX credentials
- **WHEN** mode 為 `auto`（或未指定）且 config 有 TDX 憑證
- **THEN** 系統 SHALL 使用 TDX 為 primary、eBus 為 fallback（等同現有行為）

#### Scenario: Auto mode without TDX credentials
- **WHEN** mode 為 `auto`（或未指定）且 config 無 TDX 憑證
- **THEN** 系統 SHALL 使用 eBus 為 primary，無 fallback

### Requirement: Error on tdx mode without credentials
選擇 `tdx` 模式但缺少 TDX 憑證時，系統 SHALL 報錯退出，不做 silent fallback。

#### Scenario: TDX mode missing credentials
- **WHEN** 使用者指定 `--provider tdx` 但 config 無 TDX client_id 或 client_secret
- **THEN** 系統 SHALL 輸出錯誤訊息並以非零狀態碼退出

#### Scenario: Invalid provider mode
- **WHEN** 使用者指定 `--provider foo`（無效值）
- **THEN** 系統 SHALL 輸出錯誤訊息並以非零狀態碼退出

### Requirement: Provider config field with override chain
Config 層 SHALL 支援 `provider` 欄位，優先順序為：CLI flag > `BUS_PROVIDER` 環境變數 > config.yaml `provider` 欄位 > 預設 `auto`。

#### Scenario: Environment variable overrides config
- **WHEN** config.yaml 設 `provider: ebus` 但環境變數 `BUS_PROVIDER=tdx`
- **THEN** 系統 SHALL 使用 `tdx` 模式

#### Scenario: CLI flag overrides everything
- **WHEN** config.yaml 設 `provider: ebus`、環境變數 `BUS_PROVIDER=tdx`、CLI flag `--provider auto`
- **THEN** 系統 SHALL 使用 `auto` 模式

#### Scenario: No provider specified anywhere
- **WHEN** config.yaml 無 `provider` 欄位、無環境變數、無 CLI flag
- **THEN** 系統 SHALL 使用 `auto` 模式（向後相容）

### Requirement: Centralized provider factory
provider 初始化邏輯 SHALL 集中在 `internal/provider` 套件的 `Build()` 函式中，CLI 和 Web server 共用同一函式。

#### Scenario: CLI uses Build function
- **WHEN** CLI 啟動並決定 provider mode
- **THEN** CLI SHALL 呼叫 `provider.Build(cfg, mode)` 取得 primary 和 fallback，不自行建構 provider

#### Scenario: Build returns structured result
- **WHEN** 呼叫 `provider.Build(cfg, mode)`
- **THEN** 函式 SHALL 回傳 `(primary BusDataSource, fallback BusDataSource, error)`
