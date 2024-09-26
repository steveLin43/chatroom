package server

import (
	"chatroom/logic"
	"log"
	"net/http"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func WebSocketHandleFunc(w http.ResponseWriter, req *http.Request) {
	// Accept 從用戶端接收 WebSocket 驗證，並將連接升級到 WebSocket
	// 如果 Origin 域跟主機不同，Accept 將拒絕驗證，除非設定了 InsecureSkipVerify 選項(透過第三個參數 AcceptOptions 來設定)
	// 換句話說，預設情況下，它不允許跨源請求。如果發生錯誤， Accept 將始終寫入適當的回應。
	conn, err := websocket.Accept(w, req, nil)
	if err != nil {
		log.Println("websocket accept error:", err)
		return
	}

	// 1. 新使用者近來，建置該使用者實例
	nickname := req.FormValue("nickname")
	if l := len(nickname); l < 2 || l > 20 {
		log.Println("nickname illegal: ", nickname)
		wsjson.Write(req.Context(), conn, logic.NewErrorMessage("非法暱稱，暱稱長度：2-20"))
		conn.Close(websocket.StatusUnsupportedData, "nickname illegal!")
		return
	}
	if !logic.Broadcaster.CanEnterRoom(nickname) {
		log.Println("暱稱已經存在：", nickname)
		wsjson.Write(req.Context(), conn, logic.NewErrorMessage("該暱稱已經存在！"))
		conn.Close(websocket.StatusUnsupportedData, "nickname exists!")
		return
	}

	user := logic.NewUser(conn, nickname, req.RemoteAddr)
	// 2. 開啟給使用者發送訊息的 goroutine
	go user.SendMessage(req.Context())

	// 3. 給新使用者發送歡迎資訊
	user.MessageChannel <- logic.NewWelcomeMessage(nickname)

	// 向所有使用者告知新使用者到來
	msg := logic.NewNoticeMessage(nickname + "加入了聊天室")
	logic.Broadcaster.Broadcast(msg)

	// 4. 將該使用者加入廣播器的使用者列表中
	logic.Broadcaster.UserEntering(user)
	log.Println("user:", nickname, "joins chat")

	// 5. 接收使用者訊息
	err = user.ReceiveMessage(req.Context())

	// 6. 使用者離開
	logic.Broadcaster.UserLeaving(user)
	msg = logic.NewNoticeMessage(user.NickName + "離開了聊天室")
	logic.Broadcaster.Broadcast(msg)
	log.Println("user:", nickname, "leaves chat")

	// 根據讀取時的錯誤執行不同的 Close
	if err != nil {
		conn.Close(websocket.StatusNormalClosure, "")
	} else {
		log.Println("read from client error:", err)
		conn.Close(websocket.StatusInternalError, "Read from client error")
	}
}
