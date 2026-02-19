package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	sharedModels "github.com/k1ngalph0x/atlas/shared/models"
	"gorm.io/gorm"
)

type APIKeyMiddleware struct {
	DB *gorm.DB
}

func NewAuthMiddleware(db *gorm.DB) *APIKeyMiddleware{
	return &APIKeyMiddleware{
		DB: db,
	}
}

func (a *APIKeyMiddleware) ValidateAPIKey() gin.HandlerFunc{
	return func(c *gin.Context){
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == ""{
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing API key"})
			c.Abort()
			return
		}

		hashedKey := sharedModels.HashAPIKey(apiKey)

		var project sharedModels.Project
		//result := a.DB.Where("api_key = ?", apiKey).First(&project)
		result := a.DB.Where("api_key = ?", hashedKey).First(&project)
		if result.Error != nil{
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		c.Set("project_id", project.ID)
		c.Set("organization_id", project.OrganizationID)
		c.Next()
	}
}