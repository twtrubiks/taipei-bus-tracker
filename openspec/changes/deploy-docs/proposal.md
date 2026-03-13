## Why

README 的部署段落過於簡略，缺少完整的 server 建置流程、靜態編譯說明、nginx 設定、以及常見問題排除。使用者（含未來的自己）在新 server 上部署時需要一份可直接照做的文件。

## What Changes

- 補充 README 部署段落：完整的建置、複製、啟動流程
- 新增靜態編譯說明（CGO_ENABLED=0 解決 glibc 版本不符）
- 新增 server 端 build 說明（`-buildvcs=false`）
- 補充 nginx 反向代理 + HTTPS 設定說明
- 補充 config.yaml / 環境變數設定範例
- 補充常見問題排除（glibc 不符、provider 切換、僅 HTTP 測試）

## Capabilities

### New Capabilities

（無新增 capability，純文件補充）

### Modified Capabilities

（無 spec 層級的行為變更）

## Impact

- `README.md`：部署段落重寫擴充
- `deploy/` 目錄：可能新增說明或範例檔案
