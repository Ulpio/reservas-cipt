package routes

import (
	"github.com/Ulpio/reservas-cipt/handlers"
	"github.com/Ulpio/reservas-cipt/middleware"
	"github.com/gin-gonic/gin"
)

// SetupReservationRoutes configures routes for reservation operations.
func SetupReservationRoutes(r *gin.Engine) {
	group := r.Group("/reservas")
	group.Use(middleware.JWTAuthMiddleware(), middleware.OnlyReceptionist())
	{
		group.POST("/", handlers.CreateReservationHandler)
		group.GET("/", handlers.GetAllReservationsHandler)
		group.GET("/:id", handlers.GetReservationByIDHandler)
	}
}
