package file

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

func TestMultiple(t *testing.T) {
	s1 := "foo"
	s2 := "bar"
	s3 := "3.14"

	m := NewAttrMap()
	i1 := m.ToIndex(s1)
	i2 := m.ToIndex(s2)
	i3 := m.ToIndex(s3)

	if m.Length() != 3 {
		t.Errorf("Length: got %d, wanted %d", m.Length(), 3)
	}

	// make sure indexes work
	str1, exist := m.AtIndex(i1)
	if !exist {
		t.Errorf("AtIndex: index %d for attr %s is missing", i1, s1)
	}

	if s1 != str1 {
		t.Errorf("AtIndex: index %d got attr %s, expected %s", i1, str1, s1)
	}

	str2, exist := m.AtIndex(i2)
	if !exist {
		t.Errorf("AtIndex: index %d for attr %s is missing", i2, s2)
	}

	if s2 != str2 {
		t.Errorf("AtIndex: index %d got attr %s, expected %s", i2, str2, s2)
	}

	str3, exist := m.AtIndex(i3)
	if !exist {
		t.Errorf("AtIndex: index %d for attr %s is missing", i3, s3)
	}

	if s3 != str3 {
		t.Errorf("AtIndex: index %d got attr %s, expected %s", i3, str3, s3)
	}
}

func TestDelete(t *testing.T) {
	s1 := "foo"
	s2 := "bar"
	s3 := "3.14"

	m := NewAttrMap()
	i1 := m.ToIndex(s1)
	i2 := m.ToIndex(s2)

	if m.Length() != 2 {
		t.Errorf("Length: got %d, wanted %d", m.Length(), 2)
	}

	// make sure indexes work
	str1, exist := m.AtIndex(i1)
	if !exist {
		t.Errorf("AtIndex: index %d for attr %s is missing", i1, s1)
	}

	if s1 != str1 {
		t.Errorf("AtIndex: index %d got attr %s, expected %s", i1, str1, s1)
	}

	str2, exist := m.AtIndex(i2)
	if !exist {
		t.Errorf("AtIndex: index %d for attr %s is missing", i2, s2)
	}

	if s2 != str2 {
		t.Errorf("AtIndex: index %d got attr %s, expected %s", i2, str2, s2)
	}

	// test deletion
	m.DeleteAttr(s2)

	if m.Length() != 1 {
		t.Errorf("Length: got %d, wanted %d", m.Length(), 1)
	}

	m.DeleteAttr(s2) // should be a no-op
	if m.Length() != 1 {
		t.Errorf("Length: got %d, wanted %d", m.Length(), 1)
	}

	_, exist = m.AtIndex(i2)
	if exist {
		t.Errorf("AtIndex: index %d for attr %s still exists", i2, s2)
	}

	// test index when adding another
	i3 := m.ToIndex(s3)
	if i3 != 2 {
		t.Errorf("ToIndex: %s was given index %d, expected %d", s3, i3, 2)
	}
}
