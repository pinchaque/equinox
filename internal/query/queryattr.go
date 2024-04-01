package query

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// This is the generic interface for querying against attributes that is used
// within the Query object as well as in composite attribute queries.
type QueryAttr interface {
	// Returns true if the specified attributes match this filter
	Match(attrs map[string]string) bool

	// Human-readable string representation of the query
	String() string

	// Implements TextMarshaler interface
	MarshalText() (text []byte, err error)

	// Implements TextUnmarshaler interface
	UnmarshalText(text []byte) error
}

/****************************************************************************
	QATrue
****************************************************************************/

// Always returns true; useful as a no-op
type QATrue struct{}

// Always returns true
func (qa *QATrue) Match(attrs map[string]string) bool { return true }

func (qa *QATrue) String() string { return "true" }

// TODO: Implements TextMarshaler interface
func (qa *QATrue) MarshalText() (text []byte, err error) {
	return []byte(""), nil
}

// TODO: Implements TextUnmarshaler interface
func (qa *QATrue) UnmarshalText(text []byte) error {
	return nil
}

// Returns new QATrue object
func True() *QATrue {
	return &QATrue{}
}

/****************************************************************************
	QAExist - Exists
****************************************************************************/

// Represents whether an attribute exists
type QAExists struct {
	k string
}

// Returns true if this query attribute exists (can be empty)
func (qa *QAExists) Match(attrs map[string]string) bool {
	_, exists := attrs[qa.k]
	return exists
}

func (qa *QAExists) String() string {
	return fmt.Sprintf("%s exists", qa.k)
}

type jsonExpr struct {
	Op    string            `json:"op"`
	Attr  string            `json:"attr,omitempty"`
	Val   string            `json:"val,omitempty"`
	Expr  json.RawMessage   `json:"expr,omitempty"`
	Exprs []json.RawMessage `json:"exprs,omitempty"`
}

// Implements TextMarshaler interface
func (qa *QAExists) MarshalText() (text []byte, err error) {
	return json.Marshal(jsonExpr{Op: "exists", Attr: qa.k})
}

// Implements TextUnmarshaler interface
func (qa *QAExists) UnmarshalText(text []byte) error {
	qa.k = string(text)
	return nil
}

// Returns new QAExists object with specified attribute key
func Exists(k string) *QAExists {
	return &QAExists{k: k}
}

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

// TODO: Implements TextMarshaler interface
func (qa *QAEqual) MarshalText() (text []byte, err error) {
	return []byte(""), nil
}

// TODO: Implements TextUnmarshaler interface
func (qa *QAEqual) UnmarshalText(text []byte) error {
	return nil
}

// Returns new QAEqual object with specified attribute key and value
func Equal(k string, v string) *QAEqual {
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

// TODO: Implements TextMarshaler interface
func (qa *QARegex) MarshalText() (text []byte, err error) {
	return []byte(""), nil
}

// TODO: Implements TextUnmarshaler interface
func (qa *QARegex) UnmarshalText(text []byte) error {
	return nil
}

// Returns new QARegex object with specified attribute key and regex to use
// when comparing against values.
func Regex(k string, regex string) *QARegex {
	re := regexp.MustCompile(regex)
	return &QARegex{k: k, re: re}
}

/****************************************************************************
	QANot - Logical NOT
****************************************************************************/

// Represents logical inversion (NOT)
type QANot struct {
	qa QueryAttr
}

// Returns logical inversion (NOT) of the contained QueryAttr
func (qa *QANot) Match(attrs map[string]string) bool {
	return !qa.qa.Match(attrs)
}

func (qa *QANot) String() string {
	return fmt.Sprintf("!(%s)", qa.qa.String())
}

// TODO: Implements TextMarshaler interface
func (qa *QANot) MarshalText() (text []byte, err error) {
	return []byte(""), nil
}

// TODO: Implements TextUnmarshaler interface
func (qa *QANot) UnmarshalText(text []byte) error {
	return nil
}

// Returns new QANot object that's the logical inversion of the specified
// QueryAttr
func Not(qa QueryAttr) *QANot {
	return &QANot{qa: qa}
}

/****************************************************************************
	QAAnd - Logical AND
****************************************************************************/

// Represents logical conjunction (AND)
type QAAnd struct {
	qa []QueryAttr
}

// Returns logical conjunction (AND) of the contained QueryAttrs
func (qa *QAAnd) Match(attrs map[string]string) bool {
	if len(qa.qa) == 0 {
		return false
	}

	for i := 0; i < len(qa.qa); i++ {
		if !qa.qa[i].Match(attrs) {
			return false
		}
	}

	return true
}

func (qa *QAAnd) String() string {
	var ret []string

	for i := 0; i < len(qa.qa); i++ {
		ret = append(ret, "("+qa.qa[i].String()+")")
	}

	return strings.Join(ret, " && ")
}

// TODO: Implements TextMarshaler interface
func (qa *QAAnd) MarshalText() (text []byte, err error) {
	return []byte(""), nil
}

// TODO: Implements TextUnmarshaler interface
func (qa *QAAnd) UnmarshalText(text []byte) error {
	return nil
}

// Returns new QAAnd object that's the logical inversion of the specified
// QueryAttr
func And(qa ...QueryAttr) *QAAnd {
	return &QAAnd{qa: qa}
}

/****************************************************************************
	QAOr - Logical OR
****************************************************************************/

// Represents logical disjunction (OR)
type QAOr struct {
	qa []QueryAttr
}

// Returns logical disjunction (OR) of the contained QueryAttrs
func (qa *QAOr) Match(attrs map[string]string) bool {
	if len(qa.qa) == 0 {
		return false
	}

	for i := 0; i < len(qa.qa); i++ {
		if qa.qa[i].Match(attrs) {
			return true
		}
	}

	return false
}

func (qa *QAOr) String() string {
	var ret []string

	for i := 0; i < len(qa.qa); i++ {
		ret = append(ret, "("+qa.qa[i].String()+")")
	}

	return strings.Join(ret, " || ")
}

// TODO: Implements TextMarshaler interface
func (qa *QAOr) MarshalText() (text []byte, err error) {
	return []byte(""), nil
}

// TODO: Implements TextUnmarshaler interface
func (qa *QAOr) UnmarshalText(text []byte) error {
	return nil
}

// Returns new QAOr object that's the logical inversion of the specified
// QueryAttr
func Or(qa ...QueryAttr) *QAOr {
	return &QAOr{qa: qa}
}
