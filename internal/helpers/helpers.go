package helpers

import (
	"bytes"
	"compress/gzip"
	"log"
	"net/http"
	"strconv"

	"github.com/bobgromozeka/metrics/internal/hash"
)

// StrToInt parses string into int or returns 0
func StrToInt(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		log.Println("Error when converting string to int: ", err)
		return 0 //ignore error and return 0
	}

	return v
}

// Gzip compresses provided bytes or return error if occurs.
func Gzip(b []byte) ([]byte, error) {
	var resBuf bytes.Buffer
	gzw := gzip.NewWriter(&resBuf)
	_, err := gzw.Write(b)
	gzw.Close()

	if err != nil {
		return resBuf.Bytes(), err
	}

	return resBuf.Bytes(), nil
}

// SignResponse hashes provided bytes and adds them as header to http.ResponseWriter
func SignResponse(w http.ResponseWriter, body []byte, key string, header string) {
	if key != "" {
		h := hash.New(key)
		w.Header().Set(header, h.Sha256(string(body)))
	}
}
