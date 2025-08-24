package routes

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes() {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	SetupAuthRoutes(api)
	SetupUserRoutes(api)
	SpaceRoutes(api)
	SetupClientRoutes(api)
	SetupReservationRoutes(api)
	SetupStrikeRoutes(api)

	r.Run(":8080")
}
