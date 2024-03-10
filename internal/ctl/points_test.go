package ctl_test

import (
	"bytes"
	"encoding/json"
	"equinox/internal/core"
	"equinox/internal/engine"
	"equinox/internal/models"
	"equinox/internal/mw"
	"equinox/internal/routers"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func testNewPoint() *core.Point {
	ts := time.Date(2024, 01, 10, 23, 1, 2, 0, time.UTC)
	p := core.NewPointEmptyId(ts)
	p.Attrs["shape"] = "square"
	p.Attrs["color"] = "red"
	p.Vals["area"] = 43.1
	p.Vals["temp"] = 21.1
	return p
}

func setupDataSeries(id string) {
	mgr := mw.GetSeriesMgr()
	s := &models.Series{Id: id, IO: engine.NewMemTree()}
	mgr.Add(s)
}

func TestPointsAdd(t *testing.T) {
	sid := "foobar"
	setupDataSeries(sid)
	router := routers.SetupRouter()

	// create the point
	p := testNewPoint()
	assert.Nil(t, p.Id) // shouldn't have an Id yet
	data, err := json.Marshal(p)
	assert.NoError(t, err)

	// save the point
	path := fmt.Sprintf("/series/%s/points", sid)
	req, err := http.NewRequest("POST", path, bytes.NewReader(data))
	assert.NoError(t, err)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	result := make(map[string]any)
	err = json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)

	for k, v := range result {
		assert.Equal(t, "xxx", fmt.Sprintf("`%s`: `%s`", k, v))
	}

	// round trip should deliver the same JSON we sent
	//exp := fmt.Sprintf(`{"data":{"point":%s},"status":"success"}`, data)
	//assert.Equal(t, exp, rec.Body.String())

}
