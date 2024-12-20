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

type ProjectController struct {
	projectCollection *mongo.Collection
	userCollection   *mongo.Collection
}

func NewProjectController(db *mongo.Database) *ProjectController {
	return &ProjectController{
		projectCollection: db.Collection("projects"),
		userCollection:    db.Collection("users"),
	}
}

// CreateProject handles the creation of a new project
func (pc *ProjectController) CreateProject(c *gin.Context) {
	// Extract userId from cookie
	userID, err := c.Cookie("userId")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in cookie"})
		return
	}

	// Convert userID (string) to primitive.ObjectID
	ownerObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Parse project details from request body
	var project models.ProjectModel
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set project details
	project.ProjectID = primitive.NewObjectID()
	project.OwnerID = ownerObjectID
	project.CreatedAt = time.Now()
	project.UpdatedAt = time.Now()

	// Set default visibility if not provided
	if project.Visibility == "" {
		project.Visibility = "private"
	}

	// Initialize empty slices if not provided
	if project.Collaborators == nil {
		project.Collaborators = []primitive.ObjectID{}
	}
	if project.Files == nil {
		project.Files = []primitive.ObjectID{}
	}

	// Insert project into databaseOOO
	_, err = pc.projectCollection.InsertOne(context.Background(), project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Project created successfully",
		"project": project,
	})
}

// UpdateProject updates an existing project
func (pc *ProjectController) UpdateProject(c *gin.Context) {
	// Extract userId from cookie
	userID, err := c.Cookie("userId")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in cookie"})
		return
	}

	// Convert userID (string) to primitive.ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Get project ID from URL parameter
	projectID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Check if project exists and user is the owner
	var existingProject models.ProjectModel
	err = pc.projectCollection.FindOne(context.Background(), bson.M{
		"_id":      objectID,
		"owner_id": userObjectID,
	}).Decode(&existingProject)
	
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Project not found or you're not authorized to update"})
		return
	}

	// Parse update data
	var updates bson.M
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Always update the UpdatedAt timestamp
	updates["updated_at"] = time.Now()

	// Prevent changing OwnerID
	delete(updates, "owner_id")
	delete(updates, "_id")

	// Perform update
	_, err = pc.projectCollection.UpdateOne(
		context.Background(), 
		bson.M{"_id": objectID}, 
		bson.M{"$set": updates},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project updated successfully"})
}

// DeleteProject removes a project
func (pc *ProjectController) DeleteProject(c *gin.Context) {
	// Extract userId from cookie
	userID, err := c.Cookie("userId")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in cookie"})
		return
	}

	// Convert userID (string) to primitive.ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Get project ID from URL parameter
	projectID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Delete project (only if user is the owner)
	result, err := pc.projectCollection.DeleteOne(context.Background(), bson.M{
		"_id":      objectID,
		"owner_id": userObjectID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Project not found or you're not authorized to delete"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

// AddCollaborator adds a collaborator to a project
func (pc *ProjectController) AddCollaborator(c *gin.Context) {
	// Extract userId from cookie (owner who is adding the collaborator)
	userID, err := c.Cookie("userId")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in cookie"})
		return
	}

	// Convert userID (string) to primitive.ObjectID
	ownerObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid owner ID format"})
		return
	}

	// Get project ID from URL parameter
	projectID := c.Param("id")
	projectObjectID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Parse collaborator details
	var collaboratorInput struct {
		CollaboratorID string `json:"collaborator_id"`
	}
	if err := c.ShouldBindJSON(&collaboratorInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert collaborator ID
	collaboratorObjectID, err := primitive.ObjectIDFromHex(collaboratorInput.CollaboratorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collaborator ID format"})
		return
	}

	// Start a transaction
	session, err := pc.projectCollection.Database().Client().StartSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start database session"})
		return
	}
	defer session.EndSession(context.Background())

	// Perform transactional operation
	err = mongo.WithSession(context.Background(), session, func(sessionContext mongo.SessionContext) error {
		// Verify project exists and user is the owner
		filter := bson.M{
			"_id":      projectObjectID,
			"owner_id": ownerObjectID,
		}
		var project models.ProjectModel
		err := pc.projectCollection.FindOne(sessionContext, filter).Decode(&project)
		if err != nil {
			return err
		}

		// Update project to add collaborator
		update := bson.M{
			"$addToSet": bson.M{"collaborators": collaboratorObjectID},
			"$set":      bson.M{"updated_at": time.Now()},
		}
		
		_, err = pc.projectCollection.UpdateOne(
			sessionContext, 
			filter, 
			update,
		)
		return err
	})

	// Check transaction result
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Failed to add collaborator. Project not found or you're not the owner."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Collaborator added successfully",
	})
}