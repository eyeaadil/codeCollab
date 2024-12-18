package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)


// Secret key used to sign tokens (replace with your secret key)
var jwtSecret = []byte("123456")

// AuthMiddleware validates the access token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the token from cookies
		tokenString, err := c.Cookie("access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No access token provided"})
			c.Abort()
			return
		}

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the algorithm
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return jwtSecret, nil
		})
	
		// if err != nil {
		// 	return  err
		// }

		println("token",token)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid token", "details": err.Error()})
			c.Abort()
			return
		}

		// Extract claims (user info) and set in context
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID, idExists := claims["sub"].(string)
			// userRole, roleExists := claims["role"].(string)

			if !idExists  {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid claims structure"})
				c.Abort()
				return
			}

			// Set user information in context
			c.Set("user_id", userID)
			// c.Set("user_role", userRole)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid claims"})
			c.Abort()
			return
		}

		// Proceed to the next handler
		c.Next()
	}
}
