package equinox

import (
	"container/list"
)

type MemList struct {
	buf *list.List
}

func NewMemList() *MemList {
	ml := MemList{}
	ml.buf = list.New()
	return &ml
}

func (ml *MemList) String() string {
	return "MemList"
}

func (ml *MemList) Add(p []*Point) error {
	return nil
}

func (ml *MemList) Vacuum() error {
	return nil
}

type MemListCursor struct {
}

func NewMemListCursor() *MemListCursor {
	mlc := MemListCursor{}
	return &mlc
}

func (mlc *MemListCursor) Fetch(n int) ([]*Point, error) {
	return nil, nil
}

func (ml *MemList) Search(q *Query) (Cursor, error) {
	return NewMemListCursor(), nil
}
