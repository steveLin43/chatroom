package server

import (
	"net/http"
)

func RegisterHandle() {
	global.inferRootDir()

	// 廣播訊息處理
	go logic.Broadcaster.Start()

	http.HandleFunc("/", homeHandleFunc)
	http.HandleFunc("/ws", WebSocketHandleFunc)
}
