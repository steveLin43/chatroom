package logic

type broadcaster struct { // 單例模式:一種類別僅有一個實例，有積極模式與懶散模式兩種。
	// 所有聊天室使用者
	users map[string]*User

	// 所有 channel 統一管理，可以避免外部亂用

	enteringChannel chan *User
	leavingChannel  chan *User
	messageChannel  chan *Message

	// 判斷該暱稱使用者是否可進入聊天室(重複與否)
	checkUserChannel      chan string
	checkUserCanInChannel chan bool
}

var Broadcaster = &broadcaster{ // Go 較推薦積極模式，此處也用積極模式實現。
	users: make(map[string]*User),

	enteringChannel: make(chan *User),
	leavingChannel:  make(chan *User),
	messageChannel:  make(chan *Message, MessageQueueLen),

	checkUserChannel:      make(chan string),
	checkUserCanInChannel: make(chan bool),
}

// Start 啟動廣播器
// 需要在一個 goroutine 中執行，因為它不會傳回
// select-case 架構是專門為 channel 設計的
func (b *broadcaster) Start() {
	for {
		select {
		case user := <-b.enteringChannel:
			// 新使用者進入
			b.users[user.NickName] = user

			b.sendUserList()
			OfflineProcessor.Send(user)
		case user := <-b.leavingChannel:
			// 使用者離開
			delete(b.users, user.NickName)
			// 避免 goroutine 洩漏
			user.CloseMessageChannel()

			b.sendUserList()
		case msg := <-b.messageChannel:
			// 給所有線上使用者發送訊息
			for _, user := range b.users {
				if user.UID == msg.User.UID {
					continue
				}
				user.MessageChannel <- msg
			}
			OfflineProcessor.Save(msg)
		case nickname := <-b.checkUserChannel:
			if _, ok := b.users[nickname]; ok {
				b.checkUserCanInChannel <- false
			} else {
				b.checkUserCanInChannel <- true
			}
		}

	}
}

func (b *broadcaster) UserEntering(u *User) {
	b.enteringChannel <- u
}

func (b *broadcaster) UserLeaving(u *User) {
	b.leavingChannel <- u
}

func (b *broadcaster) Broadcast(msg *Message) {
	b.messageChannel <- msg
}

func (b *broadcaster) CanEnterRoom(nickname string) bool {
	b.checkUserChannel <- nickname

	return <-b.checkUserCanInChannel
}
