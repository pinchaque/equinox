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

	// read the JSON into a point
	p := core.NewPointEmpty()
	err = c.BindJSON(p)
	if err != nil {
		c.JSON(http.StatusBadRequest, mw.Error(err.Error()))
		return
	}

	// data validation - id should be empty
	if p.Id != nil {
		// TODO this should probably be Fail() instead but I need to figure
		// out how to represent the error message
		c.JSON(http.StatusBadRequest, mw.Error("ID cannot be specified in the request"))
		return
	}
	p.GenerateId()

	// data validation - timestamp should be specified or else we just use "now"
	empty_ts := time.Time{}
	if p.Ts == empty_ts {
		p.Ts = time.Now().UTC()
	}
	// save the point
	err = s.IO.Add(p)
	if err != nil {
		c.JSON(http.StatusBadRequest, mw.Error(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, mw.Success(gin.H{"point": p}))
}
