package routes

import "github.com/gin-gonic/gin"

func SetupRoutes() {
	r := gin.Default()
	SetupUserRoutes(r)
	SetupAuthRoutes(r)
	SpaceRoutes(r)
	SetupClientRoutes(r)
	SetupReservationRoutes(r)
	SetupStrikeRoutes(r)
	r.Run(":8080")
}
