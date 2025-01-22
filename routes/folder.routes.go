package routes

import (
	"codeCollab-backend/controllers"
	"codeCollab-backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterFolderRoutes(router *gin.Engine, folderController *controllers.FolderController) {
	folder := router.Group("/folders")
	
	// Apply authentication middleware to all folder routes
	folder.Use(middleware.AuthMiddleware()) 
	{
		folder.POST("/", folderController.CreateFolder)               // Create a new folder
		folder.GET("/:user_id", folderController.GetFolders)       // Get all folders in a session
		folder.PUT("/:id", folderController.UpdateFolder)             // Update a folder by ID
		folder.DELETE("/:id", folderController.DeleteFolder)          // Delete a folder by ID
	}
}
