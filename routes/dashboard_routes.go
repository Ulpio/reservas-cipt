package routes

import (
	"github.com/Ulpio/reservas-cipt/handlers"
	"github.com/Ulpio/reservas-cipt/middleware"
	"github.com/gin-gonic/gin"
)

// SetupDashboardRoutes configura as rotas do dashboard
func SetupDashboardRoutes(r *gin.RouterGroup) {
	dashboard := r.Group("/dashboard")
	dashboard.Use(middleware.JWTAuthMiddleware())
	{
		dashboard.GET("/stats", handlers.GetDashboardStatsHandler)
	}
}

