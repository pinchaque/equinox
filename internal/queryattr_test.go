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

func TestAttrBasic(t *testing.T) {
	attrs := getAttrs()

	runQATest(t, attrs, NewQAEqual("color", "blue"), true)
	runQATest(t, attrs, NewQAEqual("color", "red"), false)
}
