package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	UserID    string `gorm:"type:uuid;primaryKey" json:"id"`
	Email     string `gorm:"unique;not null" json:"email"`
	Password  string `gorm:"not null" json:"-"`
	CreatedAt time.Time	`gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Organization struct{
	ID string `gorm:"type:uuid;primaryKey" json:"id"`
	OrganizationName string `gorm:"unique;not null" json:"organization_name"`
	UserID string `gorm:"type:uuid;not null" json:"user_id"`
	CreatedAt time.Time	`gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Project struct{
	ID string `gorm:"type:uuid;primaryKey" json:"id"`
	ProjectName string `gorm:"not null" json:"project_name"`
	OrganizationID string `gorm:"type:uuid;not null;index" json:"organization_id"`
	CreatedAt time.Time	`gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func(u *User) BeforeCreate(tx *gorm.DB) error {
	if u.UserID == ""{
		u.UserID = uuid.New().String()
	}
	return nil
}

func(o *Organization) BeforeCreate(tx *gorm.DB) error {
	if o.ID == ""{
		o.ID = uuid.New().String()
	}
	return nil
}

func (p *Project) BeforeCreate(tx *gorm.DB) error{
	if p.ID == ""{
		p.ID = uuid.New().String()
	}
	return nil
}