package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FileType defines different types of files in a session
type FileType string

const (
	TypeSourceCode FileType = "source_code"
	TypeMarkdown   FileType = "markdown"
	TypeConfig     FileType = "configuration"
	TypeText       FileType = "text"
)

// File represents a file within a coding session
type File struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	// SessionID  primitive.ObjectID `bson:"session_id" json:"session_id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	FolderID   primitive.ObjectID `bson:"folder_id" json:"folder_id"` 
	// File Details
	Name        string    `bson:"name" json:"name" validate:"required"`
	Content     string    `bson:"content" json:"content"`
	Type        FileType  `bson:"type" json:"type"`
	Language    string    `bson:"language" json:"language"`
	
	// Version Control
	Version     int       `bson:"version" json:"version"`
	ParentFileID primitive.ObjectID `bson:"parent_file_id,omitempty" json:"parent_file_id,omitempty"`
	
	// Metadata
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
	
	// Collaboration Tracking
	LastEditedBy primitive.ObjectID `bson:"last_edited_by" json:"last_edited_by"`
}

// FileVersion tracks different versions of a file
type FileVersion struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FileID     primitive.ObjectID `bson:"file_id" json:"file_id"`
	Content    string             `bson:"content" json:"content"`
	Version    int                `bson:"version" json:"version"`
	EditedBy   primitive.ObjectID `bson:"edited_by" json:"edited_by"`
	EditedAt   time.Time          `bson:"edited_at" json:"edited_at"`
}