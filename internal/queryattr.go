package equinox

import "fmt"

type QueryAttr interface {
	// returns true if the specified attributes match this filter
	Match(attrs map[string]string) bool

	// returns human-readable string representation of the query
	String() string
}

/*
	not
	and
	or

	cmp: eq, regex, contains
*/

type QAEqual struct {
	k string
	v string
}

func (qa *QAEqual) Match(attrs map[string]string) bool {
	for k, v := range attrs {
		if k == qa.k {
			return v == qa.v
		}
	}

	// if the attribute doesn't exist then it's not a match
	return false
}

func (qa *QAEqual) String() string {
	return fmt.Sprintf("%s == '%s'", qa.k, qa.v)
}

func NewQAEqual(k, v string) *QAEqual {
	qa := QAEqual{k: k, v: v}
	return &qa
}
