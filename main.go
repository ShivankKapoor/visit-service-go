package main

import (
	"log"
	"visit-service/internal/database"
	"visit-service/internal/middleware"
	"visit-service/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	if err := database.Init(); err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// Creates a router with default logging and recovery middleware
	r := gin.Default()
	r.Use(middleware.RateLimit())

	routes.Configure(r)

	r.Run(":8088")
}
