package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CollaboratorRole defines roles of collaborators in a session
type CollaboratorRole string

const (
	RoleViewer CollaboratorRole = "viewer"
	RoleEditor CollaboratorRole = "editor"
	// RoleAdmin  CollaboratorRole = "admin"
)

// Collaborator represents a participant in a session
type Collaborator struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SessionID    primitive.ObjectID `bson:"session_id" json:"session_id"`
	UserID       primitive.ObjectID `bson:"user_id" json:"user_id"`
	// Role         CollaboratorRole   `bson:"role" json:"role"` // viewer/editor/admin
	AddedBy      primitive.ObjectID `bson:"added_by" json:"added_by"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	LastModified time.Time          `bson:"last_modified" json:"last_modified"`
}
