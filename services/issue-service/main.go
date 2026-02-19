package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/k1ngalph0x/atlas/services/issue-service/api"
	"github.com/k1ngalph0x/atlas/services/issue-service/config"
	"github.com/k1ngalph0x/atlas/services/issue-service/db"
	"github.com/k1ngalph0x/atlas/services/issue-service/kafka"
	"github.com/k1ngalph0x/atlas/services/issue-service/models"
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
	
	writer := api.NewIssueUpdateWriter(config)
	defer writer.Close()

	handler := api.NewIssueHandler(conn,config, writer)

	err = conn.AutoMigrate(&models.Issue{})
    if err != nil {
        log.Fatalf("Failed to migrate issue table: %v", err)
    }


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

	router.GET("/projects/:project_id/issues", handler.GetProjectIssue)
	router.GET("/projects/:project_id/issues/:issue_id", handler.GetIssueDetail)
	router.GET("/projects/:project_id/overview", handler.GetProjectOverview)
	router.Run(":8082") 
}

