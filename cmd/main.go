package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ping(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "")
}

func main() {
	router := gin.Default()
	router.GET("/ping", ping)

	router.Run("localhost:8080")
}
