package helpers

import (
	"bytes"
	"compress/gzip"
	"io"
	"strings"
	"testing"
)

func TestGzipBytes(t *testing.T) {
	str := strings.Repeat("Test string", 20)
	gzippedString, _ := GzipBytes([]byte(str))

	gzr, _ := gzip.NewReader(bytes.NewReader(gzippedString))
	unzippedStr, _ := io.ReadAll(gzr)
	gzr.Close()

	if string(unzippedStr) != str {
		t.Errorf("unzipped string should be equal with initial")
	}
}
