package equinox

// Internal interface used by QueryExec to retrieve results from the diferent
// data stores.
type Cursor interface {
	// Fetches the next n results from the cursor. Returns a nil slice if there
	// are no more to return.
	Fetch(n int) ([]*Point, error)
}

type QueryExec struct {
	q    *Query
	cur  Cursor
	done bool // whether we've hit end of query already
}

func NewQueryExec(q *Query, cur Cursor) *QueryExec {
	qe := QueryExec{q: q, cur: cur, done: false}
	return &qe
}

// Fetches the next n results from the query. Returns empty slice (nil) if
// there are no more. Returns an error if we aren't done but there was an
// error in running the query.
func (qe *QueryExec) Fetch(n int) ([]*Point, error) {
	return nil, nil
}

// Returns true if we've returned all results from this query, false otherwise.
func (qe *QueryExec) Done() bool {
	return qe.done
}
