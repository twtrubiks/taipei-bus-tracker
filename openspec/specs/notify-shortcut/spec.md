### Requirement: 快捷儲存
互動式選站完成後，系統 SHALL 提示使用者輸入快捷名稱，將選站結果儲存至 `notify.yaml`。

#### Scenario: 儲存新快捷
- **WHEN** 使用者完成互動選站，輸入快捷名稱「上班」
- **THEN** 系統將路線、方向、站點、閾值儲存至 `notify.yaml`，顯示儲存成功訊息

#### Scenario: 跳過儲存
- **WHEN** 使用者完成互動選站，直接按 Enter
- **THEN** 系統不儲存，直接進入監控模式

#### Scenario: 覆蓋同名快捷
- **WHEN** 使用者輸入的名稱已存在
- **THEN** 系統顯示現有設定內容，詢問是否覆蓋（y/N），使用者確認後覆蓋

### Requirement: 快捷啟動
使用者提供快捷名稱作為 CLI 參數時，系統 SHALL 跳過互動流程，直接載入設定並啟動監控。

#### Scenario: 載入快捷成功
- **WHEN** 使用者執行 `bus-notify 上班` 且快捷存在
- **THEN** 系統載入設定，顯示載入的路線和站點資訊，直接進入監控模式

#### Scenario: 快捷不存在
- **WHEN** 使用者執行 `bus-notify 不存在` 且快捷不存在
- **THEN** 系統顯示錯誤訊息，列出可用的快捷名稱（如有），然後退出

#### Scenario: notify.yaml 不存在或解析失敗
- **WHEN** 使用者指定快捷但 notify.yaml 不存在或格式錯誤
- **THEN** 系統顯示錯誤訊息並退出

### Requirement: 列出快捷
使用者執行 `bus-notify --list` 時，系統 SHALL 顯示帶編號的快捷列表，並提示使用者選擇。

#### Scenario: 選擇快捷啟動
- **WHEN** 使用者執行 `bus-notify --list`，列表顯示後輸入數字
- **THEN** 系統載入對應快捷，初始化資料源，進入監控模式

#### Scenario: 取消選擇
- **WHEN** 使用者執行 `bus-notify --list`，列表顯示後按 Enter
- **THEN** 系統結束，不進入監控

#### Scenario: 無效輸入
- **WHEN** 使用者輸入超出範圍的數字或非數字文字
- **THEN** 系統顯示「無效選擇」後結束

#### Scenario: 無快捷
- **WHEN** 無已儲存的快捷
- **THEN** 系統顯示「尚無快捷設定」後結束（不顯示選擇提示）

### Requirement: 刪除快捷
使用者執行 `bus-notify --delete <名稱>` 時，系統 SHALL 刪除指定快捷。

#### Scenario: 刪除成功
- **WHEN** 使用者執行 `bus-notify --delete 上班` 且快捷存在
- **THEN** 系統從 notify.yaml 移除該快捷，顯示刪除成功

#### Scenario: 刪除不存在的快捷
- **WHEN** 使用者執行 `bus-notify --delete 不存在`
- **THEN** 系統顯示「快捷不存在」錯誤

### Requirement: 互動式選站流程
互動式選站流程結束後，新增儲存快捷步驟（可跳過），再進入監控模式。原有互動流程不變。
