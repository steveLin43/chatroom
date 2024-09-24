# chatroom
Golang project for chatroom

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
