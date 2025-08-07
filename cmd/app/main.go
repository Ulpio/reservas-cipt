package main

import (
	"log"

	"github.com/Ulpio/reservas-cipt/internal/config"
	models "github.com/Ulpio/reservas-cipt/internal/domain/user"
	handler "github.com/Ulpio/reservas-cipt/internal/handler/user"
	"github.com/Ulpio/reservas-cipt/internal/repository/postgres"
	"github.com/Ulpio/reservas-cipt/internal/server"
	"github.com/Ulpio/reservas-cipt/internal/usecase/user"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal(err.Error())
	}
	db.AutoMigrate(&models.User{})
	repo := postgres.NewUserRepository(db)
	usecase := user.NewUserUseCase(repo)
	handler := handler.NewHandler(usecase)
	router := gin.Default()
	server.SetupRoutes(router, handler)
	router.Run()

}
