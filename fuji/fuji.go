package fuji

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/tjamet/goraw"
	"github.com/tjamet/goraw/io"
	"github.com/tjamet/goraw/jpeg"
)

var magic = []byte("FUJIFILMCCD-RAW")

// Raw holds the context to decode a fujifilm raw file
type Raw struct {
	readerAt   io.ReaderAt
	version    string
	camera     string
	jpegStart  uint32
	jpegLength uint32
}

// Open instanciates a fujifilm raw handler from a file
func Open(filename string) (*Raw, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return New(fd)
}

// New instanciates a fujifilm raw handler from ReaderAt
func New(r io.ReaderAt) (*Raw, error) {
	// https://libopenraw.freedesktop.org/wiki/Fuji_RAF/
	header := make([]byte, 108)
	_, err := r.ReadAt(header, 0)
	if err != nil {
		return nil, err
	}
	if bytes.Compare(header[:len(magic)], magic) != 0 {
		return nil, fmt.Errorf("the content is not a fujifilm raw file. Magic mismatch")
	}
	return &Raw{
		readerAt:   r,
		version:    string(header[16:20]),
		camera:     string(header[24 : 24+bytes.Index(header[24:24+32], []byte{'\x00'})]),
		jpegStart:  binary.BigEndian.Uint32(header[84:88]),
		jpegLength: binary.BigEndian.Uint32(header[88:92]),
	}, nil
}

// ExifReaderAt returns a direct reader to the Exif inside the raw
func (r *Raw) ExifReaderAt() (io.ReaderAt, error) {
	j, err := jpeg.New(gorawio.NewReaderAt(r.readerAt, int64(r.jpegStart)))
	if err != nil {
		return nil, err
	}
	return j.ExifReaderAt()
}

func init() {
	goraw.RegisterFormat("fujifilm", magic, func(r io.ReaderAt) (goraw.Decoder, error) { return New(r) })
}
