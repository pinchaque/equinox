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

func TestAttrString(t *testing.T) {

	fn := func(t *testing.T, q *QACmp, exp string) {
		if q.String() != exp {
			t.Errorf("operator string: expected %s got %s", exp, q.String())
		}
	}

	fn(t, NewQACmp("color", "blue", Equal), "color == 'blue'")
	fn(t, NewQACmp("color", "blue", Regex), "color =~ 'blue'")
}

func TestAttrEqual(t *testing.T) {
	a := getAttrs()
	op := Equal
	runQATest(t, a, NewQACmp("color", "blue", op), true)
	runQATest(t, a, NewQACmp("color", "red", op), false)
	runQATest(t, a, NewQACmp("flavor", "sour", op), false) // missing
	runQATest(t, a, NewQACmp("color", "Blue", op), false)  // case sensitive
	runQATest(t, a, NewQACmp("Color", "blue", op), false)  // case sensitive
	runQATest(t, a, NewQACmp("color", "blu", op), false)   // contains
}
