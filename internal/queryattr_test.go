package equinox

import (
	"testing"
)

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

	fn(t, NewQATrue(), "true")
	fn(t, NewQAEqual("color", "blue"), "color == 'blue'")
	fn(t, NewQARegex("color", "blue"), "color =~ /blue/")
	fn(t, NewQAExists("color"), "color exists")

	t1 := NewQAEqual("color", "blue")
	t2 := NewQAEqual("animal", "moose")
	t3 := NewQAEqual("shape", "square")

	fn(t, NewQANot(t1), "!(color == 'blue')")
	fn(t, NewQAOr(t2, t3), "(animal == 'moose') || (shape == 'square')")
	fn(t, NewQAAnd(t2, t3), "(animal == 'moose') && (shape == 'square')")
	fn(t, NewQAAnd(t2, NewQANot(t3)), "(animal == 'moose') && (!(shape == 'square'))")
}

func TestAttrTrue(t *testing.T) {
	a := testGetAttrs()
	runQATest(t, a, NewQATrue(), true)
}

func TestAttrExists(t *testing.T) {
	a := testGetAttrs()
	runQATest(t, a, NewQAExists("color"), true)
	runQATest(t, a, NewQAExists("animal"), true)
	runQATest(t, a, NewQAExists("shape"), true)
	runQATest(t, a, NewQAExists("Color"), false) // case sensitive
	runQATest(t, a, NewQAExists("flavor"), false)
	runQATest(t, a, NewQAExists("hue"), false)
}

func TestAttrEqual(t *testing.T) {
	a := testGetAttrs()
	runQATest(t, a, NewQAEqual("color", "blue"), true)
	runQATest(t, a, NewQAEqual("color", "red"), false)
	runQATest(t, a, NewQAEqual("flavor", "sour"), false) // missing
	runQATest(t, a, NewQAEqual("color", "Blue"), false)  // case sensitive
	runQATest(t, a, NewQAEqual("Color", "blue"), false)  // case sensitive
	runQATest(t, a, NewQAEqual("color", "blu"), false)   // contains
}

func TestAttrRegex(t *testing.T) {
	a := testGetAttrs()
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

func TestAttrNot(t *testing.T) {
	a := testGetAttrs()

	// things that are true
	t1 := NewQAEqual("color", "blue")
	t2 := NewQAEqual("animal", "moose")
	t3 := NewQARegex("index", `^\d+$`)
	t4 := NewQAEqual("shape", "square")

	// things that are false
	f1 := NewQARegex("color", "^x")
	f2 := NewQARegex("animal", "mo{3,5}se")
	f3 := NewQAEqual("index", "777")
	f4 := NewQARegex("flavor", "sour")

	runQATest(t, a, NewQANot(t1), false)
	runQATest(t, a, NewQANot(t2), false)
	runQATest(t, a, NewQANot(t3), false)
	runQATest(t, a, NewQANot(t4), false)
	runQATest(t, a, NewQANot(f1), true)
	runQATest(t, a, NewQANot(f2), true)
	runQATest(t, a, NewQANot(f3), true)
	runQATest(t, a, NewQANot(f4), true)
}

func TestAttrAnd(t *testing.T) {
	a := testGetAttrs()

	// things that are true
	t1 := NewQAEqual("color", "blue")
	t2 := NewQAEqual("animal", "moose")
	t3 := NewQARegex("index", `^\d+$`)
	t4 := NewQAEqual("shape", "square")

	// things that are false
	f1 := NewQARegex("color", "^x")
	f2 := NewQARegex("animal", "mo{3,5}se")
	f3 := NewQAEqual("index", "777")
	f4 := NewQARegex("flavor", "sour")

	runQATest(t, a, NewQAAnd(), false)
	runQATest(t, a, NewQAAnd(t1), true)
	runQATest(t, a, NewQAAnd(t1, t2), true)
	runQATest(t, a, NewQAAnd(t1, t2, t3, t4), true)
	runQATest(t, a, NewQAAnd(f1), false)
	runQATest(t, a, NewQAAnd(t1, f1), false)
	runQATest(t, a, NewQAAnd(t1, t2, t3, t4, f1), false)
	runQATest(t, a, NewQAAnd(f1, f2, f3, f4), false)
}

func TestAttrOr(t *testing.T) {
	a := testGetAttrs()

	// things that are true
	t1 := NewQAEqual("color", "blue")
	t2 := NewQAEqual("animal", "moose")
	t3 := NewQARegex("index", `^\d+$`)
	t4 := NewQAEqual("shape", "square")

	// things that are false
	f1 := NewQARegex("color", "^x")
	f2 := NewQARegex("animal", "mo{3,5}se")
	f3 := NewQAEqual("index", "777")
	f4 := NewQARegex("flavor", "sour")

	runQATest(t, a, NewQAOr(), false)
	runQATest(t, a, NewQAOr(t1), true)
	runQATest(t, a, NewQAOr(t1, t2), true)
	runQATest(t, a, NewQAOr(t1, t2, t3, t4), true)
	runQATest(t, a, NewQAOr(f1), false)
	runQATest(t, a, NewQAOr(t1, f1), true)
	runQATest(t, a, NewQAOr(t1, t2, t3, t4, f1), true)
	runQATest(t, a, NewQAOr(f1, f2, f3, f4), false)
}
func TestAttrLogicCombo(t *testing.T) {
	a := testGetAttrs()

	// things that are true
	t1 := NewQAEqual("color", "blue")
	t2 := NewQAEqual("animal", "moose")
	t3 := NewQARegex("index", `^\d+$`)
	t4 := NewQAEqual("shape", "square")

	// things that are false
	f1 := NewQARegex("color", "^x")
	f2 := NewQARegex("animal", "mo{3,5}se")
	f3 := NewQAEqual("index", "777")
	f4 := NewQARegex("flavor", "sour")

	t5 := NewQAOr(t1, f1)
	t6 := NewQAOr(f2, f3, f4, t2, t3, t4)
	t7 := NewQAAnd(NewQANot(f1), t1)

	f5 := NewQANot(NewQAAnd(t1, t4))
	f6 := NewQAAnd(NewQANot(f4), f5)
	f7 := NewQAOr(f2, f4)

	runQATest(t, a, NewQAOr(f5, f6, f7, t5), true)
	runQATest(t, a, NewQAAnd(f5, f6, f7, t5), false)
	runQATest(t, a, NewQAOr(f5, f6, f7), false)
	runQATest(t, a, NewQAAnd(f5, f6, f7), false)
	runQATest(t, a, NewQAOr(t5, t6, t7), true)
	runQATest(t, a, NewQAAnd(t5, t6, t7), true)
}
