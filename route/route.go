package route

import (
	"github.com/gin-gonic/gin"
	"websocket/internal/controller"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/ws", controller.WebsocketController)
	return r
}
