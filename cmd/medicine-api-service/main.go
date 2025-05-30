package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/undy45/medicine-webapi/api"
	"github.com/undy45/medicine-webapi/internal/db_service"
	"github.com/undy45/medicine-webapi/internal/medicine"
	"log"
	"os"
	"strings"
	"time"
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
	corsMiddleware := cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{""},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	})
	engine.Use(corsMiddleware)

	// setup context update  middleware
	ambulanceSvc := db_service.NewMongoService[medicine.Ambulance](db_service.MongoServiceConfig{
		Collection: "ambulance",
	})
	defer ambulanceSvc.Disconnect(context.Background())
	statusSvc := db_service.NewMongoService[medicine.Status](db_service.MongoServiceConfig{
		Collection: "status",
	})
	defer statusSvc.Disconnect(context.Background())
	engine.Use(func(ctx *gin.Context) {
		ctx.Set("db_service_ambulance", ambulanceSvc)
		ctx.Set("db_service_status", statusSvc)
		ctx.Next()
	})
	//engine.Use(func(ctx *gin.Context) {
	//	ctx.Set("db_service", dbService)
	//	ctx.Next()
	//})
	// request routings
	handleFunctions := &medicine.ApiHandleFunctions{
		OrderStatusesAPI:     medicine.NewOrderStatusesApi(),
		MedicineInventoryAPI: medicine.NewMedicineInventoryAPI(),
		MedicineOrderAPI:     medicine.NewMedicineOrderAPI(),
		AmbulancesAPI:        medicine.NewAmbulancesAPI(),
	}
	medicine.NewRouterWithGinEngine(engine, *handleFunctions)
	engine.GET("/openapi", api.HandleOpenApi)
	engine.Run(":" + port)
}
