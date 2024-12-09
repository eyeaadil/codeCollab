package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRole defines different user roles in the platform
type UserRole string

const (
	RoleUser     UserRole = "user"
	RoleAdmin    UserRole = "admin"
	RolePremium  UserRole = "premium"
)

// User represents a user in the collaborative coding platform
type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username       string             `bson:"username" json:"username" validate:"required,min=3,max=50"`
	Email          string             `bson:"email" json:"email" validate:"required,email"`
	PasswordHash   string             `bson:"password_hash" json:"-"`
	ProfilePicture string             `bson:"profile_picture" json:"profile_picture,omitempty"`
	Role           UserRole           `bson:"role" json:"role"`
	
	// User Statistics
	TotalProjectsCreated int `bson:"total_projects_created" json:"total_projects_created"`
	TotalCollaborations  int `bson:"total_collaborations" json:"total_collaborations"`

	// Account Management
	IsVerified     bool      `bson:"is_verified" json:"is_verified"`
	LastLoginAt   time.Time `bson:"last_login_at" json:"last_login_at"`
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at" json:"updated_at"`
}

// UserPublicProfile represents a sanitized version of user data for public display
type UserPublicProfile struct {
	ID       primitive.ObjectID `json:"id"`
	Username string             `json:"username"`
	Role     UserRole           `json:"role"`
}