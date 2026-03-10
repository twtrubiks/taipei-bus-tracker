## MODIFIED Requirements

### Requirement: 列出快捷（修改）
`bus-notify --list` SHALL 顯示帶編號的快捷列表，並提示使用者選擇。

#### Scenario: 選擇快捷啟動
- **WHEN** 使用者執行 `bus-notify --list`，列表顯示後輸入數字 `2`
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
