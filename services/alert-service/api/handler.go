package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/k1ngalph0x/atlas/services/alert-service/config"
	"github.com/k1ngalph0x/atlas/services/alert-service/models"
	"gorm.io/gorm"
)

type AlertHandler struct {
	DB     *gorm.DB
	Config *config.Config
	//Writer *kafka.Writer
}

func NewAlertHandler(db *gorm.DB, config *config.Config) *AlertHandler {
	return &AlertHandler{
		DB: db,
		Config: config,
	}
}

type CreateRuleRequest struct {
	Name      string `json:"name"      binding:"required"`
	Condition string `json:"condition" binding:"required,oneof=new_issue critical_error count_threshold"`
	Threshold int    `json:"threshold"`
}

func (h *AlertHandler) CreateAlertRule(c *gin.Context) {
	projectID := c.Param("project_id")

	var req CreateRuleRequest
	err := c.ShouldBindJSON(&req)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Condition == "count_threshold" && req.Threshold <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "threshold must be > 0 for count_threshold condition"})
		return
	}

	rule := models.AlertRule{
		ProjectID: projectID,
		Name:      req.Name,
		Condition: req.Condition,
		Threshold: req.Threshold,
		IsActive:  true,
	}

	result := h.DB.Create(&rule)
	if result.Error != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create alert rule"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"rule": rule})
}

func (h *AlertHandler) GetAlertRules(c *gin.Context) {
	projectID := c.Param("project_id")

	var rules []models.AlertRule
	result :=  h.DB.Where("project_id = ?", projectID).Order("created_at desc").Find(&rules)
	if result.Error != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch alert rules"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rules": rules})
}

func (h *AlertHandler) DeleteAlertRule(c *gin.Context) {
	projectID := c.Param("project_id")
	ruleID := c.Param("rule_id")

	result := h.DB.Where("id = ? AND project_id = ?", ruleID, projectID).Delete(&models.AlertRule{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete alert rule"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rule not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}


func(h *AlertHandler) ProcessAlert(e models.IssueUpdateEvent){
	var rules []models.AlertRule

	result := h.DB.Where("project_id = ? AND is_active = true", e.ProjectID).Find(&rules)
	if result.Error != nil{
		log.Printf("Failed to fetch alert rules: %v", result.Error)
		return
	}

	if len(rules) == 0{
		return 
	}

	for _, rule := range rules{
		if h.checkRule(rule, e){
			h.Alert(rule, e)
		}
	}
}

func (h *AlertHandler) checkRule(rule models.AlertRule, e models.IssueUpdateEvent) bool {
	switch rule.Condition{
	case "new_issue":
		return e.Count == 1

	case "critical_error":
		return e.Level == "critical" || e.Level == "error"
	
	case "count_threshold": 
		return e.Count > rule.Threshold

	default:
		return false
	}
}

func(h *AlertHandler) Alert(rule models.AlertRule, e models.IssueUpdateEvent){
	var existing models.AlertLog

	result := h.DB.Where("rule_id = ? AND issue_id = ?", rule.ID, e.IssueID).First(&existing)
	if result.Error == nil && rule.Condition != "count_threshold" {
		return
	}

	alertLog := models.AlertLog{
		RuleID:    rule.ID,
		IssueID:   e.IssueID,
		ProjectID: e.ProjectID,
		Message:   buildMessage(rule, e),
		FiredAt:   time.Now(),
	}

	result =  h.DB.Create(&alertLog)
	if result.Error != nil{
		log.Printf("Failed to save alert log: %v", result.Error)
		return
	}

	log.Printf("ALERT FIRED [%s] %s", rule.Name, alertLog.Message)
}

func buildMessage(rule models.AlertRule, e models.IssueUpdateEvent) string {
	switch rule.Condition {
	case "new_issue":
		return "New issue detected in project " + e.ProjectID
	case "critical_error":
		return "Critical/error level issue detected: " + e.IssueID
	case "count_threshold":
		return fmt.Sprintf("Issue %s exceeded threshold of %d", e.IssueID, rule.Threshold)
	default:
		return "Alert triggered"
	}
}

func (h *AlertHandler) GetProjectAlerts(c *gin.Context){
	projectID := c.Param("project_id")

	var alerts []models.AlertLog
	result := h.DB.Where("project_id = ?", projectID).Order("fired_at desc").Limit(50).Find(&alerts)

	if result.Error != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch alerts"})
		return 
	}

	c.JSON(http.StatusOK, gin.H{"alerts": alerts})
}

func (h *AlertHandler) GetUnreadAlerts(c *gin.Context) {
	projectID := c.Param("project_id")

	var count int64

	result := h.DB.Model(&models.AlertLog{}).Where("project_id = ? AND acknowledged = false", projectID).Count(&count)

	if result.Error != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch alert count"})
		return 
	}

	c.JSON(http.StatusOK, gin.H{"unread_count": count})
}

func (h *AlertHandler) AcknowledgeAlert(c *gin.Context) {
	alertID := c.Param("alert_id")

	result := h.DB.Model(&models.AlertLog{}).Where("id = ?", alertID).Update("acknowledged", true)

	if result.Error != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to acknowledged alert"})
		return 
	}

	c.JSON(http.StatusOK, gin.H{"status": "acknowledged"})
}