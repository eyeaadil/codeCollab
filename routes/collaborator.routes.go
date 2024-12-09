package routes

import (
	"codeCollab-backend/controllers"
	"codeCollab-backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterCollaboratorRoutes(router *gin.Engine, collaboratorController *controllers.CollaboratorController) {
	collaborator := router.Group("/collaborators")
	{
		// Apply authentication middleware
		collaborator.Use(middleware.AuthMiddleware())
		{
			// Add collaborator - only admin/host can add
			collaborator.POST("/", middleware.RoleMiddleware("admin"), collaboratorController.AddCollaborator)

			// Get collaborators - authenticated users only
			collaborator.GET("/:session_id", collaboratorController.GetCollaborators)

			// Update collaborator role - only admin can update
			collaborator.PUT("/:id", middleware.RoleMiddleware("admin"), collaboratorController.UpdateCollaboratorRole)

			// Remove collaborator - only admin/host can remove
			collaborator.DELETE("/:id", middleware.RoleMiddleware("admin"), collaboratorController.RemoveCollaborator)
		}
	}
}