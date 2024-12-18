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

type CollaboratorController struct {
	collaboratorCollection *mongo.Collection
	sessionCollection     *mongo.Collection  // Add this line

}

// Constructor for CollaboratorController
func NewCollaboratorController(db *mongo.Database) *CollaboratorController {
	return &CollaboratorController{
		collaboratorCollection: db.Collection("collaborators"),
		sessionCollection:      db.Collection("sessions"),
	}
}
func (cc *CollaboratorController) AddCollaborator(c *gin.Context) {
	var collaborator models.Collaborator

	// Extract token from cookies
	token, err := c.Cookie("access_token")

	println("tooooooooooken", token)
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

	// Convert userID (string) to primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	if err := c.ShouldBindJSON(&collaborator); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate session ID is provided
	if collaborator.SessionID == primitive.NilObjectID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	// Set fields
	collaborator.ID = primitive.NewObjectID()
	collaborator.UserID = objectID // Assign the converted ObjectID
	collaborator.CreatedAt = time.Now()
	collaborator.LastModified = time.Now()

	// Start a database session for transaction
	session, err := cc.collaboratorCollection.Database().Client().StartSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start database session"})
		return
	}
	defer session.EndSession(context.Background())

	// Perform transactional operation
	err = mongo.WithSession(context.Background(), session, func(sessionContext mongo.SessionContext) error {
		// Insert collaborator
		_, err := cc.collaboratorCollection.InsertOne(sessionContext, collaborator)
		if err != nil {
			return err
		}

		// Update session to add collaborator ID
		update := bson.M{
			"$addToSet": bson.M{"collaborators": collaborator.ID},
			"$set":      bson.M{"last_active_at": time.Now()},
		}
		
		_, err = cc.sessionCollection.UpdateOne(
			sessionContext, 
			bson.M{"_id": collaborator.SessionID}, 
			update,
		)
		return err
	})

	// Check transaction result
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add collaborator to session"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Collaborator added successfully",
		"collaborator": collaborator,
	})
}


// GetCollaborators retrieves all collaborators in a session
func (cc *CollaboratorController) GetCollaborators(c *gin.Context) {
	sessionID := c.Param("session_id")
	objectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	cursor, err := cc.collaboratorCollection.Find(context.Background(), bson.M{"session_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve collaborators"})
		return
	}
	defer cursor.Close(context.Background())

	var collaborators []models.Collaborator
	if err := cursor.All(context.Background(), &collaborators); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode collaborators"})
		return
	}

	c.JSON(http.StatusOK, collaborators)
}

// UpdateCollaboratorRole updates a collaborator's role
// func (cc *CollaboratorController) UpdateCollaboratorRole(c *gin.Context) {
// 	collaboratorID := c.Param("id")
// 	objectID, err := primitive.ObjectIDFromHex(collaboratorID)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collaborator ID"})
// 		return
// 	}

// 	var update bson.M
// 	if err := c.ShouldBindJSON(&update); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	update["last_modified"] = time.Now()

// 	_, err = cc.collaboratorCollection.UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{"$set": update})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update collaborator"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Collaborator role updated successfully"})
// }

// RemoveCollaborator removes a collaborator from a session
func (cc *CollaboratorController) RemoveCollaborator(c *gin.Context) {
	collaboratorID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(collaboratorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collaborator ID"})
		return
	}

	_, err = cc.collaboratorCollection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove collaborator"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Collaborator removed successfully"})
}
