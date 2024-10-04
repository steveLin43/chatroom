# chatroom
參考資料：《用 Go 語言完成 6 個大型專案》第四章節 + https://github.com/go-programming-tour-book/chatroom

### 安裝 websocket 套件
```
go get github.com/coder/websocket
```

### 啟動(一個)服務端與(多個)用戶端測試
```
go run cmd/tcp/server.go
go run cmd/tcp/client.go
go run cmd/websocket/server.go
go run cmd/websocket/client.go
```
http://localhost:2021
