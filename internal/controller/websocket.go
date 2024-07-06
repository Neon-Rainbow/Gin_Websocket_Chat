package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"websocket/internal/service"
)

func WebsocketController(c *gin.Context) {
	sendID := c.Query("send_id")
	receiveID := c.Query("receive_id")
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { // CheckOrigin解决跨域问题
			return true
		}}).Upgrade(c.Writer, c.Request, nil) // 升级成ws协议
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	service.WebsocketHandle(sendID, receiveID, conn)
	return
}
