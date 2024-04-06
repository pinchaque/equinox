package engine

import (
	"container/list"
	"equinox/internal/core"
	"equinox/internal/query"
	"fmt"
	"strings"
	"time"
)

// Stores data points in a buffer until time has passed or the buffer is full
// and then flushes them to archival storage. This is designed to be used to
// create efficient persistent storage by combining an in-memory buffer with
// a datafile-backed storage option.
type Buffered struct {
	// The PointIO interface to use as the buffer
	Buf PointIO

	// The PointIO interface to use as the archival storage
	Archive PointIO

	// The minimum time duration items will stay in the buffer. Items that are
	// older than this will be flushed to Archive the next time Flush() is
	// called.
	Dur time.Duration
}

func NewBuffered(buf PointIO, archive PointIO) *Buffered {
	d, err := time.ParseDuration("15m")
	if err != nil {
		panic("Failed to initailize duration")
	}
	b := Buffered{Buf: buf, Archive: archive, Dur: d}
	return &b
}

func (b *Buffered) Name() string {
	return "Buffered"
}

func (b *Buffered) String() string {
	var pstr []string

	pstr = append(pstr, fmt.Sprintf("Buffered %s to %s with %s delay", b.Buf.Name(), b.Archive.Name(), b.Dur))
	pstr = append(pstr, b.Buf.String())
	pstr = append(pstr, b.Archive.String())
	return strings.Join(pstr, "\n")
}

func (b *Buffered) Add(ps ...*core.Point) error {
	for _, p := range ps {
		// if the time is on or after what's in archive, then we can simply
		// add it to the buffer
		last := b.Archive.Last()
		if last == nil || p.Ts.UnixMicro() >= last.Ts.UnixMicro() {
			err := b.Buf.Add(p)
			if err != nil {
				return err
			}
		} else {
			err := b.Archive.Add(p)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *Buffered) First() *core.Point {
	p := b.Archive.First()
	if p != nil {
		return p
	}
	return b.Buf.First()
}

func (b *Buffered) Last() *core.Point {
	p := b.Buf.Last()
	if p != nil {
		return p
	}
	return b.Archive.Last()
}

func (b *Buffered) Len() int {
	return b.Buf.Len() + b.Archive.Len()
}

func (b *Buffered) Flush() error {
	// TODO: need to implement this
	/*
		General plan:
		(a) Iterate over all the elements of b.Buf that are older than b.Dur
		(b) Add them to b.Archive
		(c) Remove them from b.Buf

		We need primitives on PointIO to make this possible.
	*/
	return nil
}

func (b *Buffered) Vacuum() error {
	err := b.Buf.Vacuum()
	if err != nil {
		return err
	}

	return b.Archive.Vacuum()
}

func (b *Buffered) Move(start time.Time, end time.Time, dest PointIO) (int, error) {
	return 0, fmt.Errorf("not implemented")
}

type BufferedCursor struct {
	buf    *Buffered        // reference to Buffered object
	qebuf  *query.QueryExec // executing query on buf.Buf
	qearch *query.QueryExec // executing query on buf.Archive
	ptbuf  *list.List       // point buffer that we fill for executing the query
	q      *query.Query     // query params

}

// Fills ptbuf with points from the specified QueryExec. Ensures there are at
// least n points in ptbuf by the time it returns as long as there are that
// many results to return from qe. In other words, ptbuf will only be <n if we
// are at the end of the query results.
// This function and ptbuf exists because we are always reading in batches from
// the underlying data engines. If we read 100 points but only need 10 then we
// use this buf to store the 90 we didn't use in that fetch.
func (bc *BufferedCursor) fillPointBuf(qe *query.QueryExec, n int) error {
	// nothing to do if we already have enough
	if bc.ptbuf.Len() >= n {
		return nil
	}

	// fetch from the QueryExec in batches of n
	for {
		if qe.Done() {
			return nil
		}

		pts, err := qe.Fetch(n)
		if err != nil {
			return err
		}

		for _, p := range pts {
			// we assume the underlying engines are verifying the match so we
			// accept the points as-is
			bc.ptbuf.PushBack(p)
		}

		// we have enough - return now
		if bc.ptbuf.Len() >= n {
			return nil
		}
	}
}

func (bc *BufferedCursor) Fetch(n int) ([]*core.Point, error) {
	// prealloc buffer for points
	r := make([]*core.Point, 0, n)

	// nothing to do
	if n <= 0 {
		return r, nil
	}

	// Fill up the point buffer from Archive and then Buf. Note that this
	// will exhaust Archive before using Buf.
	err := bc.fillPointBuf(bc.qearch, n)
	if err != nil {
		return r, err
	}

	err = bc.fillPointBuf(bc.qebuf, n)
	if err != nil {
		return r, err
	}

	// Now we fill r from ptbuf by popping elements off the front.
	// Note that we might not have enough points in ptbuf if we are at the end
	// of the result list. Or we might have more than enough.
	for {
		e := bc.ptbuf.Front()
		if e == nil {
			break // empty list
		}

		// remove the first point
		p := bc.ptbuf.Remove(e).(*core.Point)

		// add it to the results we'll return
		r = append(r, p)

		// if we have enough results then we can return them
		if len(r) >= n {
			break
		}
	}

	// return what we got - may be less than n
	return r, nil
}

func (b *Buffered) Search(q *query.Query) (*query.QueryExec, error) {
	// initialize the queries across the archive and buffer
	qearch, err := b.Archive.Search(q)
	if err != nil {
		return nil, err
	}

	qebuf, err := b.Buf.Search(q)
	if err != nil {
		return nil, err
	}

	bc := &BufferedCursor{buf: b, q: q, qebuf: qebuf, qearch: qearch, ptbuf: list.New()}
	return query.NewQueryExec(q, bc), nil

}
