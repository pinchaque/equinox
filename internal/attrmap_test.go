package equinox

import (
	"testing"
)

func TestBasic(t *testing.T) {
	s1 := "foo"
	//s2 := "bar"
	//s3 := "3.14"

	m := NewAttrMap()

	if m.Length() != 0 {
		t.Errorf("Length: got %d, wanted %d", m.Length(), 0)
	}

	if m.HasAttr(s1) {
		t.Errorf("Exist: attr %s exists when it shouldn't", s1)
	}

	idx1, exists := m.AddAttr(s1)
	if exists {
		t.Errorf("Add: attribute %s already existed: %s", s1)
	}

	if idx1 != 0 {
		t.Errorf("Add: got index %d for attr %s expected %d", idx1, s1, 0)
	}

	if !m.HasIndex(idx1) {
		t.Errorf("Exist: index %d for attr %s is missing", idx1, s1)
	}

	if !m.HasAttr(s1) {
		t.Errorf("Exist: attr %s is missing", s1)
	}

	if m.Length() != 1 {
		t.Errorf("Length: got %d, wanted %d", m.Length(), 1)
	}

}
