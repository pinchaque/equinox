package equinox

import (
	"fmt"
	"regexp"
)

// This is the generic interface for querying against attributes that is used
// within the Query object as well as in composite attribute queries.
type QueryAttr interface {
	// Returns true if the specified attributes match this filter
	Match(attrs map[string]string) bool

	// Human-readable string representation of the query
	String() string
}

/*
	TODO: not and or
*/

/****************************************************************************
	QAEqual - Equality
****************************************************************************/

// Represents an attribute equality comparison for a given key
type QAEqual struct {
	k string
	v string
}

// Returns true if this query attribute comparison matches the given attribute
// map.
func (qa *QAEqual) Match(attrs map[string]string) bool {
	v, exists := attrs[qa.k]

	// if the attribute doesn't exist then it's not a match
	// otherwise must be exact match
	return exists && (v == qa.v)
}

func (qa *QAEqual) String() string {
	return fmt.Sprintf("%s == '%s'", qa.k, qa.v)
}

// Returns new QAEqual object with specified attribute key and value
func NewQAEqual(k string, v string) *QAEqual {
	return &QAEqual{k: k, v: v}
}

/****************************************************************************
	QARegex - Regular Expression
****************************************************************************/

// Represents an attribute equality comparison for a given key
type QARegex struct {
	k  string
	re *regexp.Regexp
}

// Returns true if this query attribute comparison matches the given attribute
// map.
func (qa *QARegex) Match(attrs map[string]string) bool {
	v, exists := attrs[qa.k]

	// if the attribute doesn't exist then it's not a match
	// compare value to compiled regexp
	return exists && qa.re.MatchString(v)
}

func (qa *QARegex) String() string {
	return fmt.Sprintf("%s =~ /%s/", qa.k, qa.re.String())
}

// Returns new QARegex object with specified attribute key and regex to use
// when comparing against values.
func NewQARegex(k string, regex string) *QARegex {
	re := regexp.MustCompile(regex)
	return &QARegex{k: k, re: re}
}
