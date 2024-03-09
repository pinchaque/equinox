package mw

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testJsendData struct {
	Idx  int
	Attr map[string]string
	Vals map[string]int
	Ts   time.Time
	Name string
}

func testGetJSendData() testJsendData {
	r := testJsendData{}
	r.Idx = 345
	r.Ts = time.Date(2024, 01, 10, 23, 1, 2, 123456789, time.UTC) // add microsecs
	r.Name = "Best Struct Ever"
	r.Vals = make(map[string]int)
	r.Vals["a"] = 4
	r.Vals["b"] = 982
	return r
}

func TestJSendSuccess(t *testing.T) {
	r := Success("basic string")
	b, err := json.Marshal(r)
	assert.NoError(t, err)
	exp := `{"data":"basic string","status":"success"}`
	assert.Equal(t, exp, string(b))
}

func TestJSendFail(t *testing.T) {
	r := Fail("basic string")
	b, err := json.Marshal(r)
	assert.NoError(t, err)
	exp := `{"data":"basic string","status":"fail"}`
	assert.Equal(t, exp, string(b))
}
func TestJSendSuccess2(t *testing.T) {
	r := Success(testGetJSendData())
	b, err := json.Marshal(r)
	assert.NoError(t, err)
	exp := `{"data":{"Idx":345,"Attr":null,"Vals":{"a":4,"b":982},"Ts":"2024-01-10T23:01:02.123456789Z","Name":"Best Struct Ever"},"status":"success"}`
	assert.Equal(t, exp, string(b))
}
