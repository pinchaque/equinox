package query

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
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

// Helper data struct that is used to marshal/unmarshal the JSON representation
// of FilterAttrs.
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

// This is the generic interface for querying against attributes that is used
// within the Query object as well as in composite attribute queries.
type FilterAttr interface {
	// Returns true if the specified attributes match this filter
	Match(attrs map[string]string) bool

	// Human-readable string representation of the query
	String() string

	// Implements TextMarshaler interface
	MarshalText() ([]byte, error)

	// Unmarshals this object from a FilterAttrJson struct, returning an error
	// if the JSON struct doesn't have the correct fields.
	unmarshalStruct(j *FilterAttrJson) error
}

// Helper function that creates a slice of json.RawMessage attributes from a
// slice of FilterAttrs
func createExprs(fa ...FilterAttr) ([]json.RawMessage, error) {
	exprs := make([]json.RawMessage, 0, len(fa))
	for _, v := range fa {
		var r json.RawMessage
		var err error
		r, err = v.MarshalText()
		if err != nil {
			return []json.RawMessage{}, err
		}
		exprs = append(exprs, r)
	}
	return exprs, nil
}

// Factory function that creates the correct FilterAttr object for the
// specified FilterOp.
func createFilterAttr(op FilterOp) (FilterAttr, error) {

	switch op {
	case OpTrue:
		return &FATrue{}, nil
	case OpExists:
		return &FAExists{}, nil
	case OpEqual:
		return &FAEqual{}, nil
	case OpRegex:
		return &FARegex{}, nil
	case OpAnd:
		return &FAAnd{}, nil
	case OpOr:
		return &FAOr{}, nil
	case OpNot:
		return &FANot{}, nil

	default:
		return nil, fmt.Errorf("unrecognized filter operator %s", op)
	}
}

// Unmarshals the specified JSON string into a hierarchy of FilterAttr objects.
// We can't use the standard UnmarshalText() functions directly because we don't
// know what type of FilterAttr object to instantiate.
func UnmarshalFilterAttr(text []byte) (FilterAttr, error) {
	j := newFilterAttrJson()
	err := json.Unmarshal(text, j)
	if err != nil {
		return nil, err
	}

	fa, err := createFilterAttr(j.Op)
	if err != nil {
		return nil, err
	}

	err = fa.unmarshalStruct(j)
	if err != nil {
		return nil, err
	}

	return fa, nil
}

/****************************************************************************
	FATrue
****************************************************************************/

// Always returns true; useful as a no-op
type FATrue struct{}

// Always returns true
func (fa *FATrue) Match(attrs map[string]string) bool { return true }

func (fa *FATrue) String() string { return "true" }

// Implements TextMarshaler interface
func (fa *FATrue) MarshalText() ([]byte, error) {
	return json.Marshal(FilterAttrJson{Op: OpTrue})
}

// Unmarshals this object from a FilterAttrJson struct, returning an error
// if the JSON struct doesn't have the correct fields.
func (fa *FATrue) unmarshalStruct(j *FilterAttrJson) error {
	s := "Invalid JSON for FATrue"
	if j.Op != OpTrue {
		return fmt.Errorf("%s: Op must be %s", s, OpTrue)
	}
	if len(j.Exprs) != 0 {
		return fmt.Errorf("%s: Exprs must be empty", s)
	}
	if j.Attr != "" {
		return fmt.Errorf("%s: Attr must be empty", s)
	}
	if j.Val != "" {
		return fmt.Errorf("%s: Val must be empty", s)
	}

	return nil
}

// Returns new QATrue object
func True() *FATrue {
	return &FATrue{}
}

/****************************************************************************
	FAExist - Exists
****************************************************************************/

// Represents whether an attribute exists
type FAExists struct {
	k string
}

// Returns true if this query attribute exists (can be empty)
func (fa *FAExists) Match(attrs map[string]string) bool {
	_, exists := attrs[fa.k]
	return exists
}

func (fa *FAExists) String() string {
	return fmt.Sprintf("%s exists", fa.k)
}

// Implements TextMarshaler interface
func (fa *FAExists) MarshalText() ([]byte, error) {
	return json.Marshal(FilterAttrJson{Op: OpExists, Attr: fa.k})
}

// Unmarshals this object from a FilterAttrJson struct, returning an error
// if the JSON struct doesn't have the correct fields.
func (fa *FAExists) unmarshalStruct(j *FilterAttrJson) error {
	s := "Invalid JSON for FAExists"
	if j.Op != OpExists {
		return fmt.Errorf("%s: Op must be %s", s, OpExists)
	}
	if len(j.Exprs) != 0 {
		return fmt.Errorf("%s: Exprs must be empty", s)
	}
	if j.Attr == "" {
		return fmt.Errorf("%s: Attr cannot be empty", s)
	}
	if j.Val != "" {
		return fmt.Errorf("%s: Val must be empty", s)
	}

	fa.k = j.Attr
	return nil
}

