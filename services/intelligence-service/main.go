package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/k1ngalph0x/atlas/services/identity-service/db"
	"github.com/k1ngalph0x/atlas/services/intelligence-service/api"
	"github.com/k1ngalph0x/atlas/services/intelligence-service/config"
	"github.com/k1ngalph0x/atlas/services/intelligence-service/kafka"
	"github.com/k1ngalph0x/atlas/services/intelligence-service/models"
	"github.com/k1ngalph0x/atlas/services/intelligence-service/rabbitmq"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	conn, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("DB error: %v", err)
	}

	if err := conn.AutoMigrate(&models.IssueInsight{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	rabbitConn, rabbitCh, err := rabbitmq.Connect(config)
	if err != nil {
		log.Fatalf("RabbitMQ error: %v", err)
	}
	defer rabbitConn.Close()
	defer rabbitCh.Close()

	handler := api.NewAIHandler(conn, config, rabbitCh)


	go handler.StartWorkers(3)

	go kafka.Consume(handler)

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})
	router.GET("/projects/:project_id/insights", handler.GetProjectInsights)
	router.GET("/issues/:issue_id/insight", handler.GetIssueInsight)
	router.Run(":8083")

}