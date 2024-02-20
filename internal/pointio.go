package equinox

type PointIO interface {
	Add(p []*Point) error
	Vacuum() error
	Search(q *Query) (Cursor, error)
	String() string
}
