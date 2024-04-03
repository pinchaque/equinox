package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testGetAttrs() map[string]string {
	r := make(map[string]string)
	r["color"] = "blue"
	r["animal"] = "moose"
	r["shape"] = "square"
	r["index"] = "74"
	return r
}

func TestAttrString(t *testing.T) {

	fn := func(t *testing.T, q FilterAttr, exp string) {
		assert.Equal(t, exp, q.String())
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

func runQATest(t *testing.T, attrs map[string]string, qa FilterAttr, exp bool) {
	assert.Equal(t, exp, qa.Match(attrs))
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

func TestFilterAttrJson(t *testing.T) {
	f := func(fa FilterAttr, exp string) {
		// first test marshaling
		b, err := fa.MarshalText()
		assert.NoError(t, err)
		assert.Equal(t, exp, string(b))

		// now try unmarshaling that same string
		fa2, err := UnmarshalFilterAttr([]byte(exp))
		assert.NoError(t, err)

		// strings should be equal
		assert.Equal(t, fa.String(), fa2.String())

		// round trip should be equal
		b, err = fa.MarshalText()
		assert.NoError(t, err)
		assert.Equal(t, exp, string(b))
	}

	// basic exprs
	f(True(), `{"op":"true"}`)
	f(Exists("color"), `{"op":"exists","attr":"color"}`)
	return // TODO finish this
	f(Equal("color", "blue"), `{"op":"equal","attr":"color","val":"blue"}`)
	f(Regex("animal", "mo{3,5}se"), `{"op":"regex","attr":"animal","val":"mo{3,5}se"}`)

	// more complex exprs
	e1 := Equal("color", "blue")
	j1 := `{"op":"equal","attr":"color","val":"blue"}`
	f(e1, j1)

	e2 := Regex("animal", "mo{3,5}se")
	j2 := `{"op":"regex","attr":"animal","val":"mo{3,5}se"}`
	f(e2, j2)

	f(Not(e1), `{"op":"not","exprs":[`+j1+`]}`)
	f(Not(e2), `{"op":"not","exprs":[`+j2+`]}`)
	f(Or(e2), `{"op":"or","exprs":[`+j2+`]}`)
	f(Or(e2, e1), `{"op":"or","exprs":[`+j2+","+j1+`]}`)
	f(And(e2, e1), `{"op":"and","exprs":[`+j2+","+j1+`]}`)

	// even more complex
	e3 := Or(e2, e1)
	j3 := `{"op":"or","exprs":[` + j2 + "," + j1 + `]}`
	f(e3, j3)
	e4 := Not(e2)
	j4 := `{"op":"not","exprs":[` + j2 + `]}`
	f(e4, j4)
	f(And(e3, e4), `{"op":"and","exprs":[`+j3+","+j4+`]}`)
	f(Or(e3, e4, Not(True())), `{"op":"or","exprs":[`+j3+","+j4+`,{"op":"not","exprs":[{"op":"true"}]}]}`)
}
