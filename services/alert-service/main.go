package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/k1ngalph0x/atlas/services/alert-service/api"
	"github.com/k1ngalph0x/atlas/services/alert-service/config"
	"github.com/k1ngalph0x/atlas/services/alert-service/db"
	"github.com/k1ngalph0x/atlas/services/alert-service/kafka"
	"github.com/k1ngalph0x/atlas/services/alert-service/models"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil{
		log.Fatalf("Error loading config: %v", err)
	}

	conn, err := db.ConnectDB()
	if err != nil{
		log.Fatalf("Error connecting to database: %v", err)
	}
	err = conn.AutoMigrate(&models.AlertRule{})
	if err != nil{
		log.Fatalf("Failed to migrate alert rule table: %v", err)
	}

	err = conn.AutoMigrate(&models.AlertLog{})
	if err != nil {
		log.Fatalf("Failed to migrate alert log table: %v", err)
	}

	handler := api.NewAlertHandler(conn, cfg)

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

	router.POST("/projects/:project_id/rules", handler.CreateAlertRule)
	router.GET("/projects/:project_id/rules", handler.GetAlertRules)
	router.DELETE("/projects/:project_id/rules/:rule_id", handler.DeleteAlertRule)
	router.GET("/projects/:project_id/alerts", handler.GetProjectAlerts)
	router.GET("/projects/:project_id/alerts/unread", handler.GetUnreadAlerts)
	router.POST("/alerts/:alert_id/acknowledge", handler.AcknowledgeAlert)
	router.Run(":8084")
}
