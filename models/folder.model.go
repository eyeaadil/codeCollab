package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Folder struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID   `bson:"user_id" json:"user_id"`          // User who created the folder
	SessionID primitive.ObjectID   `bson:"session_id" json:"session_id"`    // Associated session ID
	Name      string               `bson:"name" json:"name"`                // Name of the folder
	FileIDs   []primitive.ObjectID `bson:"file_ids" json:"file_ids"`        // List of file IDs associated with the folder
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time            `bson:"updated_at" json:"updated_at"`
	Version   int                  `bson:"version" json:"version"`          // Version of the folder (for updates)
}
