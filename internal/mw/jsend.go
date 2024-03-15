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

// Returns true if this JSend object represents a successful response.
func (js *JSend) IsSuccess() bool {
	return js.Status == "success"
}

// Returns true if this JSend object represents a failed response. This would
// mean there was a problem with the data submitted, or some pre-condition of
// the API call wasn't satisfied.
func (js *JSend) IsFail() bool {
	return js.Status == "fail"
}

// Returns true if this JSend object represents an error response. This would
// mean an error occurred in processing the request, i.e. an exception was
// thrown. If this is true then the caller can check Message to see details.
func (js *JSend) IsError() bool {
	return js.Status == "error"
}
