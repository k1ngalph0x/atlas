package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/k1ngalph0x/atlas/services/identity-service/api"
	"github.com/k1ngalph0x/atlas/services/identity-service/config"
	"github.com/k1ngalph0x/atlas/services/identity-service/db"
	"github.com/k1ngalph0x/atlas/services/identity-service/middleware"
	identityModels "github.com/k1ngalph0x/atlas/services/identity-service/models"
	sharedModels "github.com/k1ngalph0x/atlas/shared/models"
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

	err = conn.AutoMigrate(&identityModels.User{})
	if err != nil{
		log.Fatalf("Failed to migrate user table: %v", err)
	}

	err = conn.AutoMigrate(&identityModels.Organization{})
	if err != nil{
		log.Fatalf("Failed to migrate organization table: %v", err)
	}

	err = conn.AutoMigrate(&sharedModels.Project{})
	if err != nil{
		log.Fatalf("Failed to migrate Project table: %v", err)
	}

	authHandler := api.NewAuthHandler(conn, config)
	authMiddleware := middleware.NewAuthMiddleware(config.TOKEN.JwtKey)
	projectHandler := api.NewProjectHandler(conn, config)

	router := gin.Default()
	router.Use(gin.Logger())

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

	auth := router.Group("/auth")
	{
		auth.POST("/signup", authHandler.SignUp)
		auth.POST("/signin", authHandler.SignIn)
	}

	router.Use(authMiddleware.RequireAuth())
	project := router.Group("/project")
	{
		project.POST("/create-organization", projectHandler.CreateOrganization)
		project.POST("/create-project", projectHandler.CreateProject)
		project.GET("/organizations", projectHandler.GetOrganizations) 
		project.GET("/projects", projectHandler.GetProjects)  
	}

	fmt.Println("Running Identity service")
	router.Run(":8080")

}

