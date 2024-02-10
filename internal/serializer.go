package equinox

import (
	"bytes"
	"encoding/binary"
	"time"
)

type Serializer struct {
	valkey  AttrMap
	attrkey AttrMap
	attrval AttrMap
}

func (s *Serializer) Deserialize(b []byte) (*Point, error) {
	return NewPoint(time.Now()), nil
}

func (s *Serializer) Serialize(p *Point) ([]byte, error) {
	var buf bytes.Buffer
	ord := binary.BigEndian

	// timestamp
	err := binary.Write(&buf, ord, p.ts.UnixMicro())
	if err != nil {
		return nil, err
	}

	// values: key -> value pairs
	err = binary.Write(&buf, ord, len(p.vals))
	if err != nil {
		return nil, err
	}
	for key, val := range p.vals {
		// write key
		err = binary.Write(&buf, ord, s.valkey.GetIndex(key))
		if err != nil {
			return nil, err
		}

		// write value
		err = binary.Write(&buf, ord, val)
		if err != nil {
			return nil, err
		}
	}

	// attributes: key -> value pairs
	err = binary.Write(&buf, ord, len(p.attrs))
	if err != nil {
		return nil, err
	}
	for key, val := range p.attrs {
		// write key
		err = binary.Write(&buf, ord, s.attrkey.GetIndex(key))
		if err != nil {
			return nil, err
		}

		// write value
		err = binary.Write(&buf, ord, s.attrval.GetIndex(val))
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}
