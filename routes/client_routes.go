package routes

import (
	"github.com/Ulpio/reservas-cipt/handlers"
	"github.com/Ulpio/reservas-cipt/middleware"
	"github.com/gin-gonic/gin"
)

func SetupClientRoutes(r *gin.RouterGroup) {
	grupo := r.Group("/clientes")
	grupo.Use(middleware.JWTAuthMiddleware())
	{
		grupo.PATCH("/:id", handlers.UpdateClientHandler)
		grupo.GET("/:cpf", handlers.BuscarClientePorCPF)
		grupo.POST("/buscar-criar", handlers.BuscarOuCriarClienteHandler)
		grupo.GET("", handlers.GetAllClientes)
	}
}
