package jpeg

import (
	"bytes"
	"fmt"
	"io"

	"github.com/tjamet/goraw"
	"github.com/tjamet/goraw/io"
)

const jpegAPP1 = 0xE1

// JPEG holds the context to handle JPEG within a RAW
type JPEG struct {
	readerAt io.ReaderAt
}

// New instanciates a handler for JPEG within a RAW
func New(r io.ReaderAt) (*JPEG, error) {
	return &JPEG{
		readerAt: r,
	}, nil
}

// ExifReaderAt returns a direct reader to the Exif inside the raw
func (r *JPEG) ExifReaderAt() (io.ReaderAt, error) {
	offset, err := findJPEGExifOffset(jpegAPP1, r.readerAt)
	if err != nil {
		return nil, fmt.Errorf("could not find Exif data: %s", err.Error())
	}
	return gorawio.NewReaderAt(r.readerAt, int64(offset)), nil
}

// findJPEGExifOffset finds marker in r and returns the offset of the Exif data
func findJPEGExifOffset(marker byte, r io.ReaderAt) (uint32, error) {
	// read the file 1MB per 1MB to locate the marker
	bufferLength := 4096
	// Ensure there is a small overlap between buffers to handle the case
	// where the marker is located exactly in-between 2 reads
	seeker := make([]byte, bufferLength+2)
	appHeader := []byte{0xFF, marker}

	// seek to marker
	var markerIndex int
	i := 0
	for {
		count, err := r.ReadAt(seeker, int64(i*bufferLength))
		if err != nil && err != io.EOF {
			return 0, err
		}

		markerIndex = bytes.Index(seeker, appHeader)
		// stop at the end of the file or when the marker is found
		if markerIndex >= 0 || count < bufferLength {
			break
		}
		i++
	}
	if markerIndex < 0 {
		return 0, fmt.Errorf("Unable to find the JPEG Exif marker")
	}
	// Skip the content length
	markerIndex += i*bufferLength + 4
	exifHeader := make([]byte, 6)
	_, err := r.ReadAt(exifHeader, int64(markerIndex))
	if err != nil {
		return 0, err
	}
	if !bytes.Equal(exifHeader, []byte("Exif\000\000")) {
		return 0, fmt.Errorf("Unable to find the JPEG Exif marker")
	}
	return uint32(markerIndex + 6), nil
}

func init() {
	goraw.RegisterFormat("jpeg", []byte{0xff, 0xd8}, func(r io.ReaderAt) (goraw.Decoder, error) { return New(r) })
}
