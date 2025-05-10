package tiff_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	tools "github.com/tjamet/goraw/test-tools"
	"github.com/tjamet/goraw/tiff"
)

func TestTiffClose(t *testing.T) {
	testFile := tools.DownloadRAW("http://www.rawsamples.ch/raws/nikon/RAW_NIKON_D1.NEF")
	r, err := tiff.Open(testFile)
	assert.NoError(t, err)

	// Test normal close
	err = r.Close()
	assert.NoError(t, err)

	// Test double close (should not error)
	err = r.Close()
	assert.NoError(t, err)

	// Test use after close (should error)
	_, err = r.ExifReaderAt()
	assert.Error(t, err)

	// Test nil receiver close (should not error)
	var nilRaw *tiff.Raw
	err = nilRaw.Close()
	assert.NoError(t, err)
}
