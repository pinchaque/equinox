package mw

// Helper functions that make it easy to return responses compliant with JSend
// as described here: https://github.com/omniti-labs/jsend

// All went well, and (usually) some data was returned.
func Success(data any) map[string]any {
	h := make(map[string]any, 2)
	h["status"] = "success"
	h["data"] = data
	return h
}

// There was a problem with the data submitted, or some pre-condition of the
// API call wasn't satisfied
func Fail(data any) map[string]any {
	h := make(map[string]any, 2)
	h["status"] = "fail"
	h["data"] = data
	return h
}

// An error occurred in processing the request, i.e. an exception was thrown
func Error(msg string) map[string]any {
	h := make(map[string]any, 2)
	h["status"] = "error"
	h["message"] = msg
	return h
}
