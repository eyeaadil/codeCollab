package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SessionType defines different types of coding sessions
type SessionType string

const (
	TypePrivate   SessionType = "private"
	TypePublic    SessionType = "public"
	TypeWorkspace SessionType = "workspace"
)

// SessionLanguage represents supported programming languages
type SessionLanguage string

const (
	LanguagePython   SessionLanguage = "python"
	LanguageGo       SessionLanguage = "go"
	LanguageJavaScript SessionLanguage = "javascript"
	LanguageRust     SessionLanguage = "rust"
	LanguageJava     SessionLanguage = "java"
)

// Session represents a collaborative coding session
type Session struct {
	ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Title           string               `bson:"title" json:"title" validate:"required,min=3,max=100"`
	Description     string               `bson:"description" json:"description,omitempty"`
	HostUserID      primitive.ObjectID   `bson:"host_user_id" json:"host_user_id"`
	
	// Collaboration Details
	// CollaboratorIDs []primitive.ObjectID `bson:"collaborator_ids" json:"collaborator_ids"`
	MaxParticipants int                  `bson:"max_participants" json:"max_participants"`
	
	// Session Configuration
	Type           SessionType       `bson:"type" json:"type"`
	Language       SessionLanguage   `bson:"language" json:"language"`
	IsPasswordProtected bool         `bson:"is_password_protected" json:"is_password_protected"`
	SessionPassword   string         `bson:"session_password,omitempty" json:"-"`
	
	// Tracking and Metadata
	CreatedAt       time.Time `bson:"created_at" json:"created_at"`
	LastActiveAt    time.Time `bson:"last_active_at" json:"last_active_at"`
	ExpiresAt       time.Time `bson:"expires_at" json:"expires_at"`
	
	// Session Status
	Status          string `bson:"status" json:"status"`
}