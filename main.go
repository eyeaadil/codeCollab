// filepath: /path/to/your/go/server/main.go
package main

import (
    "log"
    "os"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
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
    // projectController := controllers.NewProjectController(config.DB) // Initialize ProjectController

    // Set Gin to release mode for production
    gin.SetMode(gin.ReleaseMode)

    // Initialize Gin router
    router := gin.Default()

    // Middleware
    router.Use(gin.Logger())
    router.Use(gin.Recovery())

    // CORS middleware
    router.Use(cors.Default())

    // Register routes
    routes.RegisterAuthRoutes(router, authController)
    routes.RegisterSessionRoutes(router, sessionController)
    routes.RegisterFileRoutes(router, fileController)
    routes.RegisterCollaboratorRoutes(router, collaboratorController)
    // routes.ProjectRoutes(router, projectController) // Add project routes

    // Register WebSocket routes
    routes.RegisterWebSocketRoutes(router)

    // Define a /data endpoint
    router.GET("/data", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "items": []map[string]interface{}{
                {"id": 1, "name": "Item 1"},
                {"id": 2, "name": "Item 2"},
            },
        })
    })

    // Start the server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8000" // Default port
    }
    log.Printf("🚀 Server is running at http://localhost:%s", port)
    if err := router.Run(":" + port); err != nil {
        log.Fatalf("❌ Failed to start server: %v", err)
    }
}