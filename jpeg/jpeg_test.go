package jpeg_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tjamet/goraw/jpeg"
)

func TestJpegClose(t *testing.T) {
	// Use a small local JPEG file or create a dummy one for the test
	file, err := os.CreateTemp("", "testfile-*.jpg")
	assert.NoError(t, err)
	defer os.Remove(file.Name())
	_, err = file.Write([]byte{0xff, 0xd8, 0xff, 0xe1, 0x00, 0x10, 'E', 'x', 'i', 'f', 0x00, 0x00})
	assert.NoError(t, err)
	file.Close()

	f, err := os.Open(file.Name())
	assert.NoError(t, err)
	jpegDecoder, err := jpeg.New(f)
	assert.NoError(t, err)

	// Test normal close
	err = jpegDecoder.Close()
	assert.NoError(t, err)

	// Test double close (should not error)
	err = jpegDecoder.Close()
	assert.NoError(t, err)

	// Test use after close (should error)
	_, err = jpegDecoder.ExifReaderAt()
	assert.Error(t, err)

	// Test nil receiver close (should not error)
	var nilJpeg *jpeg.JPEG
	err = nilJpeg.Close()
	assert.NoError(t, err)
}
