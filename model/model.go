package model

import (
	"github.com/gorilla/websocket"
)

// SendMsg 定义了发送消息的结构体
type SendMsg struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
}

// ReplyMsg 定义了回复消息的结构体
type ReplyMsg struct {
	From    string      `json:"from"`
	Code    int         `json:"code"`
	Content interface{} `json:"content"`
}

// Client 定义了客户端的结构体
type Client struct {
	ID     string
	SendID string
	Socket *websocket.Conn
	Send   chan []byte
}

// Broadcast 定义了广播的结构体
type Broadcast struct {
	Client  *Client
	Message []byte
	Type    int
}

// ClientManager 定义了客户端管理的结构体
type ClientManager struct {
	Clients    map[string]*Client
	Broadcast  chan *Broadcast
	Reply      chan *Client
	Register   chan *Client
	Unregister chan *Client
}

// Message 定义了消息的结构体
type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

var Manager = ClientManager{
	Clients:    make(map[string]*Client), // 参与连接的用户，出于性能的考虑，需要设置最大连接数
	Broadcast:  make(chan *Broadcast),
	Register:   make(chan *Client),
	Reply:      make(chan *Client),
	Unregister: make(chan *Client),
}
