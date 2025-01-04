package routes

import (
	"github.com/gin-gonic/gin"
	"codeCollab-backend/controllers"
	"codeCollab-backend/middleware"
)

func SetupProfileRoutes(r *gin.Engine, profileController *controllers.ProfileController) {
	// Grouping the profile routes
	profileGroup := r.Group("/profile")
	profileGroup.Use(middleware.AuthMiddleware()) // Apply the AuthMiddleware to all routes in this group

	{
		// Protected route: Edit profile
		profileGroup.PUT("/edit", profileController.EditProfile)
	}
}
