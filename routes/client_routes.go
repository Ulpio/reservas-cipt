package routes

import (
	"github.com/Ulpio/reservas-cipt/handlers"
	"github.com/Ulpio/reservas-cipt/middleware"
	"github.com/gin-gonic/gin"
)

func SetupClientRoutes(r *gin.Engine) {
	grupo := r.Group("/client")
	grupo.Use(middleware.JWTAuthMiddleware())
	{
		grupo.POST("/buscar-criar", handlers.BuscarOuCriarClienteHandler)
	}
}
