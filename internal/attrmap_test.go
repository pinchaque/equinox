package equinox

import (
	"testing"
)

func TestBasic(t *testing.T) {
	s1 := "foo"
	i1 := uint32(0)
	//s2 := "bar"
	//s3 := "3.14"

	m := NewAttrMap()

	if m.Length() != 0 {
		t.Errorf("Length: got %d, wanted %d", m.Length(), 0)
	}

	if m.HasIndex(i1) {
		t.Errorf("Exist: index %d exists when it shouldn't", i1)
	}

	if m.HasAttr(s1) {
		t.Errorf("Exist: attr %s exists when it shouldn't", s1)
	}

	i1 = m.ToIndex(s1)

	if i1 != 0 {
		t.Errorf("Add: got index %d for attr %s expected %d", i1, s1, 0)
	}

	if !m.HasIndex(i1) {
		t.Errorf("Exist: index %d for attr %s is missing", i1, s1)
	}

	if !m.HasAttr(s1) {
		t.Errorf("Exist: attr %s is missing", s1)
	}

	if m.Length() != 1 {
		t.Errorf("Length: got %d, wanted %d", m.Length(), 1)
	}

}
