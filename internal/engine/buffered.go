package engine

import (
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
	return nil
}

func (b *Buffered) Vacuum() error {
	return nil
}

type BufferedCursor struct {
	buf *Buffered    // reference to MemTree object
	q   *query.Query // query params
}

func (bc *BufferedCursor) Fetch(n int) ([]*core.Point, error) {
	return nil, nil
}

func (b *Buffered) Search(q *query.Query) (*query.QueryExec, error) {
	bc := &BufferedCursor{buf: b, q: q}
	return query.NewQueryExec(q, bc), nil

}
