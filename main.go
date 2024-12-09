package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"codeCollab-backend/config"
	"codeCollab-backend/controllers"
	"codeCollab-backend/routes"
)

func main() {
	// Set up logging
	gin.DefaultWriter = os.Stdout
	log.SetOutput(os.Stderr) // Use stderr for error logging

	// Connect to MongoDB
	config.ConnectDB()

	// Initialize controllers with the connected database
	authController := controllers.NewAuthController(config.DB)
	sessionController := controllers.NewSessionController(config.DB)
	fileController := controllers.NewFileController(config.DB)
	collaboratorController := controllers.NewCollaboratorController(config.DB)


	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	// Initialize Gin router
	router := gin.Default()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Register routes
	routes.RegisterAuthRoutes(router, authController)
	routes.RegisterSessionRoutes(router, sessionController)
	routes.RegisterFileRoutes(router, fileController)
	routes.RegisterCollaboratorRoutes(router, collaboratorController)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	log.Printf("🚀 Server is running at http://localhost:%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
