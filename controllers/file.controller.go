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

type FileController struct {
	fileCollection *mongo.Collection
}

// Constructor for FileController
func NewFileController(db *mongo.Database) *FileController {
	return &FileController{
		fileCollection: db.Collection("files"),
	}
}

// CreateFile adds a new file to a session
func (fc *FileController) CreateFile(c *gin.Context) {
	var file models.File

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

	// Bind the file data from the request body
	if err := c.ShouldBindJSON(&file); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set fields
	file.ID = primitive.NewObjectID()
	file.UserID = objectID // Set the userID extracted from the token
	file.CreatedAt = time.Now()
	file.UpdatedAt = time.Now()
	file.Version = 1

	_, err = fc.fileCollection.InsertOne(context.Background(), file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "File created successfully", "file": file})
}

// GetFiles retrieves all files in a session
func (fc *FileController) GetFiles(c *gin.Context) {
	sessionID := c.Param("session_id")
	objectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	cursor, err := fc.fileCollection.Find(context.Background(), bson.M{"session_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve files"})
		return
	}
	defer cursor.Close(context.Background())

	var files []models.File
	if err = cursor.All(context.Background(), &files); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode files"})
		return
	}

	c.JSON(http.StatusOK, files)
}

// UpdateFile updates a file's content
func (fc *FileController) UpdateFile(c *gin.Context) {
	fileID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	var updates bson.M
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates["updated_at"] = time.Now()
	updates["version"] = bson.M{"$inc": 1}

	_, err = fc.fileCollection.UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{"$set": updates})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File updated successfully"})
}

// DeleteFile removes a file by its ID
func (fc *FileController) DeleteFile(c *gin.Context) {
	fileID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	_, err = fc.fileCollection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}
