package routes

import (
	"codeCollab-backend/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterWebSocketRoutes defines WebSocket-related routes
func RegisterWebSocketRoutes(router *gin.Engine) {
	router.GET("/ws", controllers.WebSocketHandler)
}
