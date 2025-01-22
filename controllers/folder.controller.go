// controllers/folderController.go
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
	userCollection   *mongo.Collection
}

// Constructor for FolderController
func NewFolderController(db *mongo.Database) *FolderController {
	return &FolderController{
		folderCollection: db.Collection("folders"),
		userCollection:   db.Collection("users"),
	}
}

func (fc *FolderController) CreateFolder(c *gin.Context) {
	var folder models.Folder

	token, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	userID, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := c.ShouldBindJSON(&folder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	folderID := primitive.NewObjectID()
	newFolder := models.Folder{
		ID:        folderID,
		UserID:    userObjectID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      folder.Name,
	}

	// Initialize folder_ids array if it doesn't exist
	_, err = fc.userCollection.UpdateOne(
		context.Background(),
		bson.M{
			"_id":        userObjectID,
			"folder_ids": bson.M{"$exists": false},
		},
		bson.M{
			"$set": bson.M{"folder_ids": []primitive.ObjectID{}},
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize folder_ids"})
		return
	}

	// Insert folder and update user atomically
	session, err := fc.folderCollection.Database().Client().StartSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start session"})
		return
	}
	defer session.EndSession(context.Background())

	err = session.StartTransaction()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	err = mongo.WithSession(context.Background(), session, func(sc mongo.SessionContext) error {
		_, err := fc.folderCollection.InsertOne(sc, newFolder)
		if err != nil {
			return err
		}

		_, err = fc.userCollection.UpdateOne(
			sc,
			bson.M{"_id": userObjectID},
			bson.M{"$addToSet": bson.M{"folder_ids": folderID}},
		)
		return err
	})

	if err != nil {
		session.AbortTransaction(context.Background())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
		return
	}

	if err = session.CommitTransaction(context.Background()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Return updated user along with the folder
	var updatedUser models.User
	if err := fc.userCollection.FindOne(context.Background(), bson.M{"_id": userObjectID}).Decode(&updatedUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Folder created successfully",
		"folder":  newFolder,
		"user":    updatedUser,
	})
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

	updates["updated_at"] = time.Now().Unix()

	_, err = fc.folderCollection.UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{"$set": updates})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update folder"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Folder updated successfully"})
}

// DeleteFolder removes a folder by its ID and also updates the User model
func (fc *FolderController) DeleteFolder(c *gin.Context) {
	folderID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(folderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder ID"})
		return
	}

	// Find the folder to delete
	var folder models.Folder
	err = fc.folderCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&folder)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Folder not found"})
		return
	}

	// Remove the folder from the User model
	_, err = fc.userCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": folder.UserID},                    // Find the user associated with the folder
		bson.M{"$pull": bson.M{"folder_ids": objectID}}, // Remove the folder ID from the folder_ids array
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove folder ID from user"})
		return
	}

	// Delete the folder from the folders collection
	_, err = fc.folderCollection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete folder"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Folder deleted successfully"})
}
