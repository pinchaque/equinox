package equinox

import (
	"fmt"
	"os"
)

type DataFile struct {
	path        string
	num_records int
	header_size int
	record_size int
	fd          *os.File
	ser         Serializer
}

func NewDataFile(path string, ser Serializer) *DataFile {
	df := DataFile{}
	df.path = path
	df.num_records = 0
	df.header_size = 0
	df.ser = ser
	return &df
}

func (df *DataFile) Open() error {
	var err error
	df.fd, err = os.Create(df.path)
	return err
}

func (df *DataFile) Close() error {
	return df.fd.Close()
}

func (df *DataFile) getOffset(idx int) int64 {
	return 0 // TODO implement this
}

func (df *DataFile) Write(idx int, p *Point) error {
	data, err := df.ser.Serialize(p)
	if err != nil {
		return err
	}

	_, err = df.fd.Seek(df.getOffset(idx), 0)
	if err != nil {
		return err
	}

	_, err = df.fd.Write(data)
	if err != nil {
		return err
	}

	err = df.fd.Sync()
	if err != nil {
		return err
	}

	return nil
}

func (df *DataFile) Read(idx int) (*Point, error) {
	_, err := df.fd.Seek(df.getOffset(idx), 0)
	if err != nil {
		return nil, err
	}

	data := make([]byte, df.record_size)
	var n int
	n, err = df.fd.Read(data)
	if err != nil {
		return nil, err
	}

	if n != df.record_size {
		return nil, fmt.Errorf("bytes read (%d) doesn't match record size (%d)", n, df.record_size)
	}

	var p *Point
	p, err = df.ser.Deserialize(data)

	return p, nil
}
