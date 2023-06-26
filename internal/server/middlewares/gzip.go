package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

const Gzip = "gzip"

type gzipWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func (gzw *gzipWriter) Header() http.Header {
	return gzw.w.Header()
}

func (gzw *gzipWriter) Write(p []byte) (int, error) {
	return gzw.zw.Write(p)
}

func (gzw *gzipWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		gzw.w.Header().Set("Content-Encoding", Gzip)
	}
	gzw.w.WriteHeader(statusCode)
}

func (gzw *gzipWriter) Close() error {
	return gzw.zw.Close()
}

type gzipReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func (gzr *gzipReader) Read(p []byte) (int, error) {
	return gzr.zr.Read(p)
}

func (gzr *gzipReader) Close() error {
	return gzr.zr.Close()
}

func Gzippify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resultWriter := w
		acceptList := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptList, Gzip)
		if supportsGzip {
			gzw := newGzipWriter(w)
			resultWriter = gzw
			defer gzw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		gotGzip := strings.Contains(contentEncoding, Gzip)
		if gotGzip {
			gzr, err := newGzipReader(r.Body)
			if err != nil {
				resultWriter.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = gzr
			defer gzr.Close()
		}

		next.ServeHTTP(resultWriter, r)
	})
}

func newGzipWriter(w http.ResponseWriter) *gzipWriter {
	gzw := gzip.NewWriter(w)

	return &gzipWriter{
		w,
		gzw,
	}
}

func newGzipReader(r io.ReadCloser) (*gzipReader, error) {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &gzipReader{
		r,
		gzr,
	}, nil
}
