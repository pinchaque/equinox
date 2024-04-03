package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterAttrJsonxxx(t *testing.T) {
	f := func(qa FilterAttr, exp string) {
		b, err := qa.MarshalText()
		assert.NoError(t, err)
		assert.Equal(t, exp, string(b))
	}

	f(Exists("color"), `{"op":"exists","attr":"color"}`)
}
