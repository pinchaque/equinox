package ctl

import (
	"equinox/internal/core"
	"net/http"
	"time"

	"equinox/internal/mw"

	"github.com/gin-gonic/gin"
)

func PointAdd(c *gin.Context) {
	p := core.NewPoint(time.Now().UTC())
	err := c.BindJSON(p)
	if err != nil {
		// TODO how to handle this error?
		c.JSON(http.StatusBadRequest, mw.Error(err.Error()))
	}

	// TODO: save the point

	// TODO: should be returning "p" as part of a map that looks like: { "point" : {... p ... }}
	c.JSON(http.StatusCreated, mw.Success(p))
}
