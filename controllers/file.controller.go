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
	"codeCollab-backend/utils"
)

type FileController struct {
	fileCollection       *mongo.Collection
	fileVersionCollection *mongo.Collection
}

// Constructor for FileController
func NewFileController(db *mongo.Database) *FileController {
	return &FileController{
		fileCollection:       db.Collection("files"),
		fileVersionCollection: db.Collection("file_versions"),
	}
}

// CreateFile adds a new file to a folder
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

	// Convert userID to primitive.ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Bind the file data from the request body
	if err := c.ShouldBindJSON(&file); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate FolderID
	if file.FolderID.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Folder ID is required"})
		return
	}

	// Set fields
	file.ID = primitive.NewObjectID()
	file.UserID = userObjectID
	file.CreatedAt = time.Now()
	file.UpdatedAt = time.Now()
	file.Version = 1
	file.LastEditedBy = userObjectID

	_, err = fc.fileCollection.InsertOne(context.Background(), file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
		return
	}

	// Save initial version
	version := models.FileVersion{
		ID:        primitive.NewObjectID(),
		FileID:    file.ID,
		Content:   file.Content,
		Version:   file.Version,
		EditedBy:  userObjectID,
		EditedAt:  file.CreatedAt,
	}
	_, err = fc.fileVersionCollection.InsertOne(context.Background(), version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file version"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "File created successfully", "file": file})
}

// GetFiles retrieves all files within a folder
func (fc *FileController) GetFiles(c *gin.Context) {
	folderID := c.Param("folder_id")
	objectID, err := primitive.ObjectIDFromHex(folderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder ID"})
		return
	}

	cursor, err := fc.fileCollection.Find(context.Background(), bson.M{"folder_id": objectID})
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

// UpdateFile updates a file's content and tracks versions
func (fc *FileController) UpdateFile(c *gin.Context) {
	fileID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	var updates struct {
		Content  string `json:"content" binding:"required"`
		Language string `json:"language"`
	}

	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract token from cookies to identify the editor
	token, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
		return
	}

	userID, parseErr := utils.ParseToken(token)
	if parseErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	editorID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Retrieve current file
	var currentFile models.File
	err = fc.fileCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&currentFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File not found"})
		return
	}

	// Increment version and update fields
	newVersion := currentFile.Version + 1
	updatesMap := bson.M{
		"content":       updates.Content,
		"language":      updates.Language,
		"version":       newVersion,
		"updated_at":    time.Now(),
		"last_edited_by": editorID,
	}

	_, err = fc.fileCollection.UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{"$set": updatesMap})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update file"})
		return
	}

	// Save version
	version := models.FileVersion{
		ID:        primitive.NewObjectID(),
		FileID:    objectID,
		Content:   updates.Content,
		Version:   newVersion,
		EditedBy:  editorID,
		EditedAt:  time.Now(),
	}
	_, err = fc.fileVersionCollection.InsertOne(context.Background(), version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file version"})
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

	_, err = fc.fileVersionCollection.DeleteMany(context.Background(), bson.M{"file_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file versions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}
