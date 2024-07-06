package service

import (
	"github.com/gorilla/websocket"
	"log"
	"strings"
	"sync"
	"websocket/code"
	"websocket/internal/dao"
	"websocket/model"
)

var wsMutex sync.Mutex

func WebsocketHandle(sengID string, receiveID string, conn *websocket.Conn) {
	client := &model.Client{
		ID:     GenerateId(sengID, receiveID),
		SendID: GenerateId(receiveID, sengID),
		Socket: conn,
		Send:   make(chan []byte),
	}
	model.Manager.Register <- client
	go Read(client)
	go Write(client)
	//go HandleUnreadMessages(conn, sengID)
}

// Read 定义了读取的方法
func Read(c *model.Client) {
	defer func() {
		model.Manager.Unregister <- c
		_ = c.Socket.Close()
	}()
	for {
		c.Socket.PongHandler()
		sendMsg := new(model.SendMsg)
		err := c.Socket.ReadJSON(sendMsg)
		if err != nil {
			model.Manager.Unregister <- c
			_ = c.Socket.Close()
			break
		}
		switch sendMsg.Type {
		case code.SengMessage:
			//r1, _ := redis.Rdb.Get(context.Background(), c.ID).Result()
			//r2, _ := redis.Rdb.Get(context.Background(), c.SendID).Result()
			//if r1 >= "3" && r2 == "0" {
			//	// 1->2发送了超过三次消息，2没有回复,此时禁止1发送消息
			//	replyMessage := &model.ReplyMsg{
			//		From:    "Server",
			//		Code:    code.WebsocketLimit.Int(),
			//		Content: "发送信息过多对方已禁止你发送消息",
			//	}
			//	writeToSocket(c.Socket, replyMessage)
			//	continue
			//} else {
			//	redis.Rdb.Incr(context.Background(), c.ID)
			//	_, _ = redis.Rdb.Expire(context.Background(), c.ID, time.Hour*24*30*3).Result()
			//}
			model.Manager.Broadcast <- &model.Broadcast{
				Client:  c,
				Message: []byte(sendMsg.Content),
			}

		case code.GetHistoryMessage:
			generatedID := c.ID //1->2
			parts := strings.Split(generatedID, "->")
			sendID := parts[1]    // 1
			receiveID := parts[0] // 2

			result, _ := dao.GetHistoryMessage(sendID, receiveID, -1)
			if len(result) == 0 {
				replyMsg := &model.ReplyMsg{
					From:    "Server",
					Code:    code.WebsocketEnd.Int(),
					Content: "没有更多消息了",
				}
				writeToSocket(c.Socket, replyMsg)
				continue
			}
			replyMessage := &model.ReplyMsg{
				From:    sendID,
				Code:    code.WebsocketSuccess.Int(),
				Content: result,
			}
			dao.MarkMessagesAsRead(receiveID)
			writeToSocket(c.Socket, replyMessage)

			//for _, m := range result {
			//	replyMessage := &model.ReplyMsg{
			//		From:    sendID,
			//		Code:    code.WebsocketSuccess.Int(),
			//		Content: m.Content,
			//	}
			//	writeToSocket(c.Socket, replyMessage)
			//}
		case code.GetUnreadHistoryMessage:
			generatedID := c.ID //1->2
			parts := strings.Split(generatedID, "->")
			sendID := parts[1]    // 1
			receiveID := parts[0] // 2
			result, _ := dao.GetUnreadMessages(sendID, receiveID)
			if len(result) == 0 {
				replyMsg := &model.ReplyMsg{
					From:    "Server",
					Code:    code.WebsocketEnd.Int(),
					Content: "没有更多消息了",
				}
				writeToSocket(c.Socket, replyMsg)
				continue
			}
			replyMessage := &model.ReplyMsg{
				From:    sendID,
				Code:    code.WebsocketSuccess.Int(),
				Content: result,
			}
			dao.MarkMessagesAsRead(receiveID)
			writeToSocket(c.Socket, replyMessage)

		case code.GetUnreadMessageCounts:
			generatedID := c.ID //1->2
			parts := strings.Split(generatedID, "->")
			sendID := parts[0] // 1
			HandleUnreadMessages(c.Socket, sendID)
		}
	}
}

