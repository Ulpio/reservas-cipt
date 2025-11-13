package errors

import (
	"github.com/gin-gonic/gin"
)

// SendError envia uma resposta de erro padronizada
func SendError(c *gin.Context, err *APIError) {
	c.JSON(err.StatusCode, gin.H{
		"status_code": err.StatusCode,
		"message":     err.Message,
		"code":        err.Code,
		"field":       err.Field,
		"details":     err.Details,
	})
}

// SendErrorWithMessage envia um erro com mensagem customizada
func SendErrorWithMessage(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"status_code": statusCode,
		"message":     message,
	})
}

// SendValidationError envia um erro de validação com campo específico
func SendValidationError(c *gin.Context, field, message string) {
	c.JSON(400, gin.H{
		"status_code": 400,
		"message":     message,
		"field":       field,
		"code":        "VALIDATION_ERROR",
	})
}

