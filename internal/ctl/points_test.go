package ctl

import (
	"bytes"
	"encoding/json"
	"equinox/internal/core"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func testNewPoint() *core.Point {
	ts := time.Date(2024, 01, 10, 23, 1, 2, 0, time.UTC)
	p := core.NewPoint(ts)
	p.Attrs["shape"] = "square"
	p.Attrs["color"] = "red"
	p.Vals["area"] = 43.1
	p.Vals["temp"] = 21.1
	return p
}

func TestPointsAdd(t *testing.T) {
	router := gin.Default()
	router.POST("/points", PointAdd)

	p := testNewPoint()
	data, err := json.Marshal(p)
	assert.NoError(t, err)
	buf := bytes.NewReader(data)

	req, err := http.NewRequest("POST", "/points", buf)
	assert.NoError(t, err)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "{\"message\":\"Hello World\"}", rec.Body.String())

}
