package tiff

import (
	"fmt"
	"io"
	"os"

	"github.com/tjamet/goraw"
)

// Raw holds the context to decode a TIFF based raw file
type Raw struct {
	fd       *os.File
	readerAt io.ReaderAt
}

// Open instanciates a TIFF based raw handler from a file
func Open(filename string) (*Raw, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return &Raw{readerAt: fd, fd: fd}, nil
}

func Close(r *Raw) error {
	return r.Close()
}

// New instanciates a TIFF based raw handler from ReaderAt
func New(r io.ReaderAt) (*Raw, error) {
	return &Raw{readerAt: r}, nil
}

// ExifReaderAt returns a direct reader to the Exif inside the raw
func (r *Raw) ExifReaderAt() (io.ReaderAt, error) {
	if r == nil {
		return nil, fmt.Errorf("raw is nil")
	}
	if r.readerAt == nil {
		return nil, os.ErrClosed
	}
	return r.readerAt, nil
}

func (r *Raw) Close() error {
	if r == nil {
		return nil
	}
	if r.fd != nil {
		r.readerAt = nil
		err := r.fd.Close()
		if err != nil {
			return err
		}
		r.fd = nil
	}
	return nil
}

func init() {
	goraw.RegisterFormat("tiff", []byte("MM"), func(r io.ReaderAt) (goraw.Decoder, error) { return New(r) })
	goraw.RegisterFormat("tiff", []byte("II"), func(r io.ReaderAt) (goraw.Decoder, error) { return New(r) })
}