// Returns new QAExists object with specified attribute key
func Exists(k string) *FAExists {
	return &FAExists{k: k}
}

/****************************************************************************
	FAEqual - Equality
****************************************************************************/

// Represents an attribute equality comparison for a given key
type FAEqual struct {
	k string
	v string
}

// Returns true if this query attribute comparison matches the given attribute
// map.
func (fa *FAEqual) Match(attrs map[string]string) bool {
	v, exists := attrs[fa.k]

	// if the attribute doesn't exist then it's not a match
	// otherwise must be exact match
	return exists && (v == fa.v)
}

func (fa *FAEqual) String() string {
	return fmt.Sprintf("%s == '%s'", fa.k, fa.v)
}

// Implements TextMarshaler interface
func (fa *FAEqual) MarshalText() ([]byte, error) {
	return json.Marshal(FilterAttrJson{Op: OpEqual, Attr: fa.k, Val: fa.v})
}

// Unmarshals this object from a FilterAttrJson struct, returning an error
// if the JSON struct doesn't have the correct fields.
func (fa *FAEqual) unmarshalStruct(j *FilterAttrJson) error {
	s := "Invalid JSON for FAEqual"
	if j.Op != OpEqual {
		return fmt.Errorf("%s: Op must be %s", s, OpEqual)
	}
	if len(j.Exprs) != 0 {
		return fmt.Errorf("%s: Exprs must be empty", s)
	}
	if j.Attr == "" {
		return fmt.Errorf("%s: Attr cannot be empty", s)
	}
	if j.Val == "" {
		return fmt.Errorf("%s: Val cannot be empty", s)
	}

	fa.k = j.Attr
	fa.v = j.Val
	return nil
}

// Returns new QAEqual object with specified attribute key and value
func Equal(k string, v string) *FAEqual {
	return &FAEqual{k: k, v: v}
}

/****************************************************************************
	FARegex - Regular Expression
****************************************************************************/

// Represents an attribute equality comparison for a given key
type FARegex struct {
	k  string
	re *regexp.Regexp
}

// Returns true if this query attribute comparison matches the given attribute
// map.
func (fa *FARegex) Match(attrs map[string]string) bool {
	v, exists := attrs[fa.k]

	// if the attribute doesn't exist then it's not a match
	// compare value to compiled regexp
	return exists && fa.re.MatchString(v)
}

func (fa *FARegex) String() string {
	return fmt.Sprintf("%s =~ /%s/", fa.k, fa.re.String())
}

// Implements TextMarshaler interface
func (fa *FARegex) MarshalText() ([]byte, error) {
	return json.Marshal(FilterAttrJson{Op: OpRegex, Attr: fa.k, Val: fa.re.String()})
}

// Unmarshals this object from a FilterAttrJson struct, returning an error
// if the JSON struct doesn't have the correct fields.
func (fa *FARegex) unmarshalStruct(j *FilterAttrJson) error {
	s := "Invalid JSON for FARegex"
	if j.Op != OpRegex {
		return fmt.Errorf("%s: Op must be %s", s, OpRegex)
	}
	if len(j.Exprs) != 0 {
		return fmt.Errorf("%s: Exprs must be empty", s)
	}
	if j.Attr == "" {
		return fmt.Errorf("%s: Attr cannot be empty", s)
	}
	if j.Val == "" {
		return fmt.Errorf("%s: Val cannot be empty", s)
	}

	fa.k = j.Attr
	re, err := regexp.Compile(j.Val)
	if err != nil {
		return err
	}
	fa.re = re
	return nil
}

// Returns new QARegex object with specified attribute key and regex to use
// when comparing against values.
func Regex(k string, regex string) *FARegex {
	re := regexp.MustCompile(regex)
	return &FARegex{k: k, re: re}
}

/****************************************************************************
	FANot - Logical NOT
****************************************************************************/

// Represents logical inversion (NOT)
type FANot struct {
	fa FilterAttr
}

// Returns logical inversion (NOT) of the contained QueryAttr
func (fa *FANot) Match(attrs map[string]string) bool {
	return !fa.fa.Match(attrs)
}

func (fa *FANot) String() string {
	return fmt.Sprintf("!(%s)", fa.fa.String())
}

// Implements TextMarshaler interface
func (fa *FANot) MarshalText() ([]byte, error) {
	exprs, err := createExprs(fa.fa)
	if err != nil {
		return []byte(""), err
	}

	return json.Marshal(FilterAttrJson{Op: OpNot, Exprs: exprs})
}

// Unmarshals this object from a FilterAttrJson struct, returning an error
// if the JSON struct doesn't have the correct fields.
func (fa *FANot) unmarshalStruct(j *FilterAttrJson) error {
	s := "Invalid JSON for FANot"
	if j.Op != OpNot {
		return fmt.Errorf("%s: Op must be %s", s, OpNot)
	}
	if len(j.Exprs) != 1 {
		return fmt.Errorf("%s: Must have a single Exprs", s)
	}
	if j.Attr != "" {
		return fmt.Errorf("%s: Attr must be empty", s)
	}
	if j.Val != "" {
		return fmt.Errorf("%s: Val must be empty", s)
	}

	for _, expr := range j.Exprs {
		f, err := UnmarshalFilterAttr(expr)
		if err != nil {
			return err
		}
		fa.fa = f
	}

	return nil
}

