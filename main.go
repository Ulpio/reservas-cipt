package main

import (
	"fmt"

	"github.com/Ulpio/reservas-cipt/config"
	"github.com/Ulpio/reservas-cipt/database"
	"github.com/Ulpio/reservas-cipt/routes"
	"github.com/Ulpio/reservas-cipt/services"
)

func main() {
	config.LoadEnv()
	database.ConnectDB()
	InsertAdmin()
	routes.SetupRoutes()
}

func InsertAdmin() {
	_, err := services.CreateAdmin("Ulpio Paulo de Miranda Netto", "13366671416")
	if err != nil {
		return
	}
	fmt.Println("Admin criado com sucesso")
}
