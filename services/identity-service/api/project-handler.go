package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/k1ngalph0x/atlas/services/identity-service/config"
	"github.com/k1ngalph0x/atlas/services/identity-service/models"
	"gorm.io/gorm"
)

type ProjectHandler struct {
	DB     *gorm.DB
	Config *config.Config
}

type CreateOrgRequest struct{
	OrganizationName string `json:"organization_name" binding:"required"`
}

type CreateProjectRequest struct{
	OrganizationID string `json:"organization_id" binding:"required,uuid"`
	ProjectName string `json:"project_name" binding:"required"`
}

func NewProjectHandler(db *gorm.DB, config *config.Config) *ProjectHandler {
	return &ProjectHandler{
		DB:     db,
		Config: config,
	}
}

func(p *ProjectHandler) CreateOrganization(c *gin.Context){
	userId := c.GetString("user_id")
	if userId == ""{
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usnauthorized"})
		return
	}
	var req CreateOrgRequest	
	
	err := c.ShouldBindJSON(&req)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
	}

	//orgname := strings.TrimSpace(req.OrganizationName)

	org := models.Organization{
		OrganizationName: req.OrganizationName,
		UserID: userId,
	}

	result := p.DB.Create(&org)

	if result.Error != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create organization"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":"Organization created successfully",
		"organization": org,
	})
}

func (p *ProjectHandler) CreateProject(c *gin.Context){
	userId := c.GetString("user_id")
	if userId == ""{
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usnauthorized"})
		return
	}
	var req CreateProjectRequest
	var org models.Organization
	err := c.ShouldBindJSON(&req)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
	}

	result := p.DB.Where("id = ? AND user_id = ?", req.OrganizationID, userId).First(&org)
	if result.Error != nil{
		c.JSON(http.StatusForbidden, gin.H{"error": "Organization not found or access denied"})
		return
	}

	project := models.Project{
		ProjectName: req.ProjectName,
		OrganizationID: req.OrganizationID,
	}

	result = p.DB.Create(&project)
	if result.Error != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Project created successfully",
		"project": project,
	})
}