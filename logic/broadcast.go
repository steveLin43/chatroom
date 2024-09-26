package logic

type broadcaster struct { // 單例模式
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

func (b *broadcaster) CanEnterRoom(nickname string) bool {
	b.checkUserChannel <- nickname

	return <-b.checkUserCanInChannel
}
