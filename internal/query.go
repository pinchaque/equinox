package equinox

import (
	"fmt"
	"time"
)

// Represents the parameters for a query of points from the database. All queries
// must specify a time range as [start, end] and these are inclusive values.
// Queries may additionally specify attributes to filter on.
type Query struct {
	start time.Time
	end   time.Time
	qa    QueryAttr
}

func NewQuery(start time.Time, end time.Time, qa QueryAttr) *Query {
	q := Query{start: start, end: end, qa: qa}

	// ensure times are in correct order
	if start.UnixMicro() > end.UnixMicro() {
		q.start, q.end = q.end, q.start
	}

	return &q
}

// Returns string representation of the query
func (q *Query) String() string {
	return fmt.Sprintf("[%s-%s] [%s]", q.start.UTC(), q.end.UTC(), q.qa.String())
}

// Checks whether the given point is within the time range specified by the
// query. Returns -1 if the point is before the range, 0 if within, 1 after.
// Does not check the point against the attributes for the query.
func (q *Query) CmpTime(p *Point) int {
	if p.ts.UnixMicro() < q.start.UnixMicro() {
		return -1
	} else if p.ts.UnixMicro() > q.end.UnixMicro() {
		return 1
	} else {
		return 0
	}
}

// Returns true if the given point matches the attributes for this query,
// false otherwise. If the query has no attributes then all points will match.
// Does not check the point against the time range.
func (q *Query) MatchAttr(p *Point) bool {
	return q.qa.Match(p.attrs)
}

// Returns true if the point matches both the time range and attributes specified
// by this query, false otherwise.
func (q *Query) Match(p *Point) bool {
	return q.CmpTime(p) == 0 && q.MatchAttr(p)
}
