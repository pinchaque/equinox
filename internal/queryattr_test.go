package equinox

import (
	"sort"
	"strings"
	"testing"
)

func getAttrs() map[string]string {
	r := make(map[string]string)
	r["color"] = "blue"
	r["animal"] = "moose"
	r["shape"] = "square"
	r["index"] = "74"
	return r
}

func attrsToString(attrs map[string]string) string {
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
			attrsToString(attrs), qa.String(), exp, result)
	}
}

func TestAttrString(t *testing.T) {

	fn := func(t *testing.T, q QueryAttr, exp string) {
		if q.String() != exp {
			t.Errorf("operator string: expected %s got %s", exp, q.String())
		}
	}

	fn(t, NewQAEqual("color", "blue"), "color == 'blue'")
	fn(t, NewQARegex("color", "blue"), "color =~ /blue/")
}

func TestAttrEqual(t *testing.T) {
	a := getAttrs()
	runQATest(t, a, NewQAEqual("color", "blue"), true)
	runQATest(t, a, NewQAEqual("color", "red"), false)
	runQATest(t, a, NewQAEqual("flavor", "sour"), false) // missing
	runQATest(t, a, NewQAEqual("color", "Blue"), false)  // case sensitive
	runQATest(t, a, NewQAEqual("Color", "blue"), false)  // case sensitive
	runQATest(t, a, NewQAEqual("color", "blu"), false)   // contains
}

func TestAttrRegex(t *testing.T) {
	a := getAttrs()
	runQATest(t, a, NewQARegex("color", "blue"), true)
	runQATest(t, a, NewQARegex("color", "red"), false)
	runQATest(t, a, NewQARegex("flavor", "sour"), false) // missing
	runQATest(t, a, NewQARegex("color", "Blue"), false)  // case sensitive
	runQATest(t, a, NewQARegex("Color", "blue"), false)  // case sensitive
	runQATest(t, a, NewQARegex("color", "blu"), true)    // contains

	// now we get into regexp metacharacters
	runQATest(t, a, NewQARegex("color", "blu."), true)
	runQATest(t, a, NewQARegex("color", "^bl"), true)
	runQATest(t, a, NewQARegex("color", "ue$"), true)
	runQATest(t, a, NewQARegex("color", "^x"), false)
	runQATest(t, a, NewQARegex("color", "blu$"), false)

	runQATest(t, a, NewQARegex("animal", "mo+se"), true)
	runQATest(t, a, NewQARegex("animal", "mo.se"), true)
	runQATest(t, a, NewQARegex("animal", "mo*se"), true)
	runQATest(t, a, NewQARegex("animal", "mo{1,5}se"), true)
	runQATest(t, a, NewQARegex("animal", "mo{3,5}se"), false)
	runQATest(t, a, NewQARegex("index", `^\d+$`), true)
	runQATest(t, a, NewQARegex("index", `^\D+$`), false)
}
