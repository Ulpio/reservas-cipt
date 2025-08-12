package main

import (
	"github.com/Ulpio/reservas-cipt/config"
	"github.com/Ulpio/reservas-cipt/database"
	"github.com/Ulpio/reservas-cipt/routes"
)

func main() {
	config.LoadEnv()
	database.ConnectDB()
	routes.SetupRoutes()
}
