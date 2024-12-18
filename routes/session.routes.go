package routes

import (
	"codeCollab-backend/controllers"
	"codeCollab-backend/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterSessionRoutes sets up protected routes for session operations
func RegisterSessionRoutes(router *gin.Engine, sessionController *controllers.SessionController) {
	// Apply the AuthMiddleware to protect all /sessions routes
	session := router.Group("/sessions").Use(middleware.AuthMiddleware())
	{
		session.POST("/", sessionController.CreateSession)       // Create a new session
		session.GET("/:id", sessionController.GetSession)        // Get a session by ID
		session.PUT("/:id", sessionController.UpdateSession)     // Update a session by ID
		session.DELETE("/:id", sessionController.DeleteSession)  // Delete a session by ID
	}
}
