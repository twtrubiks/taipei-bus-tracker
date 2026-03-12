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

### 1. 建置與複製

```bash
make build

# 複製到 server
scp taipei-bus config.yaml your-server:/opt/taipei-bus-tracker/
scp -r static/ your-server:/opt/taipei-bus-tracker/static/
```

### 2. systemd 服務

```bash
sudo cp deploy/taipei-bus-tracker.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now taipei-bus-tracker
```

### 3. nginx + HTTPS (Let's Encrypt)

```bash
# 安裝 certbot
sudo apt install nginx certbot python3-certbot-nginx

# 取得 SSL 憑證
sudo certbot --nginx -d your-domain.com

# 複製 nginx 設定（修改 your-domain.com 為實際域名）
sudo cp deploy/nginx-taipei-bus.conf /etc/nginx/sites-available/taipei-bus
sudo ln -s /etc/nginx/sites-available/taipei-bus /etc/nginx/sites-enabled/
sudo nginx -t && sudo systemctl reload nginx
```

設定檔在 `deploy/` 目錄：
- `taipei-bus-tracker.service` — systemd 服務
- `nginx-taipei-bus.conf` — nginx 反向代理 + HTTPS
