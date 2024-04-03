package tools

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const testDir = "tests/raw"

func DownloadRAW(link string) string {
	destFile := testDir + strings.TrimPrefix(link, "http://www.rawsamples.ch/raws")
	destFolder := filepath.Dir(destFile)
	os.MkdirAll(destFolder, 0777)
	_, err := os.Stat(destFile)
	if os.IsNotExist(err) {
		resp, err := http.Get(link)
		if err != nil {
			log.Printf("Failed to download image %s: %s", link, err)
		} else {
			defer resp.Body.Close()
			fd, err := os.Create(destFile)
			if err != nil {
				log.Printf("Failed to download image %s: %s", link, err)
			} else {
				io.Copy(fd, resp.Body)
			}
		}
	}
	return destFile
}
