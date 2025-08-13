package routes

import (
	"github.com/Ulpio/reservas-cipt/handlers"
	"github.com/Ulpio/reservas-cipt/middleware"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine) {
	users := r.Group("/users")
	users.GET("/", middleware.JWTAuthMiddleware(), handlers.GetAllUsersHandler)
	users.GET("/me", middleware.JWTAuthMiddleware(), handlers.MeHandler)
	users.POST("/", middleware.JWTAuthMiddleware(), handlers.CreateUserHandler)
	users.GET("/:id", middleware.JWTAuthMiddleware(), handlers.GetUserByIDHandler)
	users.DELETE("/:id", middleware.JWTAuthMiddleware(), handlers.DeleteUserHandler)
}
