package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"codeCollab-backend/models"
	"codeCollab-backend/utils"
)

type ProfileController struct {
	userCollection *mongo.Collection
}

// Constructor for ProfileController
func NewProfileController(db *mongo.Database) *ProfileController {
	return &ProfileController{
		userCollection: db.Collection("users"),
	}
}

type EditProfileRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
}

// EditProfile handles the profile editing logic
func (pc *ProfileController) EditProfile(c *gin.Context) {
	// Extract user ID from the token stored in cookies
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized, missing access token"})
		return
	}

	// Decode the access token to get the user ID
	userID, err := utils.ParseToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Retrieve the user from the database
	var user models.User
	err = pc.userCollection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Bind and validate the incoming edit request
	var req EditProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update user profile
	update := bson.M{
		"$set": bson.M{
			"username":   req.Username,
			"email":      req.Email,
			"updated_at": time.Now(),
		},
	}

	// Perform the update operation
	_, err = pc.userCollection.UpdateOne(context.Background(), bson.M{"_id": user.ID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	// Fetch the updated user data
	err = pc.userCollection.FindOne(context.Background(), bson.M{"_id": user.ID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated user data"})
		return
	}

	// Return the updated user profile
	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user": gin.H{
			"id":        user.ID.Hex(),
			"username":  user.Username,
			"email":     user.Email,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
		},
	})
}
