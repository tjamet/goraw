package goraw_test

import (
	"fmt"

	"github.com/tjamet/goraw"
	_ "github.com/tjamet/goraw/fuji"
	tools "github.com/tjamet/goraw/test-tools"
)

func ExampleOpen() {
	testFile := tools.DownloadRAW("http://www.rawsamples.ch/raws/fuji/RAW_FUJI_FINEPIX_X100.RAF")
	decoder, _ := goraw.Open(testFile)
	reader, _ := decoder.ExifReaderAt()
	order := make([]byte, 2)
	reader.ReadAt(order, 0)
	fmt.Println(string(order))
	// Output: II
}
