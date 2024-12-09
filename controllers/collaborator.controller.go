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
)

type CollaboratorController struct {
	collaboratorCollection *mongo.Collection
}

// Constructor for CollaboratorController
func NewCollaboratorController(db *mongo.Database) *CollaboratorController {
	return &CollaboratorController{
		collaboratorCollection: db.Collection("collaborators"),
	}
}

// AddCollaborator adds a new collaborator to a session
func (cc *CollaboratorController) AddCollaborator(c *gin.Context) {
	var collaborator models.Collaborator
	if err := c.ShouldBindJSON(&collaborator); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set fields
	collaborator.ID = primitive.NewObjectID()
	collaborator.CreatedAt = time.Now()
	collaborator.LastModified = time.Now()

	_, err := cc.collaboratorCollection.InsertOne(context.Background(), collaborator)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add collaborator"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Collaborator added successfully", "collaborator": collaborator})
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
func (cc *CollaboratorController) UpdateCollaboratorRole(c *gin.Context) {
	collaboratorID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(collaboratorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collaborator ID"})
		return
	}

	var update bson.M
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update["last_modified"] = time.Now()

	_, err = cc.collaboratorCollection.UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{"$set": update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update collaborator"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Collaborator role updated successfully"})
}

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
