package fuji_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tjamet/goraw/fuji"
	"github.com/tjamet/goraw/test-tools"
)

func TestNonFujiFileReportsError(t *testing.T) {
	testFile := tools.DownloadRAW("http://www.rawsamples.ch/raws/canon/RAW_CANON_EOS_5DS.CR2")
	_, err := fuji.Open(testFile)
	assert.Error(t, err)
}

func TestFujiFileIsAccepted(t *testing.T) {
	testFile := tools.DownloadRAW("http://www.rawsamples.ch/raws/fuji/RAW_FUJI_FINEPIX_X100.RAF")
	_, err := fuji.Open(testFile)
	assert.NoError(t, err)
}

func TestFujiFileCanGetExif(t *testing.T) {
	testFile := tools.DownloadRAW("http://www.rawsamples.ch/raws/fuji/RAW_FUJI_FINEPIX_X100.RAF")
	r, _ := fuji.Open(testFile)
	reader, _ := r.ExifReaderAt()
	header := make([]byte, 2)
	_, err := reader.ReadAt(header, 0)
	assert.NoError(t, err)
	// fujifilm is big endian
	assert.Equal(t, []byte("II"), header)
}
