package goraw_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tjamet/goraw"
	_ "github.com/tjamet/goraw/fuji"
	tools "github.com/tjamet/goraw/test-tools"
)

func TestGenericClose(t *testing.T) {
	testFile := tools.DownloadRAW("http://www.rawsamples.ch/raws/fuji/RAW_FUJI_FINEPIX_X100.RAF")
	decoder, err := goraw.Open(testFile)
	assert.NoError(t, err)

	// Test normal close
	err = goraw.Close(decoder)
	assert.NoError(t, err)

	// Test double close (should not error)
	err = goraw.Close(decoder)
	assert.NoError(t, err)

	// Test use after close (should error)
	_, err = decoder.ExifReaderAt()
	assert.Error(t, err)

	// Test nil decoder close (should error)
	var nilDecoder goraw.Decoder
	err = goraw.Close(nilDecoder)
	assert.Error(t, err)
}
