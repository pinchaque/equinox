package ctl

import (
	"equinox/internal/core"
	"net/http"
	"time"

	"equinox/internal/mw"

	"github.com/gin-gonic/gin"
)

func PointAdd(c *gin.Context) {
	// get the data series
	sid := c.Param("id")
	s, err := mw.GetSeriesMgr().Get(sid)
	if err != nil {
		c.JSON(http.StatusBadRequest, mw.Error(err.Error()))
		return
	}

	p := core.NewPoint(time.Now().UTC())
	err = c.BindJSON(p)
	if err != nil {
		c.JSON(http.StatusBadRequest, mw.Error(err.Error()))
		return
	}

	// save the point
	err = s.IO.Add([]*core.Point{p})
	if err != nil {
		c.JSON(http.StatusBadRequest, mw.Error(err.Error()))
		return
	}

	// TODO: should be returning "p" as part of a map that looks like: { "point" : {... p ... }}
	ret := gin.H{"point": p}
	c.JSON(http.StatusCreated, mw.Success(ret))
}
