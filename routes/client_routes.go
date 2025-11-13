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
		// Listar todos
		grupo.GET("", handlers.GetAllClientes)
		
		// Criar novo cliente
		grupo.POST("", handlers.CreateClientHandler)
		
		// Buscar ou criar (endpoint legado)
		grupo.POST("/buscar-criar", handlers.BuscarOuCriarClienteHandler)
		
		// Buscar por CPF
		grupo.GET("/cpf/:cpf", handlers.BuscarClientePorCPF)
		
		// Buscar por telefone
		grupo.GET("/phone/:phone", handlers.BuscarClientePorTelefone)
		
		// Buscar por ID (deve vir depois dos endpoints específicos)
		grupo.GET("/:id", handlers.BuscarClientePorID)
		
		// Atualizar
		grupo.PATCH("/:id", handlers.UpdateClientHandler)
	}
}
