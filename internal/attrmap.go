package equinox

type AttrMap struct {
	str2int map[string]uint32
	int2str map[uint32]string
	numattr uint32
}

func NewAttrMap() *AttrMap {
	m := AttrMap{}
	m.str2int = make(map[string]uint32)
	m.int2str = make(map[uint32]string)
	m.numattr = 0
	return &m
}

func (m *AttrMap) Length() uint32 {
	return m.numattr
}

func (m *AttrMap) AtIndex(idx uint32) (string, bool) {
	v, exist := m.int2str[idx]
	return v, exist
}

func (m *AttrMap) HasIndex(idx uint32) bool {
	_, exist := m.int2str[idx]
	return exist
}

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
