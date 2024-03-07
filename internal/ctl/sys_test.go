package ctl

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSysPing(t *testing.T) {
	router := gin.Default()
	router.GET("/ping", Ping)
	req, err := http.NewRequest("GET", "/ping", nil)
	assert.NoError(t, err)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "{\"message\":\"Hello World\"}", rec.Body.String())
}
