## Context

目前 CLI（`cmd/notify/main.go`）和 Web server（`cmd/server/main.go`）各自在 main() 中硬編碼 provider 選擇邏輯：有 TDX 憑證就 primary=TDX + fallback=eBus，否則 primary=eBus + fallback=nil。使用者無法指定只用某個 provider。

兩個 cmd 的初始化邏輯幾乎相同，但各寫一份。本次新增 `--provider` 開關並將邏輯集中到共用套件。

## Goals / Non-Goals

**Goals:**
- 使用者可透過 CLI flag / config / 環境變數選擇 provider 模式（auto / tdx / ebus）
- provider 初始化邏輯只寫一次，CLI 和 Web server 共用
- 選 `tdx` 但缺憑證時明確報錯

**Non-Goals:**
- 本次不改 Web server（`cmd/server`），僅預留共用函式
- 不改 FallbackService、ebus、tdx、model 等既有套件
- 不支援執行期動態切換 provider

## Decisions

### 1. 新增 `internal/provider` 套件作為工廠

**選擇**: 建立 `internal/provider/provider.go`，提供 `Build(cfg, mode)` 函式。

**替代方案**:
- 放在 `internal/config`：config 不應 import ebus/tdx（依賴方向錯誤）
- 放在 `internal/handler`：handler 目前只認 interface，保持乾淨
- 各 cmd 各寫一份：重複且日後容易不一致

**理由**: provider 套件只做一件事——根據設定組裝 provider pair，職責清楚，依賴方向正確（provider → config, ebus, tdx, model）。

### 2. 三層覆蓋優先順序：flag > env > config > default

**選擇**: CLI flag `--provider` 優先，其次 `BUS_PROVIDER` 環境變數，再來 config.yaml 的 `provider` 欄位，最後預設 `auto`。

**理由**: 與現有 config 的 env override 模式一致（如 `TDX_CLIENT_ID`）。flag 最臨時、env 適合容器部署、config 適合本地持久化。

### 3. 模式對應表

| mode | primary | fallback | 備註 |
|------|---------|----------|------|
| `auto` | TDX（若有 key）或 eBus | eBus（若 primary 是 TDX）或 nil | 等同現有行為 |
| `tdx` | TDX | nil | 無 key 則 error |
| `ebus` | eBus | nil | |

### 4. 錯誤處理策略

選 `tdx` 但 config 無 TDX 憑證 → `Build()` 回傳 error，由呼叫端 `log.Fatalf` 退出。不做 silent fallback，讓使用者明確知道設定問題。

## Risks / Trade-offs

- **新套件成本**: 多一個 package，但只有一個檔案、一個函式，複雜度極低 → 可接受
- **config.yaml 不向後相容風險**: 無。新欄位 `provider` 預設空值等同 `auto`，完全相容既有設定檔
