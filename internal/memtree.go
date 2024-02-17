package equinox

import (
	"github.com/google/btree"
)

type MemTree struct {
	buf *btree.BTree
}

func NewMemTree() *MemTree {
	mt := MemTree{}
	mt.buf = btree.New(2)
	return &mt
}

func (mt *MemTree) Add(p []*Point) error {
	return nil
}

func (mt *MemTree) Vacuum() error {
	return nil
}

type MemTreeCursor struct {
}

func NewMemTreeCursor() *MemTreeCursor {
	mtc := MemTreeCursor{}
	return &mtc
}

func (mtc *MemTreeCursor) Fetch(n int) ([]*Point, error) {
	return nil, nil
}

func (mt *MemTree) Search(q *Query) (Cursor, error) {
	return NewMemTreeCursor(), nil
}
