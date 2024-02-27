package equinox

import (
	"time"

	"github.com/google/btree"
)

type MemTree struct {
	buf *btree.BTreeG[*Point]
}

func NewMemTree() *MemTree {
	fn := func(a, b *Point) bool { return a.Less(b) }
	mt := MemTree{}
	mt.buf = btree.NewG(2, fn)
	return &mt
}

func (mt *MemTree) Name() string {
	return "MemTree"
}

func (mt *MemTree) String() string {
	return "MemTree"
}

func (mt *MemTree) Add(ps []*Point) error {
	for _, p := range ps {
		mt.buf.ReplaceOrInsert(p)
	}
	return nil
}

func (mt *MemTree) Len() int {
	return mt.buf.Len()
}

func (mt *MemTree) Vacuum() error {
	return nil
}

type MemTreeCursor struct {
	mt   *MemTree // reference to MemTree object
	st   *Point   // point where we start the search
	end  *Point   // point where we end the search
	last *Point   // last point returned
	q    *Query   // query params
}

func (mtc *MemTreeCursor) fetch(n int) ([]*Point, error) {
	if mtc.st == nil || mtc.end == nil || mtc.end.Less(mtc.st) {
		// nothing to do if empty time range
		return nil, nil
	}

	// prealloc buffer for points
	r := make([]*Point, 0, n)

	// func that gets called on each iteration
	iter := func(p *Point) bool {
		// update starting point to current point
		mtc.st = NewPoint(p.ts)

		// if we're already full then we need to stop now and we'll try this
		// point again on the next call to fetch
		if len(r) >= n {
			return false
		}

		// add point if it matches
		if mtc.q.Match(p) {
			// don't add it if it matches the previous returned point
			// this fixes an edge case where AscendRange stops naturally on
			// an added point
			if mtc.last == nil || !mtc.last.Equal(p) {
				r = append(r, p)
				mtc.last = p // remember last point added
			}
		}

		// TODO: we need strict ordering and avoid duplicates. GUID would help!
		return true
	}

	mtc.mt.buf.AscendRange(mtc.st, mtc.end, iter)
	return r, nil
}

func (mt *MemTree) Search(q *Query) (*QueryExec, error) {
	// starting point is what was specified in the query
	st := NewPoint(q.start)

	// ending point needs to be one microsecond past the query since
	// AscendRange uses < not <=
	end := NewPoint(time.UnixMicro(q.end.UnixMicro() + 1))

	mlc := &MemTreeCursor{mt: mt, q: q, st: st, end: end}
	return NewQueryExec(q, mlc), nil

}