// Write 定义了写入的方法
func Write(c *model.Client) {
	defer func() {
		_ = c.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				writeToSocket(c.Socket, websocket.CloseMessage)
				return
			}
			replyMsg := &model.ReplyMsg{
				From:    c.ID,
				Code:    code.WebsocketSuccess.Int(),
				Content: string(message),
			}
			writeToSocket(c.Socket, replyMsg)
		}
	}
}

func Start() {
	for {
		select {
		case conn := <-model.Manager.Register:
			model.Manager.Clients[conn.ID] = conn
			replyMessage := &model.ReplyMsg{
				From:    "Server",
				Code:    code.WebsocketSuccess.Int(),
				Content: "连接成功",
			}
			writeToSocket(conn.Socket, replyMessage)
		case conn := <-model.Manager.Unregister:
			if _, ok := model.Manager.Clients[conn.ID]; ok {
				replyMessage := &model.ReplyMsg{
					From:    "Server",
					Code:    code.WebsocketEnd.Int(),
					Content: "连接断开",
				}
				writeToSocket(conn.Socket, replyMessage)
				close(conn.Send)
				delete(model.Manager.Clients, conn.ID)
			}
		case broadcast := <-model.Manager.Broadcast:
			message := broadcast.Message
			sendID := broadcast.Client.SendID
			flag := false
			for id, conn := range model.Manager.Clients {
				if id != sendID {
					continue
				}
				select {
				case conn.Send <- message:
					flag = true
				default:
					close(conn.Send)
					delete(model.Manager.Clients, conn.ID)
				}
			}
			generatedID := broadcast.Client.ID // 1->2
			parts := strings.Split(generatedID, "->")
			send := parts[0]  // 1
			reply := parts[1] // 2

			var replyMessage *model.ReplyMsg
			if flag {
				replyMessage = &model.ReplyMsg{
					From:    sendID,
					Code:    code.WebsocketOnlineReply.Int(),
					Content: "消息发送成功,对方在线",
				}
			} else {
				replyMessage = &model.ReplyMsg{
					From:    sendID,
					Code:    code.WebsocketOfflineReply.Int(),
					Content: "消息发送成功,对方离线",
				}
			}
			writeToSocket(broadcast.Client.Socket, replyMessage)
			err := dao.SaveMessage(send, reply, string(message), flag)
			if err != nil {
				log.Println("消息保存失败")
			}
		}
	}
}

func writeToSocket(conn *websocket.Conn, v interface{}) {
	if conn == nil {
		log.Println("WebSocket connection is nil, skipping write")
		return
	}

	wsMutex.Lock()
	defer wsMutex.Unlock()

	if err := conn.WriteJSON(v); err != nil {
		log.Println("Error writing to WebSocket:", err)
	}
}

// HandleUnreadMessages 处理未读消息
func HandleUnreadMessages(conn *websocket.Conn, receiveID string) {
	//messages, err := dao.GetUnreadMessages(receiveID)
	//if err != nil {
	//	log.Println("Error getting unread messages:", err)
	//	return
	//}
	//
	//for _, msg := range messages {
	//	replyMsg := &model.ReplyMsg{
	//		From:    msg.SendID,
	//		Code:    code.WebsocketSuccessMessage.Int(),
	//		Content: msg.Content,
	//	}
	//	writeToSocket(conn, replyMsg)
	//}
	//
	//// 标记消息为已读
	//if err := dao.MarkMessagesAsRead(receiveID); err != nil {
	//	log.Println("Error marking messages as read:", err)
	//}

	totalUnreadCount, unreadCount, err := dao.GetUnreadMessageCounts(receiveID)
	if err != nil {
		log.Println("Error getting unread messages:", err)
		return
	}
	replyMessage := &model.UnreadMessage{
		TotalCounts:  totalUnreadCount,
		UnreadCounts: unreadCount,
	}
	writeToSocket(conn, replyMessage)

	// dao.MarkMessagesAsRead(receiveID)
	return
}
