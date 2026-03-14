## Context

目前 README 部署段落只有 3 行指令，缺少完整的建置和部署流程。實際部署時遇到了 glibc 版本不符、VCS stamping 錯誤等問題，這些經驗應記錄在文件中。

現有部署相關檔案：
- `Makefile`：build-web + build-go + 複製 static/
- `deploy/taipei-bus-tracker.service`：systemd service
- `deploy/nginx-taipei-bus.conf`：nginx 反向代理 + HTTPS
- `config.example.yaml`：設定範例

## Goals / Non-Goals

**Goals:**
- README 部署段落涵蓋從 build 到上線的完整流程
- 記錄已知問題與解法（glibc、buildvcs）
- 說明 HTTP-only 測試方式和 provider 切換

**Non-Goals:**
- 不寫 Docker / CI/CD 流程
- 不改任何程式碼邏輯

## Decisions

### 1. 文件放在 README 而非獨立 DEPLOY.md

**理由**: 專案規模小，部署步驟不多，放在 README 一目了然。如果未來內容膨脹再拆分。

### 2. 文件結構

按部署順序組織：
1. 建置（本機 build 或 server 端 build）
2. 靜態編譯（跨平台部署）
3. 複製到 server
4. systemd 設定
5. nginx + HTTPS（可選）
6. 僅 HTTP 測試
7. 常見問題

## Risks / Trade-offs

- 文件可能隨專案演進過時 → 部署步驟簡單，維護成本低
