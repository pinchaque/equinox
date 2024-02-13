package equinox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"unsafe"
)

type DataFile struct {
	path        string
	num_records int
	header_size uint32
	record_size uint32
	fd          *os.File
	ser         *Serializer
}

func OpenExistingDF(path string, ser *Serializer) (*DataFile, error) {
	df, err := newDataFile(path, ser)
	if err != nil {
		return nil, err
	}

	// open existing file for R/W (do not create if doesn't exist)
	df.fd, err = os.OpenFile(df.path, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	err = df.parseHeader()
	if err != nil {
		return nil, err
	}

	return df, nil
}

func OpenNewDF(path string, ser *Serializer, recsize uint32) (*DataFile, error) {
	df, err := newDataFile(path, ser)
	if err != nil {
		return nil, err
	}

	// create file for read and write
	df.fd, err = os.OpenFile(df.path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return nil, err
	}

	// write the header
	df.record_size = recsize
	err = df.writeHeader()
	if err != nil {
		return nil, err
	}

	return df, nil
}

func newDataFile(path string, ser *Serializer) (*DataFile, error) {
	df := DataFile{}
	df.path = path
	df.num_records = 0
	df.header_size = uint32(unsafe.Sizeof(df.header_size))
	df.ser = ser
	df.fd = nil
	return &df, nil
}

func (df *DataFile) writeHeader() error {
	_, err := df.fd.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("writeHeader: failed to seek file start: %s", err.Error())
	}

	var buf bytes.Buffer
	err = binary.Write(&buf, binary.BigEndian, df.record_size)
	if err != nil {
		return fmt.Errorf("writeHeader: buffer write failed: %s", err.Error())
	}

	_, err = df.fd.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("writeHeader: header write failed: %s", err.Error())
		return err
	}

	err = df.fd.Sync()
	if err != nil {
		return fmt.Errorf("writeHeader: sync failed: %s", err.Error())
	}

	return nil
}

func (df *DataFile) parseHeader() error {
	_, err := df.fd.Seek(0, 0)
	if err != nil {
		return err
	}

	data := make([]byte, df.header_size)
	var n int
	n, err = df.fd.Read(data)
	if err != nil {
		return err
	}

	if uint32(n) != df.header_size {
		return fmt.Errorf("bytes read (%d) doesn't match header size (%d)", n, df.header_size)
	}

	buf := bytes.NewReader(data)
	err = binary.Read(buf, binary.BigEndian, &df.record_size)
	if err != nil {
		return err
	}

	return nil
}

func (df *DataFile) Close() error {
	if df.fd == nil {
		return nil
	}
	fd := df.fd
	df.fd = nil
	return fd.Close()
}

// File offset of the idx'th record
func (df *DataFile) getOffset(idx uint32) int64 {
	return int64(df.header_size + idx*df.record_size)
}

func (df *DataFile) Write(idx uint32, p *Point) error {
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

func (df *DataFile) Read(idx uint32) (*Point, error) {
	offset := df.getOffset(idx)
	_, err := df.fd.Seek(offset, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to seek to offset %d to read %d bytes for index %d: %s",
			offset, df.record_size, idx, err.Error())
	}

	data := make([]byte, df.record_size)
	var n int
	n, err = df.fd.Read(data)
	if err != nil {
		return nil, fmt.Errorf("failed to read %d bytes from position %d for index %d: %s",
			df.record_size, offset, idx, err.Error())
	}

	if uint32(n) != df.record_size {
		return nil, fmt.Errorf("bytes read (%d) doesn't match record size (%d)",
			n, df.record_size)
	}

	var p *Point
	p, err = df.ser.Deserialize(data)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize %d bytes fo data for index %d: %s",
			len(data), idx, err.Error())
	}

	// if the timestamp is 0 then we consider this an invalid read because
	// we must have written sparse points to the file
	if p.ts.UnixMicro() == 0 {
		return nil, fmt.Errorf("read empty timestamp at index %d: %s", idx, p.String())
	}

	return p, nil
}
