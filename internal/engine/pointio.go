package equinox

type PointIO interface {
	Add(p []*Point) error
	Len() int
	Vacuum() error
	Search(q *Query) (*QueryExec, error)
	Name() string
	String() string
}
