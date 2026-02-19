package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/k1ngalph0x/atlas/services/identity-service/config"
	"github.com/k1ngalph0x/atlas/services/identity-service/models"
	sharedModels "github.com/k1ngalph0x/atlas/shared/models"
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
	err := c.ShouldBindJSON(&req)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
	}

	var org models.Organization
	result := p.DB.Where("id = ? AND user_id = ?", req.OrganizationID, userId).First(&org)
	if result.Error != nil{
		c.JSON(http.StatusForbidden, gin.H{"error": "Organization not found or access denied"})
		return
	}

	// project := sharedModels.Project{
	// 	ProjectName: req.ProjectName,
	// 	OrganizationID: req.OrganizationID,
	// }

	project, rawKey := sharedModels.NewProject(req.ProjectName, req.OrganizationID)

	result = p.DB.Create(&project)
	if result.Error != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	// c.JSON(http.StatusCreated, gin.H{
	// 	"message": "Project created successfully",
	// 	"project": project,
	// })

	c.JSON(http.StatusCreated, gin.H{
		"message": "Project created successfully",
		"project": sharedModels.ProjectWithRawKey{
			Project:   *project,
			RawAPIKey: rawKey,
	}})
}

func(p *ProjectHandler) GetOrganizations(c *gin.Context){
	userId := c.GetString("user_id")
	if userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}	

	var organizations []models.Organization
	result := p.DB.Where("user_id = ?", userId).Find(&organizations)
	if result.Error != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch organizations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"organizations": organizations})
}

func  (p *ProjectHandler) GetProjects(c *gin.Context) {
	userId := c.GetString("user_id")
	if userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var organizations []models.Organization
	result := p.DB.Where("user_id = ?", userId).Find(&organizations)
	if result.Error != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch organizations"})
		return
	}

	orgs := make([]string, len(organizations))
	for i, org := range organizations{
		orgs[i] = org.ID
	}

	var projects []sharedModels.Project
	result = p.DB.Where("organization_id IN ?", orgs).Find(&projects)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"projects": projects})

}