package utils

import (
	// "os"
	"time"
	"errors"

	// "github.com/golang-jwt/jwt/v4"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"codeCollab-backend/models"
)

// GenerateTokens creates access and refresh tokens
func GenerateTokens(user *models.User) (string, string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID.Hex(),
		// "role": user.Role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID.Hex(),
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	accessTokenString, err := accessToken.SignedString([]byte("123456"))
	if err != nil {
		return "", "", err
	}

	refreshTokenString, err := refreshToken.SignedString([]byte("123456"))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

// SetAuthCookies sets secure cookies for access and refresh tokens
func SetAuthCookies(c *gin.Context, accessToken, refreshToken string) {
	c.SetCookie("access_token", accessToken, 86400, "/", "", true, true)
	c.SetCookie("refresh_token", refreshToken, 604800, "/", "", true, true)
}


// Secret key used to sign tokens (replace with your secret key)
var jwtSecret = []byte("123456")

// ParseToken validates a JWT and extracts the user ID
func ParseToken(tokenString string) (string, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return "", err
	}

	// Extract the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Assuming the user ID is stored in the "sub" field
		if userID, ok := claims["sub"].(string); ok {
			return userID, nil
		}
		return "", errors.New("user ID not found in token")
	}

	return "", errors.New("invalid token")
}
