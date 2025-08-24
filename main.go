// @title Reserva de Espaços - CIPT Jaraguá
// @version 1.0
// @description API interna para gestão de reservas e controle de espaços físicos no Centro de Inovação do Polo Tecnológico - Jaraguá, Maceió.
// @contact.name Ulpio Netto
// @contact.email oxetech.mcz@gmail.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /api/v1
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
