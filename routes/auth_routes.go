package routes

import (
	"github.com/Ulpio/reservas-cipt/handlers"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	auth.POST("/login", handlers.LoginHandler)
}
