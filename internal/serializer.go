package equinox

import (
	"bytes"
	"encoding/binary"
	"time"
)

type Serializer struct {
	valkey  *AttrMap
	attrkey *AttrMap
	attrval *AttrMap
}

var byteord = binary.BigEndian

func NewSerializer() *Serializer {
	s := Serializer{}
	s.valkey = NewAttrMap()
	s.attrkey = NewAttrMap()
	s.attrval = NewAttrMap()
	return &s
}

func (s *Serializer) Deserialize(b []byte) (*Point, error) {
	buf := bytes.NewReader(b)

	var umicro int64
	err := binary.Read(buf, byteord, &umicro)
	if err != nil {
		return nil, err
	}

	p := NewPoint(time.UnixMicro(umicro).UTC())

	return p, nil
}

/*
Serialization format:
timestamp: 8 bytes (64-bit)
values map length: 4 bytes (32-bit)
Then for each entry:
- key: 4 bytes (32-bit)
- value: 8 bytes (64-bit)

attributes map length: 4 bytes (32-bit)
Then for each entry:
- key: 4 bytes (32-bit)
- value: 4 bytes (32-bit)

Expected size (bytes) =  16 + 12*num_values + 8*num_attrs
*/
func (s *Serializer) Serialize(p *Point) ([]byte, error) {
	var buf bytes.Buffer

	// timestamp => 64-bit = 8 bytes
	err := binary.Write(&buf, byteord, p.ts.UnixMicro())
	if err != nil {
		return nil, err
	}

	// values: key -> value pairs
	err = binary.Write(&buf, byteord, uint32(len(p.vals)))
	if err != nil {
		return nil, err
	}
	for key, val := range p.vals {
		// write key
		err = binary.Write(&buf, byteord, s.valkey.ToIndex(key))
		if err != nil {
			return nil, err
		}

		// write value
		err = binary.Write(&buf, byteord, val)
		if err != nil {
			return nil, err
		}
	}

	// attributes: key -> value pairs
	err = binary.Write(&buf, byteord, uint32(len(p.attrs)))
	if err != nil {
		return nil, err
	}
	for key, val := range p.attrs {
		// write key
		err = binary.Write(&buf, byteord, s.attrkey.ToIndex(key))
		if err != nil {
			return nil, err
		}

		// write value
		err = binary.Write(&buf, byteord, s.attrval.ToIndex(val))
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}
