package main

import (
	"log"
	"visit-service/internal/database"
	"visit-service/internal/middleware"
	"visit-service/internal/routes"
	"visit-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	if err := database.Init(); err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	if err := service.SetupDailyReportTask(); err != nil {
		log.Fatal(err)
	}
	defer service.StopDailyReportTask()

	// Creates a router with default logging and recovery middleware
	r := gin.Default()
	r.Use(middleware.RateLimit())

	routes.Configure(r)

	r.Run(":8088")
}
