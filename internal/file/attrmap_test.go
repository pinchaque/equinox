package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	s1 := "foo"
	i1 := uint32(0)
	//s2 := "bar"
	//s3 := "3.14"

	m := NewAttrMap()
	assert.Equal(t, uint32(0), m.Length())
	assert.False(t, m.HasIndex(i1))
	assert.False(t, m.HasAttr(s1))

	i1 = m.ToIndex(s1)

	assert.Equal(t, uint32(0), i1)
	assert.True(t, m.HasIndex(i1))
	assert.True(t, m.HasAttr(s1))
	assert.Equal(t, uint32(1), m.Length())
}

func TestMultiple(t *testing.T) {
	s1 := "foo"
	s2 := "bar"
	s3 := "3.14"

	m := NewAttrMap()
	i1 := m.ToIndex(s1)
	i2 := m.ToIndex(s2)
	i3 := m.ToIndex(s3)

	assert.Equal(t, uint32(3), m.Length())

	// make sure indexes work
	str1, exist := m.AtIndex(i1)
	assert.True(t, exist)
	assert.Equal(t, s1, str1)

	str2, exist := m.AtIndex(i2)
	assert.True(t, exist)
	assert.Equal(t, s2, str2)

	str3, exist := m.AtIndex(i3)
	assert.True(t, exist)
	assert.Equal(t, s3, str3)
}

func TestDelete(t *testing.T) {
	s1 := "foo"
	s2 := "bar"
	s3 := "3.14"

	m := NewAttrMap()
	i1 := m.ToIndex(s1)
	i2 := m.ToIndex(s2)

	assert.Equal(t, uint32(2), m.Length())

	// make sure indexes work
	str1, exist := m.AtIndex(i1)
	assert.True(t, exist)
	assert.Equal(t, s1, str1)

	str2, exist := m.AtIndex(i2)
	assert.True(t, exist)
	assert.Equal(t, s2, str2)

	// test deletion
	m.DeleteAttr(s2)

	assert.Equal(t, uint32(1), m.Length())

	m.DeleteAttr(s2) // should be a no-op
	assert.Equal(t, uint32(1), m.Length())

	_, exist = m.AtIndex(i2)
	assert.False(t, exist)

	// test index when adding another
	i3 := m.ToIndex(s3)
	assert.Equal(t, uint32(2), i3)
}
