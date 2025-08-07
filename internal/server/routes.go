package server

import (
	handler "github.com/Ulpio/reservas-cipt/internal/handler/user"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, userHandler *handler.UserHandler) {
	RegisterUserRoutes(router, userHandler)
}
