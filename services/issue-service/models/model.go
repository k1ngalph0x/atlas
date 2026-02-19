package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Event struct {
	ProjectID  string    `json:"project_id"`
	Timestamp  time.Time `json:"timestamp"`
	Level      string    `json:"level"`
	Message    string    `json:"message"`
	StackTrace string    `json:"stack_trace"`
}

type Issue struct {
	ID          string    `gorm:"type:uuid;primaryKey" json:"id"`
	Fingerprint string    `gorm:"uniqueIndex:idx_project_fp;not null" json:"fingerprint"`
	ProjectID   string    `gorm:"uniqueIndex:idx_project_fp;type:uuid;not null;index" json:"project_id"`
	Title       string    `gorm:"not null" json:"title"`
	Level       string    `gorm:"not null" json:"level"`
	Count       int       `gorm:"default:1" json:"count"`
	StackTrace  string    `gorm:"type:text" json:"stack_trace"` 
	FirstSeen   time.Time `gorm:"autoCreateTime" json:"first_seen"`
	LastSeen    time.Time `gorm:"autoUpdateTime" json:"last_seen"`
	Status      string    `gorm:"default:'open'" json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}


func (i *Issue) BeforeCreate(tx *gorm.DB) error{
	if i.ID == ""{
		i.ID = uuid.New().String()
	}

	return nil
}