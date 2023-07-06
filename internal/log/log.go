package log

import (
	"log"
	"net/http"

	"go.uber.org/zap"
)

type (
	ResponseData struct {
		statusCode int
		contentLen int
	}

	ResponseWriter struct {
		*ResponseData
		http.ResponseWriter
	}
)

func (rw ResponseWriter) Write(b []byte) (int, error) {
	contentLen, err := rw.ResponseWriter.Write(b)
	rw.ResponseData.contentLen = contentLen

	return contentLen, err
}

func (rw ResponseWriter) WriteHeader(statusCode int) {
	rw.ResponseData.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw ResponseWriter) GetStatusCode() int {
	return rw.ResponseData.statusCode
}

func (rw ResponseWriter) GetContentLen() int {
	return rw.ResponseData.contentLen
}

func NewResponseWriter(w http.ResponseWriter) ResponseWriter {
	rd := ResponseData{}
	return ResponseWriter{
		&rd,
		w,
	}
}

func NewLogger(logPaths []string) *zap.Logger {
	cfg := zap.NewProductionConfig()

	cfg.OutputPaths = logPaths

	zlogger, err := cfg.Build()
	if err != nil {
		log.Fatalln("Could not create http logger: ", err)
	}

	return zlogger
}
