package routes

import (
	"github.com/Ulpio/reservas-cipt/handlers"
	"github.com/Ulpio/reservas-cipt/middleware"
	"github.com/gin-gonic/gin"
)

func SetupStrikeRoutes(r *gin.RouterGroup) {
	grupo := r.Group("/strikes")
	grupo.Use(middleware.JWTAuthMiddleware())
	{
		grupo.POST("", handlers.CreateStrikeHandler)
		grupo.GET("/client/:id", handlers.GetStrikesByClientHandler)
		grupo.DELETE("/:id", handlers.RevokeStrikeHandler)
	}
}
