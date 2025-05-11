package goraw

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// Decoder defines the interface to implement a new raw decoder
type Decoder interface {
	ExifReaderAt() (io.ReaderAt, error)
	Close() error
}

type format struct {
	name   string
	magic  []byte
	decode func(io.ReaderAt) (Decoder, error)
}

var formats []format
var maxMagicLength int

// Open instanciates a decoder from a file by detecting its format
func Open(filename string) (Decoder, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return New(fd)
}

func Close(r Decoder) error {
	if r == nil {
		return fmt.Errorf("decoder is nil")
	}
	return r.Close()
}

// New instanciates a decoder from a ReaderAt by detecting its format
func New(r io.ReaderAt) (Decoder, error) {
	magic := make([]byte, maxMagicLength)
	_, err := r.ReadAt(magic, 0)
	if err != nil {
		return nil, err
	}
	for _, format := range formats {
		if bytes.Compare(format.magic, magic[:len(format.magic)]) == 0 {
			return format.decode(r)
		}
	}
	return nil, fmt.Errorf("unsupported format")
}

func init() {
	maxMagicLength = 0
	formats = []format{}
}

// RegisterFormat registers a new raw format
func RegisterFormat(name string, magic []byte, decode func(io.ReaderAt) (Decoder, error)) {
	if len(magic) > maxMagicLength {
		maxMagicLength = len(magic)
	}
	formats = append(formats, format{
		name:   name,
		magic:  magic,
		decode: decode,
	})
}
