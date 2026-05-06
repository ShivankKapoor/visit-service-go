package main

import (
	"visit-service/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Creates a router with default logging and recovery middleware
	r := gin.Default()

	routes.Configure(r)

	r.Run(":8088")
}
