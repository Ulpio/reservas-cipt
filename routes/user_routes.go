package routes

import (
	"github.com/Ulpio/reservas-cipt/handlers"
	"github.com/Ulpio/reservas-cipt/middleware"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine) {
	users := r.Group("/users")
	users.Use(middleware.JWTAuthMiddleware())
	users.GET("/", middleware.OnlyAdmin(), handlers.GetAllUsersHandler)
	users.GET("/me", handlers.MeHandler)
	users.POST("/", middleware.OnlyAdmin(), handlers.CreateUserHandler)
	users.GET("/:id", handlers.GetUserByIDHandler)
	users.DELETE("/:id", middleware.OnlyAdmin(), handlers.DeleteUserHandler)
}
