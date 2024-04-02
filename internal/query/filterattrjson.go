package query

import (
	"encoding/json"
)

// Interface to how the filters match
type FilterMatch interface {
	Match(attrs map[string]string) bool
}

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

// This is the generic interface for querying against attributes that is used
// within the Query object as well as in composite attribute queries.
type FilterAttrJson struct {
	Op        FilterOp          `json:"op"`
	Attr      string            `json:"attr,omitempty"`
	Val       string            `json:"val,omitempty"`
	JsonExprs []json.RawMessage `json:"exprs,omitempty"`
	Exprs     []*FilterAttr     `json:"-"`
	MatchFn   FilterMatch       `json:"-"`
}

// Implements TextMarshaler interface
func (fa *FilterAttrJson) MarshalText() (text []byte, err error) {
	return []byte(""), nil

}

// Implements TextUnmarshaler interface
func (fa *FilterAttrJson) UnmarshalText(text []byte) error {
	return nil
}
