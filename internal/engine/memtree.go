package engine

import (
	"equinox/internal/core"
	"equinox/internal/query"
	"time"

	"github.com/google/btree"
)

type MemTree struct {
	buf *btree.BTreeG[*core.Point]
}

func NewMemTree() *MemTree {
	fn := func(a, b *core.Point) bool {
		r := core.PointCmp(a, b)
		if r == 0 {
			// incorporate the ID for strict ordering
			r = a.Id.Cmp(b.Id)
		}

		return r < 0
	}
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

func (mt *MemTree) Add(ps []*core.Point) error {
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
	mt   *MemTree     // reference to MemTree object
	st   *core.Point  // point where we start the search
	end  *core.Point  // point where we end the search
	last *core.Point  // last point returned
	q    *query.Query // query params
}

func (mtc *MemTreeCursor) Fetch(n int) ([]*core.Point, error) {
	if mtc.st == nil || mtc.end == nil || core.PointCmp(mtc.end, mtc.st) < 0 {
		// nothing to do if empty time range
		return nil, nil
	}

	// prealloc buffer for points
	r := make([]*core.Point, 0, n)

	// func that gets called on each iteration
	iter := func(p *core.Point) bool {
		// update starting point to current point
		mtc.st = p.Clone()

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
			if mtc.last == nil || !mtc.last.Identical(p) {
				r = append(r, p)
				mtc.last = p // remember last point added
			}
		}

		return true
	}

	mtc.mt.buf.AscendRange(mtc.st, mtc.end, iter)
	return r, nil
}

func (mt *MemTree) Search(q *query.Query) (*query.QueryExec, error) {
	// starting point is what was specified in the query
	st := core.NewPointEmptyId(q.Start)

	// ending point needs to be one microsecond past the query since
	// AscendRange uses < not <=
	end := core.NewPointEmptyId(time.UnixMicro(q.End.UnixMicro() + 1))

	mlc := &MemTreeCursor{mt: mt, q: q, st: st, end: end}
	return query.NewQueryExec(q, mlc), nil

}
