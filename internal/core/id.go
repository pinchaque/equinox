package core

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math/rand"
)

// Representation of an ID used for data points
type Id struct {
	val uint64
}

// Creates a new random ID
func NewId() *Id {
	id := Id{val: rand.Uint64()}
	return &id
}

// Creates an Id struct from the specified string, which must be a Base64 Url
// Encoding representing a uint64.
func IdFromString(s string) (*Id, error) {
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("error decoding string '%s': %s", s, err.Error())
	}

	// check byte length
	if len(b) != 8 {
		return nil, fmt.Errorf("invalid num bytes %d when decoding string '%s'", len(b), s)
	}

	// read into a uint64
	var i uint64
	buf := bytes.NewReader(b)
	err = binary.Read(buf, binary.BigEndian, &i)
	if err != nil {
		return nil, fmt.Errorf("error reading uint64 from string '%s': %s", s, err.Error())
	}

	// create Id struct and return
	id := Id{val: i}
	return &id, nil
}

// Returns string representation of an ID, which is a Base64 Url Encoding
func (id *Id) String() string {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, id.val)
	if err != nil {
		panic(fmt.Sprintf("failed to write id %d to string", id.val))
	}

	return base64.URLEncoding.EncodeToString(buf.Bytes())
}

// Compares two Id structs, return -1 if this one is less than other, 1 if this
// is greater than other, 0 if equal. This can be used to check for uniqueness
// and duplicate IDs. The current implementation just compares the underlying
// uint64 values.
func (id *Id) Cmp(oth Id) int {
	if id.val < oth.val {
		return -1
	} else if id.val > oth.val {
		return 1
	} else {
		return 0
	}
}
