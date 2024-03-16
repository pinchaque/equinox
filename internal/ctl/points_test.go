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
	assert.Equal(t, http.StatusCreated, rec.Code) // obj was created

	// unmarshal the response into a JSend response
	var js mw.JSend
	err = json.Unmarshal(rec.Body.Bytes(), &js)
	assert.NoError(t, err)

	// validate the JSend including raw JSON data that came back
	assert.Equal(t, "success", js.Status)
	assert.True(t, js.IsSuccess())
	assert.Equal(t, "", js.Message)
	assert.Equal(t, "", js.Code)
	assert.Contains(t, string(js.Data), `{"point":`)
	assert.Contains(t, string(js.Data), `"Ts":"2024-01-10T23:01:02Z"`)
	assert.Contains(t, string(js.Data), `"Vals":{"area":43.1,"temp":21.1}`)
	assert.Contains(t, string(js.Data), `"Attrs":{"color":"red","shape":"square"}`)
	assert.NotContains(t, string(js.Data), `"Id":null`) // need an ID returned

	// convert the "point" into a real struct
	pmap := make(map[string]*core.Point)
	err = json.Unmarshal(js.Data, &pmap)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(pmap))

	for k, p2 := range pmap {
		assert.Equal(t, "point", k)

		assert.Equal(t, "2024-01-10 23:01:02 +0000 UTC", p2.Ts.String())
		assert.True(t, p.Equal(p2))
		assert.False(t, p.Identical(p2)) // shouldn't be identical - id is now set
		assert.NotNil(t, p2.Id)
	}
}

// TODO test for fail when including ID
// TODO test for usage of current time if not included
