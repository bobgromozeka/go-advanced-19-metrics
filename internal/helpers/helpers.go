package helpers

import (
	"bytes"
	"compress/gzip"
	"log"
	"strconv"
)

func StrToInt(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		log.Println("Error when converting string to int: ", err)
		return 0 //ignore error and return 0
	}

	return v
}

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
