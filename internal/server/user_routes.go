package server

import (
	handler "github.com/Ulpio/reservas-cipt/internal/handler/user"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine, h *handler.UserHandler) {
	users := router.Group("/users")
	{
		users.POST("", h.RegisterUser)
		users.GET("/:id", h.GetUserByID)
		users.GET("", h.GetAllUsers)
		users.DELETE("/:id", h.DeleteByID)
	}
}
