package query

import (
	"encoding/json"
)

// Enum-style type to represent the allowed operators
type FilterOp string

const (
	OpTrue   FilterOp = "true"
	OpExists FilterOp = "exists"
	OpEqual  FilterOp = "equal"
	OpRegex  FilterOp = "regex"
	OpAnd    FilterOp = "and"
	OpOr     FilterOp = "or"
	OpNot    FilterOp = "not"
)

// Data struct that is used to marshal/unmarshal the JSON representation of
// FilterAttrs.
type FilterAttrJson struct {
	Op    FilterOp          `json:"op"`
	Attr  string            `json:"attr,omitempty"`
	Val   string            `json:"val,omitempty"`
	Exprs []json.RawMessage `json:"exprs,omitempty"`
}

func newFilterAttrJson() *FilterAttrJson {
	j := FilterAttrJson{}
	j.Exprs = make([]json.RawMessage, 0)
	return &j
}
