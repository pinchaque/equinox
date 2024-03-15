package mw

import "encoding/json"

// Implementation of the JSend response format as described here:
// https://github.com/omniti-labs/jsend
type JSend struct {
	Status  string          `json:"status,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
	Message string          `json:"message,omitempty"`
	Code    string          `json:"code,omitempty"`
}

func NewJSend() *JSend {
	return &JSend{}
}

// All went well, and (usually) some data was returned.
func Success(data any) *JSend {
	var err error
	js := NewJSend()
	js.Status = "success"
	js.Data, err = json.Marshal(data)
	if err != nil {
		panic("failed to marshal data to JSON")
	}
	return js
}

// There was a problem with the data submitted, or some pre-condition of the
// API call wasn't satisfied
func Fail(data any) *JSend {
	var err error
	js := NewJSend()
	js.Status = "fail"
	js.Data, err = json.Marshal(data)
	if err != nil {
		panic("failed to marshal data to JSON")
	}
	return js
}

// An error occurred in processing the request, i.e. an exception was thrown
func Error(msg string) *JSend {
	js := NewJSend()
	js.Status = "error"
	js.Message = msg
	return js
}
