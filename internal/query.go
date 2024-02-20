package equinox

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Cursor interface {
	// Fetches the next n results from the cursor. Returns a nil slice if there
	// are no more to return.
	Fetch(n int) ([]*Point, error)
}

type Query struct {
	start time.Time
	end   time.Time
	attrs map[string]string
}

func NewQuery(start time.Time, end time.Time) *Query {
	q := Query{start: start, end: end}
	q.attrs = make(map[string]string)
	return &q
}

func (q *Query) String() string {
	var attr []string

	for k, v := range q.attrs {
		attr = append(attr, k+": "+v)
	}
	sort.Strings(attr) // ensure consistent output

	return fmt.Sprintf("[%s-%s] [%s]",
		q.start.UTC(),
		q.end.UTC(),
		strings.Join(attr, ", "))

}
