package engine

import (
	"container/list"
	"equinox/internal/core"
	"equinox/internal/query"
	"fmt"
	"slices"
	"strings"
	"time"
)

// Maintains list of Points ordered by timestamp
type MemList struct {
	buf *list.List
}

func NewMemList() *MemList {
	ml := MemList{}
	ml.buf = list.New()
	return &ml
}

func (ml *MemList) Name() string {
	return "MemList"
}

func (ml *MemList) String() string {
	var pstr []string
	i := 0
	for e := ml.buf.Front(); e != nil; e = e.Next() {
		p := e.Value.(*core.Point)
		pstr = append(pstr, fmt.Sprintf("%d: %s", i, p.String()))
		i++
	}
	return fmt.Sprintf("%s: {\n%s\n}", ml.Name(), strings.Join(pstr, "\n"))
}

func (ml *MemList) Add(ps ...*core.Point) error {
	/*
		we add the slice of points in two steps:
		(1) sort the slice
		(2) add the points to the end of the list, which is in ascending order
		by time stamp
		this methodology is designed to work best when we're adding close to
		the end of the list, which is typical since the normal use case is
		to add points at the current time
	*/

	if len(ps) == 0 {
		// nothing to do
		return nil
	}

	// sort the points we're adding
	slices.SortFunc(ps, core.PointCmp)

	// start inserting at the back of the list
	e := ml.buf.Back()

	// iterate over all points to insert them
	for i := len(ps) - 1; i >= 0; {
		p := ps[i]

		if e == nil { // at the front of the list
			ml.buf.PushFront(p) // add to the front
			i--
			// keep e as nil because we just want to keep adding to the front
			// we are guaranteed the next point will be before the last one
		} else if core.PointCmp(e.Value.(*core.Point), p) <= 0 {
			// new point should come after existing point
			ml.buf.InsertAfter(p, e)
			i--
			// don't change e because we know the next ps[i] will come before
			// the one we just inserted
		} else {
			// new point should come before existing point, so we move e backwards
			// and run the loop again with the same i
			e = e.Prev()
		}
	}
	return nil
}

func (ml *MemList) First() *core.Point {
	if ml.Len() == 0 {
		return nil
	}
	return ml.buf.Front().Value.(*core.Point)
}

func (ml *MemList) Last() *core.Point {
	if ml.Len() == 0 {
		return nil
	}
	return ml.buf.Back().Value.(*core.Point)
}

func (ml *MemList) Len() int {
	return ml.buf.Len()
}

func (ml *MemList) validate() error {
	if ml.buf.Len() <= 1 {
		return nil
	}

	for e := ml.buf.Front(); e.Next() != nil; e = e.Next() {
		p1 := e.Value.(*core.Point)
		p2 := e.Next().Value.(*core.Point)
		if core.PointCmp(p1, p2) > 0 {
			return fmt.Errorf("point (%s) incorrectly ordered before point (%s)", p1.String(), p2.String())
		}
	}
	return nil
}

func (ml *MemList) Flush() error {
	return nil
}

func (ml *MemList) Vacuum() error {
	return nil
}

func (ml *MemList) Move(start time.Time, end time.Time, dest PointIO) (int, error) {
	n := 0

	e := ml.buf.Front()
	for {
		if e == nil {
			break // end of list
		}

		p := e.Value.(*core.Point)

		// we're past the specified range so we can stop
		if p.Ts.UnixMicro() > end.UnixMicro() {
			break
		}

		// we're before the specified range, so loop
		if p.Ts.UnixMicro() < start.UnixMicro() {
			e = e.Next()
			continue
		}

		// move this point
		err := dest.Add(p)
		if err != nil {
			return n, err
		}
		edel := e
		e = e.Next()
		ml.buf.Remove(edel)
		n++
	}

	return n, nil
}

type MemListCursor struct {
	e *list.Element // element where we start the search
	q *query.Query  // query params
}

func (mlc *MemListCursor) Fetch(n int) ([]*core.Point, error) {
	if mlc.e == nil {
		// either the list was empty or we've already fetched everything
		return nil, nil
	}

	// prealloc buffer for points
	r := make([]*core.Point, 0, n)

	// iterate until we've filled the buffer or we're at the end of the list
	for ; len(r) < n && mlc.e != nil; mlc.e = mlc.e.Next() {
		p := mlc.e.Value.(*core.Point)

		// add matching points
		if mlc.q.Match(p) {
			r = append(r, p)
		}

		// since the list is ordered by time, we know there can't be more
		// results if the current list elem is after the query end time
		if p.Ts.UnixMicro() > mlc.q.End.UnixMicro() {
			mlc.e = nil // nothing more to look at
			break
		}
	}

	return r, nil
}

func (ml *MemList) Search(q *query.Query) (*query.QueryExec, error) {
	mlc := &MemListCursor{q: q, e: ml.buf.Front()}
	return query.NewQueryExec(q, mlc), nil
}
