package routes

import "github.com/gin-gonic/gin"

func SetupRoutes() {
	r := gin.Default()
	r.Run(":8080")
}
