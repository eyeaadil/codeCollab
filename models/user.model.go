package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRole defines different user roles in the platform
type UserRole string

const (
	RoleAdmin        UserRole = "admin"
	RoleCollaborator UserRole = "collaborator"
	RoleUser         UserRole = "user" // Default role for newly registered users
)

// User represents a user in the collaborative coding platform
type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username       string             `bson:"username" json:"username" validate:"required,min=3,max=50"`
	Email          string             `bson:"email" json:"email" validate:"required,email"`
	PasswordHash   string             `bson:"password_hash" json:"-"` // Omit for security
	// Role           UserRole           `bson:"role" json:"role"`
	ProfilePicture string             `bson:"profile_picture" json:"profile_picture,omitempty"`

	// User Statistics
	TotalProjectsCreated int `bson:"total_projects_created" json:"total_projects_created"`
	TotalCollaborations  int `bson:"total_collaborations" json:"total_collaborations"`
// Collaborator Reference
    CollaboratorID primitive.ObjectID `bson:"collaborator_id,omitempty" json:"collaborator_id,omitempty"`

	// Account Management
	LastLoginAt time.Time `bson:"last_login_at,omitempty" json:"last_login_at,omitempty"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
	IsVerified  bool      `bson:"is_verified" json:"is_verified"`
}

// UserPublicProfile represents a sanitized version of user data for public display
type UserPublicProfile struct {
	ID             primitive.ObjectID `json:"id"`
	Username       string             `json:"username"`
	Role           UserRole           `json:"role"`
	ProfilePicture string             `json:"profile_picture,omitempty"`
	TotalProjects  int                `json:"total_projects_created"`
	TotalCollabs   int                `json:"total_collaborations"`
	IsVerified     bool               `json:"is_verified"`
}
