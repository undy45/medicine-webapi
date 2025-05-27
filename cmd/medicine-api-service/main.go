package main

import (
	"github.com/gin-gonic/gin"
	"github.com/undy45/medicine-webapi/api"
	"github.com/undy45/medicine-webapi/internal/medicine"
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
	handleFunctions := &medicine.ApiHandleFunctions{
		OrderStatusesAPI:     medicine.NewOrderStatusesApi(),
		MedicineInventoryAPI: medicine.NewMedicineInventoryAPI(),
		MedicineOrderAPI:     medicine.NewMedicineOrderAPI(),
	}
	medicine.NewRouterWithGinEngine(engine, *handleFunctions)
	engine.GET("/openapi", api.HandleOpenApi)
	engine.Run(":" + port)
}
