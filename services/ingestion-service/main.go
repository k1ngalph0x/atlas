package main

import (
	"log"

	"github.com/gin-gonic/gin"
	handler "github.com/k1ngalph0x/atlas/services/ingestion-service/api"
	"github.com/k1ngalph0x/atlas/services/ingestion-service/config"
	"github.com/k1ngalph0x/atlas/services/ingestion-service/db"
	"github.com/k1ngalph0x/atlas/services/ingestion-service/kafka"
	"github.com/k1ngalph0x/atlas/services/ingestion-service/middleware"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil{
		log.Fatalf("Error loading config: %v", err)
	}

	conn, err := db.ConnectDB()
	if err != nil{
		log.Fatalf("Error connecting to database: %v", err)
	}

	kafka.InitKafka(config)

	authMiddleware := middleware.NewAuthMiddleware(conn)

	router := gin.Default()
	router.Use(gin.Logger())

	router.Use(authMiddleware.ValidateAPIKey())
	api := router.Group("/api")
	{
		api.POST("/ingest/events", handler.Ingest)
	}

	router.Run(":8081")
}

