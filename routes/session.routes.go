package routes

import (
	"codeCollab-backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterSessionRoutes(router *gin.Engine, sessionController *controllers.SessionController) {
	session := router.Group("/sessions")
	{
		session.POST("/", sessionController.CreateSession)          // Create a new session
		session.GET("/:id", sessionController.GetSession)           // Get a session by ID
		session.PUT("/:id", sessionController.UpdateSession)        // Update a session by ID
		session.DELETE("/:id", sessionController.DeleteSession)     // Delete a session by ID
	}
}
