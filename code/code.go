package code

type Code int

const (
	WebsocketSuccess Code = iota + 50000
	WebsocketEnd
	WebsocketOnlineReply
	WebsocketOfflineReply
	WebsocketLimit
)

func (c Code) Int() int {
	return int(c)
}

const (
	SengMessage             int = 1
	GetHistoryMessage       int = 2
	GetUnreadHistoryMessage int = 3
	GetUnreadMessageCounts  int = 4
)
