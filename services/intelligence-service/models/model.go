package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IssueInsight struct {
	ID          string    `gorm:"type:uuid;primaryKey" json:"id"`
	IssueID     string    `gorm:"type:uuid;not null;uniqueIndex" json:"issue_id"`
	ProjectID   string    `gorm:"type:uuid;not null;index" json:"project_id"`
	Summary     string    `gorm:"type:text;not null" json:"summary"`
	RootCause   string    `gorm:"type:text" json:"root_cause"`
	Remediation string    `gorm:"type:text" json:"remediation"`
	TokensUsed  int       `gorm:"default:0" json:"tokens_used"`
	ModelUsed   string    `gorm:"default:'llama3.2:3b'" json:"model_used"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type IssueUpdateEvent struct {
	IssueID   string    `json:"issue_id"`
	ProjectID string    `json:"project_id"`
	Count     int       `json:"count"`
	Level     string    `json:"level"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AIQueue struct {
	IssueID    string `json:"issue_id"`
	ProjectID  string `json:"project_id"`
	Title      string `json:"title"`
	Level      string `json:"level"`
	Count      int    `json:"count"`
	StackTrace string `json:"stack_trace"`
	Service    string `json:"service"`
}

func (i *IssueInsight) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
	return nil
}


