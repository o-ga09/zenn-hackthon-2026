package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}
