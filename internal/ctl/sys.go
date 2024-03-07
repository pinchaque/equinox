package ctl

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	resp := make(map[string]string, 1)
	resp["message"] = "Hello World"
	c.JSON(http.StatusOK, resp)
}
