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

/****************************************************************************
	QACmp
****************************************************************************/

type QAOp int64

const (
	Equal QAOp = iota
	Regex
)

func (s QAOp) String() string {
	switch s {
	case Equal:
		return "=="
	case Regex:
		return "=~"
	}
	return "???"
}

type QACmp struct {
	k  string
	v  string
	op QAOp
}

func (qa *QACmp) Match(attrs map[string]string) bool {
	for k, v := range attrs {
		if k == qa.k {
			switch qa.op {
			case Equal:
				return v == qa.v
			case Regex:
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

func NewQACmp(k string, v string, op QAOp) *QACmp {
	qa := QACmp{k: k, v: v, op: op}
	return &qa
}
