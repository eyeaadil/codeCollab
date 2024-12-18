package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProjectModel represents the Project collection
type ProjectModel struct {
	ProjectID    primitive.ObjectID   `bson:"_id,omitempty"`          // Auto-generated project ID
	ProjectName  string               `bson:"project_name,omitempty"` // Project name
	Description  string               `bson:"description,omitempty"`  // Project description
	OwnerID      primitive.ObjectID   `bson:"owner_id,omitempty"`     // Owner of the project
	Collaborators []primitive.ObjectID `bson:"collaborators,omitempty"` // List of collaborator user IDs
	Files        []primitive.ObjectID `bson:"files,omitempty"`        // List of file IDs
	CreatedAt    time.Time            `bson:"created_at,omitempty"`   // Creation timestamp
	UpdatedAt    time.Time            `bson:"updated_at,omitempty"`   // Last modification timestamp
	Visibility   string               `bson:"visibility,omitempty"`   // Visibility: "public" or "private"
}
