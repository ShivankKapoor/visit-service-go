package routes

import (
	"visit-service/internal/handlers"

	"github.com/gin-gonic/gin"
)

func Configure(r *gin.Engine) {
	v1 := r.Group("")
	{
		// We reference the function from the handlers package
		v1.GET("/", handlers.Home)
	}
}
