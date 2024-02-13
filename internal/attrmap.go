package equinox

type AttrMap struct {
	str2int map[string]uint32
	int2str map[uint32]string
	numattr uint32
}

// Creates new empty AttrMap
func NewAttrMap() *AttrMap {
	m := AttrMap{}
	m.str2int = make(map[string]uint32)
	m.int2str = make(map[uint32]string)
	m.numattr = 0
	return &m
}

// Returns number of elements in AttrMap
func (m *AttrMap) Length() uint32 {
	return uint32(len(m.str2int))
}

// Returns the attribute at the specified index along with bool specifying
// whether it exists or not
func (m *AttrMap) AtIndex(idx uint32) (string, bool) {
	v, exist := m.int2str[idx]
	return v, exist
}

// Returns true if the AttrMap has the specified index; does not create it
// if it doesn't exist
func (m *AttrMap) HasIndex(idx uint32) bool {
	_, exist := m.int2str[idx]
	return exist
}

// Returns true if the AttrMap has the specified attribute; does not create it
// if it doesn't exist
func (m *AttrMap) HasAttr(s string) bool {
	_, exist := m.str2int[s]
	return exist
}

// Transforms given attribute to an index, creating it in the map if it doesn't
// already exist
func (m *AttrMap) ToIndex(s string) uint32 {
	idx, exist := m.str2int[s]
	if exist {
		return idx
	}

	idx = m.numattr
	m.numattr++
	m.int2str[idx] = s
	m.str2int[s] = idx
	return idx
}

// Deletes attribute
func (m *AttrMap) DeleteAttr(s string) {
	idx, exist := m.str2int[s]
	if !exist {
		return // nothing to do
	}
	delete(m.int2str, idx)
	delete(m.str2int, s)
}
