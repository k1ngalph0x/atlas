// models/alert.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IssueUpdateEvent struct {
	IssueID   string    `json:"issue_id"`
	ProjectID string    `json:"project_id"`
	Count     int       `json:"count"`
	Level     string    `json:"level"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AlertRule struct {
	ID          string    `gorm:"type:uuid;primaryKey" json:"id"`
	ProjectID   string    `gorm:"type:uuid;not null;index" json:"project_id"`
	Name        string    `gorm:"not null" json:"name"`
	Condition   string    `gorm:"not null" json:"condition"`
	Threshold   int       `gorm:"default:0" json:"threshold"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type AlertLog struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	RuleID    string    `gorm:"type:uuid;not null;index" json:"rule_id"`
	IssueID   string    `gorm:"type:uuid;not null;index" json:"issue_id"`
	ProjectID string    `gorm:"type:uuid;not null;index" json:"project_id"`
	Message   string    `gorm:"not null" json:"message"`
	Acknowledged bool   `gorm:"default:false" json:"acknowledged"`
	FiredAt   time.Time `gorm:"autoCreateTime" json:"fired_at"`
}

func (a *AlertRule) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

func (a *AlertLog) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}