package api

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/k1ngalph0x/atlas/services/issue-service/config"
	"github.com/k1ngalph0x/atlas/services/issue-service/models"
	publisher "github.com/k1ngalph0x/atlas/services/issue-service/utils"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

type IssueHandler struct{
	DB *gorm.DB
	Config *config.Config
	Writer *kafka.Writer
}

func NewIssueHandler(db *gorm.DB, config *config.Config, writer *kafka.Writer) *IssueHandler {
	return &IssueHandler{
		DB: db,
		Config: config,
		Writer: writer,
	}
}

type IssueUpdateEvent struct {
	IssueID   string    `json:"issue_id"`
	ProjectID string    `json:"project_id"`
	Count     int       `json:"count"`
	Level     string    `json:"level"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}


func NewIssueResolvedWriter(config *config.Config) *kafka.Writer{
	return &kafka.Writer{
	Addr:     kafka.TCP(config.KAFKA.Brokers...),
	Topic:    "issue-resolved",
	Balancer: &kafka.LeastBytes{},
	}
}

func NewIssueUpdateWriter(config *config.Config) *kafka.Writer{
	return &kafka.Writer{
	Addr:     kafka.TCP(config.KAFKA.Brokers...),
	Topic:    "issue-updates",
	Balancer: &kafka.LeastBytes{},
	}
}

func generateFingerprint(message string, stack *string) string {
	raw := message
	// if stack != nil{
	// 	raw += "|" + *stack
	// }

	hash := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(hash[:])
}


func(h *IssueHandler) ProcessEvents(e models.Event) {
	var issue models.Issue
	//var event models.Event
	fp := generateFingerprint(e.Message, &e.StackTrace)

	result := h.DB.Where("project_id = ? AND fingerprint = ?", e.ProjectID, fp).First(&issue)
	if result.Error == nil{
		err := h.DB.Model(&issue).Updates(map[string]interface{}{
			"count":     gorm.Expr("count + ?", 1),
			"last_seen": time.Now(),
		}).Error

		if err != nil{
			log.Printf("Failed to update issue %s: %v", issue.ID, err)
			return 
		}

		h.DB.First(&issue, "id = ?", issue.ID)

		updateEvent := IssueUpdateEvent{
			IssueID:   issue.ID,
			ProjectID: issue.ProjectID,
			Count:     issue.Count,
			Level:     issue.Level,
			Status:    issue.Status,
			UpdatedAt: time.Now(),
		}

		publisher.PublishEvent(h.Writer, issue.ProjectID, updateEvent)

		return
	}

	newIssue := models.Issue{
		ID:          uuid.New().String(),
		ProjectID:   e.ProjectID,
		Fingerprint: fp,
		Title:       e.Message,
		Level:       e.Level,
		Count:       1,
		StackTrace:  e.StackTrace,
		FirstSeen:   time.Now(),
		LastSeen:    time.Now(),
		Status:      "open",
	}

	err := h.DB.Create(&newIssue).Error
	if err != nil{
		log.Println("Failed to create issue:", err)
		return
	}

	updateEvent := IssueUpdateEvent{
		IssueID:   newIssue.ID,
		ProjectID: newIssue.ProjectID,
		Count:     newIssue.Count,
		Level:     newIssue.Level,
		Status:    newIssue.Status,
		UpdatedAt: time.Now(),
	}

	publisher.PublishEvent(h.Writer, newIssue.ProjectID, updateEvent)
}

func(i *IssueHandler) GetProjectIssue(c *gin.Context) {
		var issues []models.Issue
		projectID := c.Param("project_id")

		result := i.DB.Where("project_id = ?", projectID).Order("last_seen desc").Find(&issues)
		if result.Error != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(200, gin.H{"issues": issues}) 
}

func(i *IssueHandler) GetIssueDetail(c *gin.Context){
	projectID := c.Param("project_id")
	issueID := c.Param("issue_id")

	var issue models.Issue
	result := i.DB.Where("id = ? AND project_id = ?", issueID, projectID).First(&issue)

	if result.Error != nil{
		c.JSON(http.StatusNotFound, gin.H{"error": "Issue not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"issue": issue})
}

func (i *IssueHandler) GetProjectOverview(c *gin.Context) {
	projectID := c.Param("project_id")

	var stats struct {
		TotalIssues    int64
		OpenIssues     int64
		ResolvedIssues int64
		CriticalCount  int64
		ErrorCount     int64
	}

	result := i.DB.Model(&models.Issue{}).
		Select(`
			COUNT(*) as total_issues,
			COUNT(*) FILTER (WHERE status = 'open') as open_issues,
			COUNT(*) FILTER (WHERE status = 'resolved') as resolved_issues,
			COUNT(*) FILTER (WHERE level = 'critical') as critical_count,
			COUNT(*) FILTER (WHERE level = 'error') as error_count
		`).
		Where("project_id = ?", projectID).
		Scan(&stats)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_issues":    stats.TotalIssues,
		"open_issues":     stats.OpenIssues,
		"resolved_issues": stats.ResolvedIssues,
		"critical_count":  stats.CriticalCount,
		"error_count":     stats.ErrorCount,
	})
}
