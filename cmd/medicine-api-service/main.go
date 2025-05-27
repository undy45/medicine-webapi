package main

import (
	"github.com/gin-gonic/gin"
	"github.com/undy45/medicine-webapi/api"
	"log"
	"os"
	"strings"
)

func main() {
	log.Printf("Server started")
	port := os.Getenv("MEDICINE_API_PORT")
	if port == "" {
		port = "8080"
	}
	environment := os.Getenv("MEDICINE_API_ENVIRONMENT")
	if !strings.EqualFold(environment, "production") { // case insensitive comparison
		gin.SetMode(gin.DebugMode)
	}
	engine := gin.New()
	engine.Use(gin.Recovery())
	// request routings
	engine.GET("/openapi", api.HandleOpenApi)
	engine.Run(":" + port)
}
