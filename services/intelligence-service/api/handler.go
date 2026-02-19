package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/k1ngalph0x/atlas/services/intelligence-service/config"
	"github.com/k1ngalph0x/atlas/services/intelligence-service/models"
	"github.com/k1ngalph0x/atlas/services/intelligence-service/ollama"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type AIHandler struct {
	DB       *gorm.DB
	Config   *config.Config
	RabbitCh *amqp.Channel
	Ollama   *ollama.Client
}

func NewAIHandler(db *gorm.DB, config *config.Config, rabbitCh *amqp.Channel) *AIHandler {
	ollamaClient := ollama.NewClient(config.OLLAMA.Url, config.OLLAMA.Model)
	return &AIHandler{
		DB:       db,
		Config:   config,
		RabbitCh: rabbitCh,
		Ollama:   ollamaClient,
	}
}

func(h *AIHandler) ProcessIssue(e models.IssueUpdateEvent){
	log.Printf("ProcessIssue called for issue %s, count: %d, level: %s", e.IssueID, e.Count, e.Level)
	var existing models.IssueInsight
	result := h.DB.Where("issue_id = ?", e.IssueID).First(&existing)
	if result.Error == nil{
		log.Printf("Issue %s already has insight", e.IssueID)
		return
	}

	if e.Count < 5 && e.Level != "critical" && e.Level != "error"{
		log.Printf("Low threshold (count: %d, level: %s)", e.Count, e.Level)
		return
	}

	var issue struct{
		ID         string
		Title      string
		Level      string
		Count      int
		StackTrace string
	}


	result = h.DB.Table("issues").Select("id, title, level, count, stack_trace").Where("id = ?", e.IssueID).Scan(&issue)

	if result.Error != nil{
		log.Printf("Failed to fetch issue details: %v", result.Error)
		return
	}

	queue := models.AIQueue{
		IssueID:   e.IssueID,
		ProjectID: e.ProjectID,
		Title:     issue.Title,
		Level:     issue.Level,
		Count:     issue.Count,
		StackTrace: issue.StackTrace,
	}

	body, err := json.Marshal(queue)
	if err != nil {
		log.Printf("Failed to marshal job: %v", err)
		return
	}

	err = h.RabbitCh.PublishWithContext(
		context.Background(),
		"",
		"ai-analysis-jobs",
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)

	if err != nil {
		log.Printf("Failed to publish job: %v", err)
		return
	}


	log.Printf("Published the job")

} 


func (h *AIHandler) StartWorkers(numWorkers int){
	for i:= 0; i< numWorkers; i++{
		go h.worker(i)
	}

	log.Printf("Started %d AI workers", numWorkers)
}

func (h *AIHandler) worker(id int){
	msgs, err := h.RabbitCh.Consume(
		"ai-analysis-jobs", 
		"",                 
		false,             
		false,              
		false,             
		false,             
		nil,               
	)

	if err != nil{
		log.Fatalf("Worker %d failed to consume: %v", id, err)	
	}

	for msg := range msgs{
		var job models.AIQueue

		err := json.Unmarshal(msg.Body, &job)
		if err != nil{
			log.Printf("Worker %d: invalid job: %v", id, err)
			msg.Nack(false, false)
			continue
		}

		log.Printf("Processing")

		result, err := h.Ollama.Analyze(job.Title, job.StackTrace, job.Level, job.Count)
		if err != nil{
			log.Printf("Worker %d: Ollama error: %v", id, err)
			msg.Nack(false, true)
			continue
		}

		insight := models.IssueInsight{
			IssueID:     job.IssueID,
			ProjectID:   job.ProjectID,
			Summary:     result.Summary,
			RootCause:   result.RootCause,
			Remediation: result.Remediation,
			ModelUsed:   h.Config.OLLAMA.Model,
		}

		res := h.DB.Create(&insight)
		if res.Error != nil{
			log.Printf("Worker %d: DB save error: %v", id, err)
			msg.Nack(false, true) 
			continue
		}

		msg.Ack(false)
	}
}

func (h *AIHandler) GetProjectInsights(c *gin.Context){
	projectID := c.Param("project_id")

	var insights []models.IssueInsight
	 
	result :=  h.DB.Where("project_id = ?", projectID).Order("created_at desc").Find(&insights)
	if result.Error != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch insights"})
		return 
	}

	c.JSON(http.StatusOK, gin.H{"insights": insights})
}

func (h *AIHandler) GetIssueInsight(c *gin.Context) {
	issueID := c.Param("issue_id")

	var insight models.IssueInsight

	result :=  h.DB.Where("issue_id = ?", issueID).First(&insight)

	if result.Error != nil{
		c.JSON(http.StatusNotFound, gin.H{"error": "Insight not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"insight": insight})
}