package routes

import (
	"codeCollab-backend/controllers"
	"codeCollab-backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterFileRoutes(router *gin.Engine, fileController *controllers.FileController) {
	file := router.Group("/files")
	
	// Apply authentication middleware to all file routes
	file.Use(middleware.AuthMiddleware()) 
	{
		file.POST("/", fileController.CreateFile)                    // Create a new file
		file.GET("/:folder_id", fileController.GetFiles)            // Get all files in a folder
		file.PUT("/:id", fileController.UpdateFile)                  // Update a file by ID
		file.DELETE("/:id", fileController.DeleteFile)               // Delete a file by ID
	}
}
