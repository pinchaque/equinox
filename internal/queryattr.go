package equinox

import "fmt"

// This is the generic interface for querying against attributes that is used
// within the Query object as well as in composite attribute queries.
type QueryAttr interface {
	// Returns true if the specified attributes match this filter
	Match(attrs map[string]string) bool

	// Human-readable string representation of the query
	String() string
}

/*
	not
	and
	or
*/

// Represents an operator used for QACmp
type QAOp int64

const (
	QA_EQ QAOp = iota
	QA_REGEX
)

// String representation of the operator.
func (s QAOp) String() string {
	switch s {
	case QA_EQ:
		return "=="
	case QA_REGEX:
		return "=~"
	}
	return "???"
}

// Represents an attribute comparison for a given key using a specified
// operator.
type QACmp struct {
	k  string
	v  string
	op QAOp
}

// Returns true if this query attribute comparison matches the given attribute
// map.
func (qa *QACmp) Match(attrs map[string]string) bool {
	for k, v := range attrs {
		if k == qa.k {
			switch qa.op {
			case QA_EQ:
				return v == qa.v
			case QA_REGEX:
				return false
			default:
				return false // invalid operator
			}
		}
	}

	// if the attribute doesn't exist then it's not a match
	return false
}

func (qa *QACmp) String() string {
	return fmt.Sprintf("%s %s '%s'", qa.k, qa.op.String(), qa.v)
}

// Returns new QACmp object with specified attribute key, value, and comparison
// operator.
func NewQACmp(k string, v string, op QAOp) *QACmp {
	qa := QACmp{k: k, v: v, op: op}
	return &qa
}
