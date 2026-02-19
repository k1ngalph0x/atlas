package main

import (
	"errors"
	"math/rand"

	"github.com/gin-gonic/gin"
	atlas "github.com/k1ngalph0x/atlas-go-sdk"
)

func main() {
	client := atlas.NewClient("atlas_96b846403c49fcaca9cc7dfca209f574bdbad503fce174cda80f2393f7337055")

	router := gin.Default()

	router.Use(client.GinMiddleware())

	router.GET("/random-error", func(c *gin.Context) {
		if rand.Intn(2) == 0 {
			err := errors.New("RandomError: something went wrong randomly")
			client.CaptureError(err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Success!"})
	})

	router.GET("/db-error", func(c *gin.Context) {
		err := errors.New("DatabaseTimeoutError: timeout after 30s while connecting to user-db")
		client.CaptureError(err)
		c.JSON(500, gin.H{"error": err.Error()})
	})

	router.GET("/null-error", func(c *gin.Context) {
		err := errors.New("NullPointerException: cannot read property 'email' of null")
		client.CaptureError(err)
		c.JSON(500, gin.H{"error": err.Error()})
	})

	router.Run(":3000")

}