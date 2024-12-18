package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// VersionEntry represents a single version in file's version history
type VersionEntry struct {
	VersionID  primitive.ObjectID `bson:"version_id,omitempty"`   // Unique version ID
	Content    string             `bson:"content,omitempty"`      // File content snapshot
	Timestamp  time.Time          `bson:"timestamp,omitempty"`    // Time of the snapshot
	ModifiedBy primitive.ObjectID `bson:"modified_by,omitempty"`  // User ID who modified this version
}
