package routes

import (
	"codeCollab-backend/controllers"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func ProjectRoutes(router *gin.Engine, db *mongo.Database) {
	// Create a new ProjectController instance
	projectController := controllers.NewProjectController(db)

	// Define routes for projects
	projectRoutes := router.Group("/api/projects")
	{
		projectRoutes.POST("/", projectController.CreateProject)           // Create a new project
		projectRoutes.PUT("/:id", projectController.UpdateProject)         // Update an existing project
		projectRoutes.DELETE("/:id", projectController.DeleteProject)      // Delete a project
		projectRoutes.POST("/:id/collaborators", projectController.AddCollaborator) // Add a collaborator
	}
}
