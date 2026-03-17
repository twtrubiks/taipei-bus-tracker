## Purpose

部署相關配置，包含單一 binary build、systemd service 與 config 管理。

## Requirements

### Requirement: 單一 Binary 部署
Go backend SHALL 編譯為單一 binary，搭配 static/ 目錄（PWA build 產物）即可運行，不需額外 runtime 或 container。

#### Scenario: 部署流程
- **WHEN** 將 binary + static/ + config.yaml 放到 server
- **THEN** 執行 `./taipei-bus` 即可啟動服務

### Requirement: Systemd Service
專案 SHALL 提供 systemd service file，支援自動啟動、自動重啟、日誌管理。

#### Scenario: 系統開機自動啟動
- **WHEN** server 重啟
- **THEN** taipei-bus service 自動啟動

#### Scenario: Crash 自動重啟
- **WHEN** process crash
- **THEN** systemd SHALL 在 5 秒後自動重啟 service

### Requirement: Build Script
專案 SHALL 提供 Makefile 或 build script，一鍵完成 Go build + React build + 複製靜態檔案。

#### Scenario: 一鍵 build
- **WHEN** 執行 `make build`
- **THEN** 產出 Go binary 和 static/ 目錄，可直接部署

### Requirement: Config 範例
專案 SHALL 提供 config.example.yaml，列出所有可設定項及說明。

#### Scenario: 新部署設定
- **WHEN** 使用者首次部署
- **THEN** 複製 config.example.yaml 為 config.yaml，填入 TDX credentials 即可使用