// Returns new QANot object that's the logical inversion of the specified
// QueryAttr
func Not(fa FilterAttr) *FANot {
	return &FANot{fa: fa}
}

/****************************************************************************
	FAAnd - Logical AND
****************************************************************************/

// Represents logical conjunction (AND)
type FAAnd struct {
	fa []FilterAttr
}

// Returns logical conjunction (AND) of the contained QueryAttrs
func (fa *FAAnd) Match(attrs map[string]string) bool {
	if len(fa.fa) == 0 {
		return false
	}

	for i := 0; i < len(fa.fa); i++ {
		if !fa.fa[i].Match(attrs) {
			return false
		}
	}

	return true
}

func (fa *FAAnd) String() string {
	var ret []string

	for i := 0; i < len(fa.fa); i++ {
		ret = append(ret, "("+fa.fa[i].String()+")")
	}

	return strings.Join(ret, " && ")
}

// Implements TextMarshaler interface
func (fa *FAAnd) MarshalText() ([]byte, error) {
	exprs, err := createExprs(fa.fa...)
	if err != nil {
		return []byte(""), err
	}

	return json.Marshal(FilterAttrJson{Op: OpAnd, Exprs: exprs})
}

// Unmarshals this object from a FilterAttrJson struct, returning an error
// if the JSON struct doesn't have the correct fields.
func (fa *FAAnd) unmarshalStruct(j *FilterAttrJson) error {
	s := "Invalid JSON for FAAnd"
	if j.Op != OpAnd {
		return fmt.Errorf("%s: Op must be %s", s, OpAnd)
	}
	if len(j.Exprs) == 0 {
		return fmt.Errorf("%s: Must have at least 1 Exprs", s)
	}
	if j.Attr != "" {
		return fmt.Errorf("%s: Attr must be empty", s)
	}
	if j.Val != "" {
		return fmt.Errorf("%s: Val must be empty", s)
	}

	fas := make([]FilterAttr, 0, len(j.Exprs))
	for _, expr := range j.Exprs {
		f, err := UnmarshalFilterAttr(expr)
		if err != nil {
			return err
		}
		fas = append(fas, f)
	}
	fa.fa = fas

	return nil
}

// Returns new QAAnd object that's the logical inversion of the specified
// QueryAttr
func And(fa ...FilterAttr) *FAAnd {
	return &FAAnd{fa: fa}
}

/****************************************************************************
	FAOr - Logical OR
****************************************************************************/

// Represents logical disjunction (OR)
type FAOr struct {
	fa []FilterAttr
}

// Returns logical disjunction (OR) of the contained QueryAttrs
func (fa *FAOr) Match(attrs map[string]string) bool {
	if len(fa.fa) == 0 {
		return false
	}

	for i := 0; i < len(fa.fa); i++ {
		if fa.fa[i].Match(attrs) {
			return true
		}
	}

	return false
}

func (fa *FAOr) String() string {
	var ret []string

	for i := 0; i < len(fa.fa); i++ {
		ret = append(ret, "("+fa.fa[i].String()+")")
	}

	return strings.Join(ret, " || ")
}

// Implements TextMarshaler interface
func (fa *FAOr) MarshalText() ([]byte, error) {
	exprs, err := createExprs(fa.fa...)
	if err != nil {
		return []byte(""), err
	}

	return json.Marshal(FilterAttrJson{Op: OpOr, Exprs: exprs})
}

// Unmarshals this object from a FilterAttrJson struct, returning an error
// if the JSON struct doesn't have the correct fields.
func (fa *FAOr) unmarshalStruct(j *FilterAttrJson) error {
	s := "Invalid JSON for FAOr"
	if j.Op != OpOr {
		return fmt.Errorf("%s: Op must be %s", s, OpOr)
	}
	if len(j.Exprs) == 0 {
		return fmt.Errorf("%s: Must have at least 1 Exprs", s)
	}
	if j.Attr != "" {
		return fmt.Errorf("%s: Attr must be empty", s)
	}
	if j.Val != "" {
		return fmt.Errorf("%s: Val must be empty", s)
	}

	fas := make([]FilterAttr, 0, len(j.Exprs))
	for _, expr := range j.Exprs {
		f, err := UnmarshalFilterAttr(expr)
		if err != nil {
			return err
		}
		fas = append(fas, f)
	}
	fa.fa = fas

	return nil
}

// Returns new QAOr object that's the logical inversion of the specified
// QueryAttr
func Or(fa ...FilterAttr) *FAOr {
	return &FAOr{fa: fa}
}
