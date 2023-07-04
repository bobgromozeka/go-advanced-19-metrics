package gzip

import (
	"compress/gzip"
	"io"
	"net/http"
)

const Name = "gzip"

type Writer struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func (gzw *Writer) Header() http.Header {
	return gzw.w.Header()
}

func (gzw *Writer) Write(p []byte) (int, error) {
	return gzw.zw.Write(p)
}

func (gzw *Writer) WriteHeader(statusCode int) {
	if statusCode < 300 {
		gzw.w.Header().Set("Content-Encoding", Name)
	}
	gzw.w.WriteHeader(statusCode)
}

func (gzw *Writer) Close() error {
	return gzw.zw.Close()
}

type Reader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func (gzr *Reader) Read(p []byte) (int, error) {
	return gzr.zr.Read(p)
}

func (gzr *Reader) Close() error {
	return gzr.zr.Close()
}

func NewGzipWriter(w http.ResponseWriter) *Writer {
	gzw := gzip.NewWriter(w)

	return &Writer{
		w,
		gzw,
	}
}

func NewGzipReader(r io.ReadCloser) (*Reader, error) {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &Reader{
		r,
		gzr,
	}, nil
}
