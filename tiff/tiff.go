package tiff

import (
	"io"
	"os"

	"github.com/tjamet/goraw"
)

// Raw holds the context to decode a TIFF based raw file
type Raw struct {
	readerAt io.ReaderAt
}

// Open instanciates a TIFF based raw handler from a file
func Open(filename string) (*Raw, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return &Raw{readerAt: fd}, nil
}

// New instanciates a TIFF based raw handler from ReaderAt
func New(r io.ReaderAt) (*Raw, error) {
	return &Raw{readerAt: r}, nil
}

// ExifReaderAt returns a direct reader to the Exif inside the raw
func (r *Raw) ExifReaderAt() (io.ReaderAt, error) {
	return r.readerAt, nil
}

func init() {
	goraw.RegisterFormat("tiff", []byte("MM"), func(r io.ReaderAt) (goraw.Decoder, error) { return New(r) })
	goraw.RegisterFormat("tiff", []byte("II"), func(r io.ReaderAt) (goraw.Decoder, error) { return New(r) })
}
