package routes

import (
	"codeCollab-backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterFileRoutes(router *gin.Engine, fileController *controllers.FileController) {
	file := router.Group("/files")
	{
		file.POST("/", fileController.CreateFile)                    // Create a new file
		file.GET("/:session_id", fileController.GetFiles)            // Get all files in a session
		file.PUT("/:id", fileController.UpdateFile)                  // Update a file by ID
		file.DELETE("/:id", fileController.DeleteFile)               // Delete a file by ID
	}
}
