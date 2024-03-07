package main_test

import (
	"equinox/internal/routers"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const host = "localhost"
const port = "8080"

func launchRouter(t *testing.T) {
	r := routers.SetupRouter()
	err := r.Run(host + ":" + port)
	assert.NoError(t, err)
}

func getResponseText(t *testing.T, path string) string {
	url := fmt.Sprintf("http://%s:%s%s", host, port, path)
	resp, err := http.Get(url)
	assert.NoError(t, err)

	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	b, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	return string(b)
}

func TestPing(t *testing.T) {
	go launchRouter(t)
	act := getResponseText(t, "/ping")
	exp := `{"message":"Hello World"}`
	assert.Equal(t, exp, act)
}
