package controllers

import (
	"context"
	// "errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"codeCollab-backend/models"
	"codeCollab-backend/utils"
)

type AuthController struct {
	userCollection *mongo.Collection
}

// Constructor for AuthController
func NewAuthController(db *mongo.Database) *AuthController {
	return &AuthController{
		userCollection: db.Collection("users"),
	}
}

type AuthRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Register handles user registration
func (ac *AuthController) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User
	err := ac.userCollection.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	newUser := models.User{
		ID:           primitive.NewObjectID(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		// Role:         models.RoleUser,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err = ac.userCollection.InsertOne(context.Background(), newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	accessToken, refreshToken, err := utils.GenerateTokens(&newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	utils.SetAuthCookies(c, accessToken, refreshToken)

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": models.UserPublicProfile{
		ID:       newUser.ID,
		Username: newUser.Username,
		// Role:     newUser.Role,
		
	}})
}

// Login handles user login
func (ac *AuthController) Login(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err := ac.userCollection.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&user)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	update := bson.M{"$set": bson.M{"last_login_at": time.Now()}}
	_, err = ac.userCollection.UpdateOne(context.Background(), bson.M{"_id": user.ID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update login time"})
		return
	}

	accessToken, refreshToken, err := utils.GenerateTokens(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	utils.SetAuthCookies(c, accessToken, refreshToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "User Login successfully",
		"user": models.UserPublicProfile{
		ID:       user.ID,
		Username: user.Username,
		// Role:     user.Role,
	}})
}





// Logout handles user logout
func (ac *AuthController) Logout(c *gin.Context) {
	// Clear access token cookie
	c.SetCookie("access_token", "", -1, "/", "", true, true)

	// Clear refresh token cookie
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)

	// Return response
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
