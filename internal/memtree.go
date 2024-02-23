package equinox

import (
	"fmt"

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

func (mt *MemTree) Name() string {
	return "MemTree"
}

func (mt *MemTree) String() string {
	return "MemTree"
}

func (mt *MemTree) Add(p []*Point) error {
	return fmt.Errorf("not implemented")
}

func (mt *MemTree) Len() int {
	return 0
}

func (mt *MemTree) Vacuum() error {
	return nil
}

type MemTreeCursor struct {
}

func (mtc *MemTreeCursor) fetch(n int) ([]*Point, error) {
	return nil, fmt.Errorf("not implemented")
}

func (mt *MemTree) Search(q *Query) (*QueryExec, error) {
	return nil, fmt.Errorf("not implemented")
}
