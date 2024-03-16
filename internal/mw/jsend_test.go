package mw

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// sample data structure (that looks similar to a Point)
type testJsendData struct {
	Idx  int
	Attr map[string]string
	Vals map[string]int
	Ts   time.Time
	Name string
}

// gets a complex data structure
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
	exp := `{"status":"success","data":"basic string"}`
	assert.Equal(t, exp, string(b))

	j := NewJSend()
	json.Unmarshal(b, j)
	assert.Equal(t, r.Status, j.Status)
	assert.Equal(t, r.Data, j.Data)
	assert.Equal(t, r.Message, j.Message)
	assert.Equal(t, r.Code, j.Code)

	assert.True(t, j.IsSuccess())
	assert.False(t, j.IsFail())
	assert.False(t, j.IsError())
}

func TestJSendFail(t *testing.T) {
	r := Fail("basic string")
	b, err := json.Marshal(r)
	assert.NoError(t, err)
	exp := `{"status":"fail","data":"basic string"}`
	assert.Equal(t, exp, string(b))

	j := NewJSend()
	json.Unmarshal(b, j)
	assert.Equal(t, r.Status, j.Status)
	assert.Equal(t, r.Data, j.Data)
	assert.Equal(t, r.Message, j.Message)
	assert.Equal(t, r.Code, j.Code)

	assert.False(t, j.IsSuccess())
	assert.True(t, j.IsFail())
	assert.False(t, j.IsError())
}
func TestJSendSuccess2(t *testing.T) {
	// first we marshal the whole object to JSON
	r := Success(testGetJSendData())
	b, err := json.Marshal(r)
	assert.NoError(t, err)
	exp_data := `{"Idx":345,"Attr":null,"Vals":{"a":4,"b":982},"Ts":"2024-01-10T23:01:02.123456789Z","Name":"Best Struct Ever"}`
	exp := `{"status":"success","data":` + exp_data + `}`
	assert.Equal(t, exp, string(b))

	// now unmarshal it and make sure it matches
	j := NewJSend()
	json.Unmarshal(b, j)
	assert.Equal(t, r.Status, j.Status)
	assert.Equal(t, exp_data, string(j.Data))
	assert.Equal(t, r.Message, j.Message)
	assert.Equal(t, r.Code, j.Code)
}

func TestJSendError(t *testing.T) {
	r := Error("cool error message")
	b, err := json.Marshal(r)
	assert.NoError(t, err)
	exp := `{"status":"error","message":"cool error message"}`
	assert.Equal(t, exp, string(b))

	j := NewJSend()
	json.Unmarshal(b, j)
	assert.Equal(t, r.Status, j.Status)
	assert.Equal(t, r.Data, j.Data)
	assert.Equal(t, r.Message, j.Message)
	assert.Equal(t, r.Code, j.Code)

	assert.False(t, j.IsSuccess())
	assert.False(t, j.IsFail())
	assert.True(t, j.IsError())
}
