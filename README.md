# Taipei Bus Tracker

自架的台北公車即時到站查詢工具。Go 後端 + React PWA 前端，單一 binary 部署。

## 功能

- 路線搜尋、方向選擇、即時到站時間
- 每 15 秒自動更新
- 雙資料源 fallback（TDX 主要 + eBus 備援）
- 收藏常用路線/站點
- 到站瀏覽器通知提醒
- 深色模式
- PWA 支援（可加到手機主畫面）

## 環境需求

- Go 1.22+
- Node.js 20+
- npm 10+

## 快速開始

```bash
# 複製設定檔並填入 TDX credentials
cp config.example.yaml config.yaml

# 一鍵 build
make build

# 啟動
./taipei-bus
```

## 開發

```bash
# Go 後端開發
go run ./cmd/server

# React 前端開發
cd web && npm run dev

# Lint (Go + React)
make lint

# Test (Go + React)
make test

# Lint + Test
make check
```

## 設定

支援 YAML config file 和環境變數，環境變數優先。

| 設定項 | 環境變數 | config.yaml key | 預設值 |
|--------|----------|-----------------|--------|
| Server Port | `BUS_PORT` | `port` | `8080` |
| TDX Client ID | `TDX_CLIENT_ID` | `tdx.client_id` | - |
| TDX Client Secret | `TDX_CLIENT_SECRET` | `tdx.client_secret` | - |
| Static Files Path | `BUS_STATIC_PATH` | `static_path` | `./static` |

## 部署

```bash
# 使用 systemd
sudo cp taipei-bus.service /etc/systemd/system/
sudo systemctl enable --now taipei-bus
```
