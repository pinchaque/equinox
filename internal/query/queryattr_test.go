package query

import (
	"sort"
	"strings"
	"testing"
)

func testGetAttrs() map[string]string {
	r := make(map[string]string)
	r["color"] = "blue"
	r["animal"] = "moose"
	r["shape"] = "square"
	r["index"] = "74"
	return r
}

func testAttrsToString(attrs map[string]string) string {
	var attr []string

	for k, v := range attrs {
		attr = append(attr, k+": "+v)
	}
	sort.Strings(attr) // ensure consistent output

	return strings.Join(attr, ", ")
}

func runQATest(t *testing.T, attrs map[string]string, qa QueryAttr, exp bool) {
	result := qa.Match(attrs)

	if result != exp {
		t.Errorf("QueryAttr: attrs{{{%s}}} and query{{{%s}}} expected %t got %t",
			testAttrsToString(attrs), qa.String(), exp, result)
	}
}

func TestAttrString(t *testing.T) {

	fn := func(t *testing.T, q QueryAttr, exp string) {
		if q.String() != exp {
			t.Errorf("operator string: expected %s got %s", exp, q.String())
		}
	}

	fn(t, True(), "true")
	fn(t, Equal("color", "blue"), "color == 'blue'")
	fn(t, Regex("color", "blue"), "color =~ /blue/")
	fn(t, Exists("color"), "color exists")

	t1 := Equal("color", "blue")
	t2 := Equal("animal", "moose")
	t3 := Equal("shape", "square")

	fn(t, Not(t1), "!(color == 'blue')")
	fn(t, Or(t2, t3), "(animal == 'moose') || (shape == 'square')")
	fn(t, And(t2, t3), "(animal == 'moose') && (shape == 'square')")
	fn(t, And(t2, Not(t3)), "(animal == 'moose') && (!(shape == 'square'))")
}

func TestAttrTrue(t *testing.T) {
	a := testGetAttrs()
	runQATest(t, a, True(), true)
}

func TestAttrExists(t *testing.T) {
	a := testGetAttrs()
	runQATest(t, a, Exists("color"), true)
	runQATest(t, a, Exists("animal"), true)
	runQATest(t, a, Exists("shape"), true)
	runQATest(t, a, Exists("Color"), false) // case sensitive
	runQATest(t, a, Exists("flavor"), false)
	runQATest(t, a, Exists("hue"), false)
}

func TestAttrEqual(t *testing.T) {
	a := testGetAttrs()
	runQATest(t, a, Equal("color", "blue"), true)
	runQATest(t, a, Equal("color", "red"), false)
	runQATest(t, a, Equal("flavor", "sour"), false) // missing
	runQATest(t, a, Equal("color", "Blue"), false)  // case sensitive
	runQATest(t, a, Equal("Color", "blue"), false)  // case sensitive
	runQATest(t, a, Equal("color", "blu"), false)   // contains
}

func TestAttrRegex(t *testing.T) {
	a := testGetAttrs()
	runQATest(t, a, Regex("color", "blue"), true)
	runQATest(t, a, Regex("color", "red"), false)
	runQATest(t, a, Regex("flavor", "sour"), false) // missing
	runQATest(t, a, Regex("color", "Blue"), false)  // case sensitive
	runQATest(t, a, Regex("Color", "blue"), false)  // case sensitive
	runQATest(t, a, Regex("color", "blu"), true)    // contains

	// now we get into regexp metacharacters
	runQATest(t, a, Regex("color", "blu."), true)
	runQATest(t, a, Regex("color", "^bl"), true)
	runQATest(t, a, Regex("color", "ue$"), true)
	runQATest(t, a, Regex("color", "^x"), false)
	runQATest(t, a, Regex("color", "blu$"), false)

	runQATest(t, a, Regex("animal", "mo+se"), true)
	runQATest(t, a, Regex("animal", "mo.se"), true)
	runQATest(t, a, Regex("animal", "mo*se"), true)
	runQATest(t, a, Regex("animal", "mo{1,5}se"), true)
	runQATest(t, a, Regex("animal", "mo{3,5}se"), false)
	runQATest(t, a, Regex("index", `^\d+$`), true)
	runQATest(t, a, Regex("index", `^\D+$`), false)
}

