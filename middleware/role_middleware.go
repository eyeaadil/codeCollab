package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RoleMiddleware restricts access based on user roles
func RoleMiddleware(requiredRole string) gin.HandlerFunc {

	println("dooooooooooooooooooooooooooooooooooooooooooooonnnnnnnn",requiredRole)
	return func(c *gin.Context) {
		// Retrieve user role from context
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: User role not found"})
			c.Abort()
			return
		}


		println("Adiiiiiiiiiiiiiiil",userRole)
		// Check if the user's role matches the required role
		if userRole != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
