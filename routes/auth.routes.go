package routes

import (
	"codeCollab-backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.Engine, authController *controllers.AuthController) {

	println("register")
	auth := router.Group("/auth")
	{
		// Public routes
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)

		// Protected routes
		auth.POST("/logout", authController.Logout) // AuthMiddleware is applied directly here
	}
}