func TestAttrNot(t *testing.T) {
	a := testGetAttrs()

	// things that are true
	t1 := Equal("color", "blue")
	t2 := Equal("animal", "moose")
	t3 := Regex("index", `^\d+$`)
	t4 := Equal("shape", "square")

	// things that are false
	f1 := Regex("color", "^x")
	f2 := Regex("animal", "mo{3,5}se")
	f3 := Equal("index", "777")
	f4 := Regex("flavor", "sour")

	runQATest(t, a, Not(t1), false)
	runQATest(t, a, Not(t2), false)
	runQATest(t, a, Not(t3), false)
	runQATest(t, a, Not(t4), false)
	runQATest(t, a, Not(f1), true)
	runQATest(t, a, Not(f2), true)
	runQATest(t, a, Not(f3), true)
	runQATest(t, a, Not(f4), true)
}

func TestAttrAnd(t *testing.T) {
	a := testGetAttrs()

	// things that are true
	t1 := Equal("color", "blue")
	t2 := Equal("animal", "moose")
	t3 := Regex("index", `^\d+$`)
	t4 := Equal("shape", "square")

	// things that are false
	f1 := Regex("color", "^x")
	f2 := Regex("animal", "mo{3,5}se")
	f3 := Equal("index", "777")
	f4 := Regex("flavor", "sour")

	runQATest(t, a, And(), false)
	runQATest(t, a, And(t1), true)
	runQATest(t, a, And(t1, t2), true)
	runQATest(t, a, And(t1, t2, t3, t4), true)
	runQATest(t, a, And(f1), false)
	runQATest(t, a, And(t1, f1), false)
	runQATest(t, a, And(t1, t2, t3, t4, f1), false)
	runQATest(t, a, And(f1, f2, f3, f4), false)
}

func TestAttrOr(t *testing.T) {
	a := testGetAttrs()

	// things that are true
	t1 := Equal("color", "blue")
	t2 := Equal("animal", "moose")
	t3 := Regex("index", `^\d+$`)
	t4 := Equal("shape", "square")

	// things that are false
	f1 := Regex("color", "^x")
	f2 := Regex("animal", "mo{3,5}se")
	f3 := Equal("index", "777")
	f4 := Regex("flavor", "sour")

	runQATest(t, a, Or(), false)
	runQATest(t, a, Or(t1), true)
	runQATest(t, a, Or(t1, t2), true)
	runQATest(t, a, Or(t1, t2, t3, t4), true)
	runQATest(t, a, Or(f1), false)
	runQATest(t, a, Or(t1, f1), true)
	runQATest(t, a, Or(t1, t2, t3, t4, f1), true)
	runQATest(t, a, Or(f1, f2, f3, f4), false)
}
func TestAttrLogicCombo(t *testing.T) {
	a := testGetAttrs()

	// things that are true
	t1 := Equal("color", "blue")
	t2 := Equal("animal", "moose")
	t3 := Regex("index", `^\d+$`)
	t4 := Equal("shape", "square")

	// things that are false
	f1 := Regex("color", "^x")
	f2 := Regex("animal", "mo{3,5}se")
	f3 := Equal("index", "777")
	f4 := Regex("flavor", "sour")

	t5 := Or(t1, f1)
	t6 := Or(f2, f3, f4, t2, t3, t4)
	t7 := And(Not(f1), t1)

	f5 := Not(And(t1, t4))
	f6 := And(Not(f4), f5)
	f7 := Or(f2, f4)

	runQATest(t, a, Or(f5, f6, f7, t5), true)
	runQATest(t, a, And(f5, f6, f7, t5), false)
	runQATest(t, a, Or(f5, f6, f7), false)
	runQATest(t, a, And(f5, f6, f7), false)
	runQATest(t, a, Or(t5, t6, t7), true)
	runQATest(t, a, And(t5, t6, t7), true)
}
