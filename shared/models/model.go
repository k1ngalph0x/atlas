package models

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	ID             string    `gorm:"type:uuid;primaryKey" json:"id"`
	ProjectName    string    `gorm:"not null" json:"project_name"`
	OrganizationID string    `gorm:"type:uuid;not null;index" json:"organization_id"`
	//APIKey         string    `gorm:"type:varchar(128);unique;not null;index" json:"api_key"`
	APIKey         string    `gorm:"type:varchar(128);unique;not null;index" json:"-"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (p *Project) BeforeCreate(tx *gorm.DB) error{
	if p.ID == ""{
		p.ID = uuid.New().String()
	}
	// if p.APIKey == ""{
	// 	p.APIKey = generateAPIKey()
	// }
	return nil
}

type ProjectWithRawKey struct {
	Project
	RawAPIKey string `json:"api_key"`
}

func NewProject(name, orgID string)(*Project, string){
	key := generateAPIKey()
	return &Project{
		ProjectName:    name,
		OrganizationID: orgID,
		APIKey:         HashAPIKey(key),
	}, key
}


func generateAPIKey() string{
	randomBytes := make([]byte, 32)
	rand.Read(randomBytes)
	return "atlas_" + hex.EncodeToString(randomBytes)
}

func HashAPIKey(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

