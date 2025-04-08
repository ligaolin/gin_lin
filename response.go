package gin_lin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, message string, data any) {
	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"code":    200,
		"data":    data,
	})
}

func Error(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"code":    400,
		"data":    nil,
	})
}

func Response(c *gin.Context, code int, message string, data any) {
	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"code":    code,
		"data":    data,
	})
}
