## ADDED Requirements

### Requirement: 路線搜尋
使用者 SHALL 能在搜尋框輸入路線名稱（如 "1"、"299"），系統即時顯示匹配的路線列表供選擇。

#### Scenario: 輸入關鍵字搜尋
- **WHEN** 使用者在搜尋框輸入 "1"
- **THEN** 系統顯示所有名稱含 "1" 的路線列表（如 1、12、21、100...），可點擊進入

#### Scenario: 選擇路線後顯示方向
- **WHEN** 使用者點擊路線 "1"
- **THEN** 系統顯示兩個方向選項（如 "往松仁路" / "往萬華"），使用者選擇後進入站點列表

### Requirement: 站點到站時間顯示
選擇路線和方向後，系統 SHALL 顯示該方向所有站點的即時到站資訊。

#### Scenario: 顯示所有站的 ETA
- **WHEN** 使用者選擇路線 1、方向 0
- **THEN** 系統列出所有站點，每站顯示站名、到站狀態（約 X 分 / 進站中 / 未發車等）、車牌號碼

#### Scenario: 指定站高亮
- **WHEN** 使用者點擊某個站
- **THEN** 該站高亮顯示，資訊放大呈現

### Requirement: 自動更新
前端 SHALL 每 15 秒自動向後端請求最新的到站資料並更新畫面。

#### Scenario: 自動輪詢
- **WHEN** 使用者停留在站點列表頁面
- **THEN** 系統每 15 秒自動更新所有站的到站資訊，無需手動重新整理

#### Scenario: 頁面不可見時停止輪詢
- **WHEN** 使用者切換到其他分頁或鎖屏
- **THEN** 系統 SHALL 停止輪詢以節省資源，回到頁面時立即恢復

### Requirement: 收藏功能
使用者 SHALL 能收藏常用的「路線 + 方向 + 站」組合，收藏資料存在 localStorage。

#### Scenario: 新增收藏
- **WHEN** 使用者在站點頁面點擊收藏按鈕
- **THEN** 該「路線 + 方向 + 站」組合被儲存到 localStorage

#### Scenario: 首頁顯示收藏
- **WHEN** 使用者開啟首頁
- **THEN** 系統顯示所有收藏站的即時到站資訊，每個收藏項顯示路線名、站名、到站狀態

#### Scenario: 移除收藏
- **WHEN** 使用者對收藏項目點擊移除
- **THEN** 該收藏從 localStorage 和畫面中移除

### Requirement: 深色模式
系統 SHALL 支援淺色和深色兩種主題，預設跟隨系統設定，也可手動切換。

#### Scenario: 跟隨系統主題
- **WHEN** 使用者裝置設定為深色模式
- **THEN** PWA 自動使用深色主題

#### Scenario: 手動切換
- **WHEN** 使用者點擊主題切換按鈕
- **THEN** 主題在淺色/深色之間切換，選擇 SHALL 保存到 localStorage

### Requirement: 到站通知
使用者 SHALL 能對某站設定到站提醒，當 ETA 到達設定分鐘數時觸發瀏覽器通知。

#### Scenario: 設定到站提醒
- **WHEN** 使用者對某站設定 "3 分鐘前提醒"
- **THEN** 當該站 ETA ≤ 180 秒時，系統觸發瀏覽器 Notification

#### Scenario: 通知權限
- **WHEN** 使用者首次設定提醒
- **THEN** 系統 SHALL 請求瀏覽器通知權限，被拒絕時顯示提示訊息

### Requirement: PWA 安裝體驗
系統 SHALL 提供 manifest.json 和 Service Worker，支援「加到主畫面」安裝。

#### Scenario: 手機加到主畫面
- **WHEN** 使用者在手機瀏覽器選擇「加到主畫面」
- **THEN** 主畫面出現 App icon，點開後全螢幕開啟，無瀏覽器網址列

### Requirement: Responsive 設計
前端 SHALL 支援桌面和手機兩種螢幕尺寸，自適應排版。

#### Scenario: 桌面瀏覽器
- **WHEN** 在寬度 > 768px 的螢幕開啟
- **THEN** 使用寬版排版，可同時顯示更多站點資訊

#### Scenario: 手機瀏覽器
- **WHEN** 在寬度 ≤ 768px 的螢幕開啟
- **THEN** 使用窄版排版，單欄顯示，觸控友善
