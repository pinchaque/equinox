package equinox

import (
	"os"
)

type DataFile struct {
	path        string
	num_records int
	header_size int
	record_size int
	fd          os.File
}

func OpenFile(path string) *DataFile {
	df := DataFile{}
	df.path = path
	df.num_records = 0
	return &df
}

func (df *DataFile) CloseFile() error {
	return df.fd.Close()
}

func (df *DataFile) Write(idx int, p Point) error {
	return nil
}

func (df *DataFile) Read(idx int) (*Point, error) {
	return &Point{}, nil
}
