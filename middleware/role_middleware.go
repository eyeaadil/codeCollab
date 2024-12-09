package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RoleMiddleware restricts access based on user roles
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve user role from context
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: User role not found"})
			c.Abort()
			return
		}

		// Check if the user's role matches the required role
		if userRole != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
