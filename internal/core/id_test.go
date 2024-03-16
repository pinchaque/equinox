package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdString(t *testing.T) {
	id := Id{val: 2822340188419286878}
	exp := "Jyr3cq4KZ14="
	assert.Equal(t, exp, id.String())
}
func TestIdClone(t *testing.T) {
	id := Id{val: 2822340188419286878}
	id2 := id.Clone()
	assert.Equal(t, id.val, id2.val)
	assert.Equal(t, id.String(), id2.String())
}

func TestIdMarshal(t *testing.T) {
	id := Id{val: 2822340188419286878}
	exp := "Jyr3cq4KZ14="
	act, err := id.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, exp, string(act))

	// test unmarshaling
	id2 := Id{val: 0}
	err = id2.UnmarshalText(act)
	assert.NoError(t, err)
	assert.Equal(t, id.val, id2.val)
	assert.Equal(t, id.String(), id2.String())
}

func TestIdRoundtrip(t *testing.T) {
	// run this a bunch of times
	for i := 0; i < 500; i++ {

		id1 := NewId()
		s := id1.String()
		id2, err := IdFromString(s)
		if err != nil {
			t.Errorf("failed to create Id from string %s: %s", s, err.Error())
		}

		if id1.val != id2.val {
			t.Errorf("Id roundtrip values not equal, expected %d got %d", id1.val, id2.val)
		}

		if id1.String() != id2.String() {
			t.Errorf("Id roundtrip strings not equal, expected %s got %s", id1.String(), id2.String())
		}
	}
}

func TestIdErrors(t *testing.T) {
	var id *Id
	var err error
	var s, msg string

	// missing the trailing "=" that is padding
	s = "Jyr3cq4KZ14"
	msg = "error decoding string 'Jyr3cq4KZ14': illegal base64 data at input byte 8"
	id, err = IdFromString(s)
	if err == nil {
		t.Errorf("expected error when parsing %s, got success (%d, %s)", s, id.val, id.String())
	} else if err.Error() != msg {
		t.Errorf("unexpected error message when parsing %s: expected '%s' got '%s'",
			s, msg, err.Error())
	}

	// garbage
	s = "Jyr$*#()"
	msg = "error decoding string 'Jyr$*#()': illegal base64 data at input byte 3"
	id, err = IdFromString(s)
	if err == nil {
		t.Errorf("expected error when parsing %s, got success (%d, %s)", s, id.val, id.String())
	} else if err.Error() != msg {
		t.Errorf("unexpected error message when parsing %s: expected '%s' got '%s'",
			s, msg, err.Error())
	}

	// empty string
	s = ""
	msg = "invalid num bytes 0 when decoding string ''"
	id, err = IdFromString(s)
	if err == nil {
		t.Errorf("expected error when parsing %s, got success (%d, %s)", s, id.val, id.String())
	} else if err.Error() != msg {
		t.Errorf("unexpected error message when parsing %s: expected '%s' got '%s'",
			s, msg, err.Error())
	}

	// partial string
	s = "Jyr3cq4K"
	msg = "invalid num bytes 6 when decoding string 'Jyr3cq4K'"
	id, err = IdFromString(s)
	if err == nil {
		t.Errorf("expected error when parsing %s, got success (%d, %s)", s, id.val, id.String())
	} else if err.Error() != msg {
		t.Errorf("unexpected error message when parsing %s: expected '%s' got '%s'",
			s, msg, err.Error())
	}
}

func TestIdCmp(t *testing.T) {
	id1 := &Id{val: 12345678}
	id2 := &Id{val: 45678912}

	fn := func(i1 *Id, i2 *Id, exp int) {
		act := i1.Cmp(i2)
		if act != exp {
			t.Errorf("error comparing [%d]%s to [%d]%s: expected %d got %d",
				i1.val, i1.String(), i2.val, i2.String(), exp, act)
		}
	}

	fn(id1, id2, -1)
	fn(id2, id1, 1)
	fn(id1, id1, 0)
	fn(id2, id2, 0)
}
