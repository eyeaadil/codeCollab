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

type SessionController struct {
	sessionCollection *mongo.Collection
}

// Constructor for SessionController
func NewSessionController(db *mongo.Database) *SessionController {
	return &SessionController{
		sessionCollection: db.Collection("sessions"),
	}
}

// CreateSession creates a new coding session
func (sc *SessionController) CreateSession(c *gin.Context) {
	var session models.Session
	if err := c.ShouldBindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return							
	}

	session.ID = primitive.NewObjectID()
	session.CreatedAt = time.Now()
	session.LastActiveAt = time.Now()

	_, err := sc.sessionCollection.InsertOne(context.Background(), session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Session created successfully", "session": session})
}

// GetSession retrieves a session by its ID
func (sc *SessionController) GetSession(c *gin.Context) {
	sessionID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	var session models.Session
	err = sc.sessionCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&session)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	c.JSON(http.StatusOK, session)
}

// UpdateSession updates a session's details
func (sc *SessionController) UpdateSession(c *gin.Context) {
	sessionID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	var updates bson.M
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates["last_active_at"] = time.Now()

	_, err = sc.sessionCollection.UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{"$set": updates})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session updated successfully"})
}

// DeleteSession removes a session by its ID
func (sc *SessionController) DeleteSession(c *gin.Context) {
	sessionID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	_, err = sc.sessionCollection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session deleted successfully"})
}
