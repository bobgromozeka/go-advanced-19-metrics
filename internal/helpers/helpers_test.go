package helpers

import (
	"bytes"
	"compress/gzip"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzip(t *testing.T) {
	str := strings.Repeat("Test string", 20)
	gzippedString, _ := Gzip([]byte(str))

	gzr, _ := gzip.NewReader(bytes.NewReader(gzippedString))
	unzippedStr, _ := io.ReadAll(gzr)
	gzr.Close()

	if string(unzippedStr) != str {
		t.Errorf("unzipped string should be equal with initial")
	}
}

func TestSignResponse(t *testing.T) {
	assert.Equal(t, 123, StrToInt("123"))
	assert.Equal(t, 0, StrToInt("123sdf"))
}
