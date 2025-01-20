package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"codeCollab-backend/models"
	"codeCollab-backend/utils" // Assuming utils has the ParseToken function
)

type FolderController struct {
	folderCollection *mongo.Collection
}

// Constructor for FolderController
func NewFolderController(db *mongo.Database) *FolderController {
	return &FolderController{
		folderCollection: db.Collection("folders"),
	}
}

// CreateFolder adds a new folder to a session
func (fc *FolderController) CreateFolder(c *gin.Context) {
	var folder models.Folder

	// Extract token from cookies
	token, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
		return
	}

	// Parse the token and extract the user ID
	userID, parseErr := utils.ParseToken(token)
	if parseErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Convert userID (string) to primitive.ObjectID if needed
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Bind the folder data from the request body
	if err := c.ShouldBindJSON(&folder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set fields
	folder.ID = primitive.NewObjectID()
	folder.UserID = objectID // Set the userID extracted from the token
	folder.CreatedAt = time.Now()
	folder.UpdatedAt = time.Now()

	_, err = fc.folderCollection.InsertOne(context.Background(), folder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create folder"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Folder created successfully", "folder": folder})
}

// GetFolders retrieves all folders in a session
func (fc *FolderController) GetFolders(c *gin.Context) {
	sessionID := c.Param("session_id")
	objectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	cursor, err := fc.folderCollection.Find(context.Background(), bson.M{"session_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve folders"})
		return
	}
	defer cursor.Close(context.Background())

	var folders []models.Folder
	if err = cursor.All(context.Background(), &folders); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode folders"})
		return
	}

	c.JSON(http.StatusOK, folders)
}

// UpdateFolder updates a folder's details
func (fc *FolderController) UpdateFolder(c *gin.Context) {
	folderID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(folderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder ID"})
		return
	}

	var updates bson.M
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates["updated_at"] = time.Now()

	_, err = fc.folderCollection.UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{"$set": updates})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update folder"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Folder updated successfully"})
}

// DeleteFolder removes a folder by its ID
func (fc *FolderController) DeleteFolder(c *gin.Context) {
	folderID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(folderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder ID"})
		return
	}

	_, err = fc.folderCollection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete folder"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Folder deleted successfully"})
}
