package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func OnlyAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Acesso permitido apenas para administradores"})
			return
		}
		c.Next()
	}
}

func OnlyReceptionist() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "recepcionista" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Acesso permitido apenas para recepcionistas"})
			return
		}
		c.Next()
	}
}
