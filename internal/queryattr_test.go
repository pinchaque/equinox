package equinox

import (
	"sort"
	"strings"
	"testing"
)

func getAttrs() map[string]string {
	r := make(map[string]string)
	r["color"] = "blue"
	r["animal"] = "cat"
	r["shape"] = "square"
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

func TestAttrEqual(t *testing.T) {
	a := getAttrs()
	runQATest(t, a, NewQAEqual("color", "blue"), true)
	runQATest(t, a, NewQAEqual("color", "red"), false)
	runQATest(t, a, NewQAEqual("flavor", "sour"), false) // missing
	runQATest(t, a, NewQAEqual("color", "Blue"), false)  // case sensitive
	runQATest(t, a, NewQAEqual("Color", "blue"), false)  // case sensitive
	runQATest(t, a, NewQAEqual("color", "blu"), false)   // contains
}
