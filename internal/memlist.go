package equinox

import (
	"container/list"
	"fmt"
	"slices"
	"strings"
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
		p := e.Value.(*Point)
		pstr = append(pstr, fmt.Sprintf("%d: %s", i, p.String()))
		i++
	}
	return fmt.Sprintf("%s: {\n%s\n}", ml.Name(), strings.Join(pstr, "\n"))
}

func (ml *MemList) Add(ps []*Point) error {
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
	slices.SortFunc(ps, PointCmp)

	// start inserting at the back of the list
	e := ml.buf.Back()

	// iterate over all points to insert them
	for i := len(ps) - 1; i >= 0; {
		p := ps[i]

		if e == nil {
			ml.buf.PushFront(p)
			i--
			// keep e as nil because we just want to keep adding to the front
			// we are guaranteed the next point will be before the last one
		} else if PointCmp(e.Value.(*Point), p) <= 0 {
			ml.buf.InsertAfter(p, e)
			i--
			// don't change e because we know the next ps[i] will come before
			// the one we just inserted
		} else {
			// point p needs to come before the current e, so we move it backwards
			// and run the loop again with teh same i
			e = e.Prev()
		}
	}
	return nil
}

func (ml *MemList) Len() int {
	return ml.buf.Len()
}

func (ml *MemList) validate() error {
	if ml.buf.Len() <= 1 {
		return nil
	}

	for e := ml.buf.Front(); e.Next() != nil; e = e.Next() {
		p1 := e.Value.(*Point)
		p2 := e.Next().Value.(*Point)
		if PointCmp(p1, p2) > 0 {
			return fmt.Errorf("point (%s) incorrectly ordered before point (%s)", p1.String(), p2.String())
		}
	}
	return nil
}

func (ml *MemList) Vacuum() error {
	return nil
}

type MemListCursor struct {
	q     *Query
	eprev *list.Element // last element returned
	eof   bool          // whether we've hit end of query already
}

func NewMemListCursor(q *Query) *MemListCursor {
	mlc := MemListCursor{}
	mlc.q = q
	mlc.eprev = nil
	mlc.eof = false
	return &mlc
}

func (mlc *MemListCursor) Fetch(n int) ([]*Point, error) {
	return nil, nil
}

func (ml *MemList) Search(q *Query) (Cursor, error) {
	return NewMemListCursor(q), nil
}
