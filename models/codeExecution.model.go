package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CodeExecutionStatus represents the status of code execution
type CodeExecutionStatus string

const (
	StatusQueued     CodeExecutionStatus = "queued"
	StatusRunning    CodeExecutionStatus = "running"
	StatusCompleted  CodeExecutionStatus = "completed"
	StatusFailed     CodeExecutionStatus = "failed"
	StatusTimedOut   CodeExecutionStatus = "timed_out"
)

// CodeExecution represents a code execution request and its result
type CodeExecution struct {
	ID             primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	UserID         primitive.ObjectID     `bson:"user_id" json:"user_id"`
	SessionID      primitive.ObjectID     `bson:"session_id" json:"session_id"`
	FileID         primitive.ObjectID     `bson:"file_id" json:"file_id"`
	
	// Code Details
	Language       SessionLanguage        `bson:"language" json:"language"`
	SourceCode     string                 `bson:"source_code" json:"source_code"`
	
	// Execution Metadata
	Status         CodeExecutionStatus    `bson:"status" json:"status"`
	ExecutionTime  time.Duration          `bson:"execution_time" json:"execution_time"`
	MemoryUsed     int64                  `bson:"memory_used" json:"memory_used"`
	
	// Execution Results
	Output         string                 `bson:"output" json:"output"`
	ErrorMessage   string                 `bson:"error_message" json:"error_message,omitempty"`
	
	// Timestamps
	CreatedAt      time.Time              `bson:"created_at" json:"created_at"`
	CompletedAt    time.Time              `bson:"completed_at" json:"completed_at,omitempty"`
}