package log

import (
	"log"
	"net/http"

	"go.uber.org/zap"
)

type (
	// ResponseData Struct to save data from http.ResponseWriter.
	ResponseData struct {
		statusCode int
		contentLen int
	}

	// ResponseWriter Copies http.ResponseWriter functionality but saves some data for further use.
	ResponseWriter struct {
		*ResponseData
		http.ResponseWriter
	}
)

// Write Same as ResponseWriter.Write but saves content length.
func (rw ResponseWriter) Write(b []byte) (int, error) {
	contentLen, err := rw.ResponseWriter.Write(b)
	rw.ResponseData.contentLen = contentLen

	return contentLen, err
}

// WriteHeader Same as ResponseWriter.WriteHeader but saves status code.
func (rw ResponseWriter) WriteHeader(statusCode int) {
	rw.ResponseData.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// GetStatusCode Return saved status code.
func (rw ResponseWriter) GetStatusCode() int {
	return rw.ResponseData.statusCode
}

// GetContentLen Returns saves contentLen.
func (rw ResponseWriter) GetContentLen() int {
	return rw.ResponseData.contentLen
}

// NewResponseWriter create ResponseWriter.
func NewResponseWriter(w http.ResponseWriter) ResponseWriter {
	rd := ResponseData{}
	return ResponseWriter{
		&rd,
		w,
	}
}

// NewLogger Creates new logger.
func NewLogger(logPaths []string) *zap.Logger {
	cfg := zap.NewProductionConfig()

	cfg.OutputPaths = logPaths

	zlogger, err := cfg.Build()
	if err != nil {
		log.Fatalln("Could not create http logger: ", err)
	}

	return zlogger
}
