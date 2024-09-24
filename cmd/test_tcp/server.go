package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", ":2020")
	if err != nil {
		panic(err)
	}

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConn(conn)
	}
}

type User struct {
	ID             int
	Addr           string
	EnterAt        time.Time
	MessageChannel chan string
}

// 給使用者發送的訊息
type Message struct {
	OwnerID int
	Content string
}

var (
	// 新用户到来，通过该 channel 进行登记
	enteringChannel = make(chan *User)
	// 用户离开，通过该 channel 进行登记
	leavingChannel = make(chan *User)
	// 广播专用的用户普通消息 channel，缓冲是尽可能避免出现异常情况堵塞
	messageChannel = make(chan Message, 8)
)

// broadcaster 用於紀錄聊天室使用者，並進行訊息廣播
// 1. 新使用者進來；2. 使用者普通訊息；3. 使用者離開
func broadcaster() {
	users := make(map[*User]struct{})

	for {
		select {
		case user := <-enteringChannel:
			// 新使用者進入
			users[user] = struct{}{}
		case user := <-leavingChannel:
			// 使用者離開
			delete(users, user)
			// 避免 goroutine 洩漏
			close(user.MessageChannel)
		case msg := <-messageChannel:
			// 給所有線上使用者發送訊息
			for user := range users {
				if user.ID == msg.OwnerID {
					continue
				}
				user.MessageChannel <- msg.Content
			}
		}
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	// 1. 新使用者進來，建置該使用者實例
	user := &User{
		ID:             GenUserID(),
		Addr:           conn.RemoteAddr().String(),
		EnterAt:        time.Now(),
		MessageChannel: make(chan string, 8),
	}

	// 2. 由於目前是在一個新的 goroutine 中進行讀取操作的，所以需要開一個 goroutine 用於寫入操作
	// 讀寫 goroutine 之間可以透過 channel 進行通訊
	go sendMessage(conn, user.MessageChannel)

	// 3. 給目前使用者發送歡迎資訊，向所有使用者告知新使用者到來
	user.MessageChannel <- "Welcome, " + user.String()
	msg := Message{
		OwnerID: user.ID,
		Content: "user:`" + strconv.Itoa(user.ID) + "` has enter",
	}
	messageChannel <- msg

	// 4. 記錄到全域使用者清單中，避免用鎖
	enteringChannel <- user

	// 踢出超時使用者
	var userActive = make(chan struct{})
	go func() {
		d := 5 * time.Minute
		timer := time.NewTimer(d)
		for {
			select {
			case <-timer.C:
				conn.Close()
			case <-userActive:
				timer.Reset(d)
			}
		}
	}()

	// 5. 循環讀取使用者輸入
	input := bufio.NewScanner(conn)
	for input.Scan() {
		msg.Content = strconv.Itoa(user.ID) + ":" + input.Text()
		messageChannel <- msg

		// 使用者活躍
		userActive <- struct{}{}
	}

	if err := input.Err(); err != nil {
		log.Println("讀取錯誤:", err)
	}

	// 6. 使用者離開
	leavingChannel <- user
	msg.Content = "user:`" + strconv.Itoa(user.ID) + "` has left"
	messageChannel <- msg
}

func sendMessage(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}
