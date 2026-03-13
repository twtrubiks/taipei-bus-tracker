## 1. README 部署段落重寫

- [ ] 1.1 重寫「建置」段落：本機 `make build` + 靜態編譯 `CGO_ENABLED=0`
- [ ] 1.2 新增「Server 端 build」段落：在 server 上安裝 Go/Node.js 並直接 build（含 `-buildvcs=false`）
- [ ] 1.3 重寫「複製到 server」段落：binary + static/ + config.yaml
- [ ] 1.4 補充 systemd 段落：完整的安裝、啟用、查看日誌指令
- [ ] 1.5 補充 nginx + HTTPS 段落：安裝 certbot、取得憑證、設定反向代理
- [ ] 1.6 新增「僅 HTTP 測試」段落：不需 nginx 直接 `./taipei-bus` 測試
- [ ] 1.7 新增「Provider 切換」段落：config.yaml 或環境變數切換 ebus/tdx/auto

## 2. 常見問題

- [ ] 2.1 新增「常見問題」段落：glibc 版本不符、VCS stamping 錯誤、port 佔用
