package fuji

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/tjamet/goraw"
	gorawio "github.com/tjamet/goraw/io"
	"github.com/tjamet/goraw/jpeg"
)

const (
	uint32Length     = 4
	headerLength     = 108
	jpegStartOffset  = 84
	jpegLengthOffset = 88
	versionOffset    = 16
	versionLength    = 4
	cameraOffset     = 24
	cameraLength     = 32
)

var magic = []byte("FUJIFILMCCD-RAW")

// Raw holds the context to decode a fujifilm raw file
type Raw struct {
	fd         *os.File
	readerAt   io.ReaderAt
	version    string
	camera     string
	jpegStart  uint32
	jpegLength uint32
	jpegReader *jpeg.JPEG
}

// Open instanciates a fujifilm raw handler from a file
func Open(filename string) (*Raw, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	raw, err := New(fd)
	if err != nil {
		return nil, err
	}
	raw.fd = fd
	return raw, nil
}

func Close(r *Raw) error {
	return r.Close()
}

// New instanciates a fujifilm raw handler from ReaderAt
func New(r io.ReaderAt) (*Raw, error) {
	// https://libopenraw.freedesktop.org/wiki/Fuji_RAF/
	header := make([]byte, headerLength)
	_, err := r.ReadAt(header, 0)
	if err != nil {
		return nil, err
	}
	if bytes.Compare(header[:len(magic)], magic) != 0 {
		return nil, fmt.Errorf("the content is not a fujifilm raw file. Magic mismatch")
	}
	return &Raw{
		readerAt:   r,
		version:    string(header[versionOffset : versionOffset+versionLength]),
		camera:     string(header[cameraOffset : cameraOffset+cameraLength]),
		jpegStart:  binary.BigEndian.Uint32(header[jpegStartOffset : jpegStartOffset+uint32Length]),
		jpegLength: binary.BigEndian.Uint32(header[jpegLengthOffset : jpegLengthOffset+uint32Length]),
	}, nil
}

// ExifReaderAt returns a direct reader to the Exif inside the raw
func (r *Raw) ExifReaderAt() (io.ReaderAt, error) {
	if r == nil {
		return nil, fmt.Errorf("raw is nil")
	}
	if r.readerAt == nil {
		return nil, os.ErrClosed
	}
	if r.jpegReader == nil {
		j, err := jpeg.New(gorawio.NewReaderAt(r.readerAt, int64(r.jpegStart)))
		if err != nil {
			return nil, err
		}
		r.jpegReader = j
	}
	return r.jpegReader.ExifReaderAt()
}

func (r *Raw) ExifOffset() (int64, error) {
	if r == nil {
		return 0, fmt.Errorf("raw is nil")
	}
	if r.jpegReader == nil {
		j, err := jpeg.New(gorawio.NewReaderAt(r.readerAt, int64(r.jpegStart)))
		if err != nil {
			return 0, err
		}
		r.jpegReader = j
	}
	jpegExifOffset, err := r.jpegReader.ExifOffset()
	if err != nil {
		return 0, err
	}
	return jpegExifOffset + int64(r.jpegStart), nil
}

func (r *Raw) Close() error {
	if r == nil {
		return nil
	}
	if r.fd != nil {
		err := r.fd.Close()
		if err != nil {
			return err
		}
		r.fd = nil
	}
	r.readerAt = nil
	return nil
}

func init() {
	goraw.RegisterFormat("fujifilm", magic, func(r io.ReaderAt) (goraw.Decoder, error) { return New(r) })
}
