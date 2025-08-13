package routes

import (
	"github.com/Ulpio/reservas-cipt/handlers"
	"github.com/Ulpio/reservas-cipt/middleware"
	"github.com/gin-gonic/gin"
)

func SpaceRoutes(r *gin.Engine) {
	spaceGroup := r.Group("/espacos")
	spaceGroup.Use(middleware.JWTAuthMiddleware())

	spaceGroup.GET("/", handlers.GetAllSpacesHandler)
	spaceGroup.GET("/:id", handlers.GetSpacesByIDHandler)

	spaceGroup.PATCH("/:id/status", handlers.UpdateSpaceStatusHandler)
	spaceGroup.PATCH("/:id/aviso", handlers.UpdateSpaceNoticeHandler)

	admin := spaceGroup.Group("/")
	admin.Use(middleware.OnlyAdmin())
	admin.POST("/", handlers.CreateSpaceHandler)
	admin.PUT("/:id", handlers.UpdateSpaceHandler)
	admin.DELETE("/:id", handlers.DeleteSpaceHandler)
}
