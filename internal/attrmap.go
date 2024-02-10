package equinox

type AttrMap struct {
	str2int map[string]uint32
	int2str map[uint32]string
}

func NewAttrMap() *AttrMap {
	a := AttrMap{}
	a.str2int = make(map[string]uint32)
	a.int2str = make(map[uint32]string)
	return &a
}

func (m *AttrMap) GetIndex(s string) (uint32, error) {
}
