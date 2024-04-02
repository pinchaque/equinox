package query

import (
	"equinox/internal/core"
	"fmt"
	"time"
)

// Represents the parameters for a query of points from the database. All queries
// must specify a time range as [start, end] and these are inclusive values.
// Queries may additionally specify attributes to filter on.
type Query struct {
	Start time.Time
	End   time.Time
	QA    FilterAttr
}

func NewQuery(start time.Time, end time.Time, qa FilterAttr) *Query {
	q := Query{Start: start, End: end, QA: qa}

	// ensure times are in correct order
	if start.UnixMicro() > end.UnixMicro() {
		q.Start, q.End = q.End, q.Start
	}

	return &q
}

// Returns string representation of the query
func (q *Query) String() string {
	return fmt.Sprintf("[%s-%s] [%s]", q.Start.UTC(), q.End.UTC(), q.QA.String())
}

// Checks whether the given point is within the time range specified by the
// query. Returns -1 if the point is before the range, 0 if within, 1 after.
// Does not check the point against the attributes for the query.
func (q *Query) CmpTime(p *core.Point) int {
	if p.Ts.UnixMicro() < q.Start.UnixMicro() {
		return -1
	} else if p.Ts.UnixMicro() > q.End.UnixMicro() {
		return 1
	} else {
		return 0
	}
}

// Returns true if the given point matches the time specified by this query,
// false otherwise. Does not check the point against attributes.
func (q *Query) MatchTime(p *core.Point) bool {
	return q.CmpTime(p) == 0
}

// Returns true if the given point matches the attributes for this query,
// false otherwise. If the query has no attributes then all points will match.
// Does not check the point against the time range.
func (q *Query) MatchAttr(p *core.Point) bool {
	return q.QA.Match(p.Attrs)
}

// Returns true if the point matches both the time range and attributes specified
// by this query, false otherwise.
func (q *Query) Match(p *core.Point) bool {
	return q.MatchTime(p) && q.MatchAttr(p)
}
