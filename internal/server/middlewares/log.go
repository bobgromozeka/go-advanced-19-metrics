package middlewares

import (
	"log"
	"net/http"
	"time"

	"go.uber.org/zap"
)

var logger = zap.NewNop()
var logPath = "./http.log"

type (
	responseData struct {
		statusCode int
		contentLen int
	}

	logResponseWriter struct {
		*responseData
		http.ResponseWriter
	}
)

func (l logResponseWriter) Write(b []byte) (int, error) {
	contentLen, err := l.ResponseWriter.Write(b)
	l.responseData.contentLen = contentLen

	return contentLen, err
}

func (l logResponseWriter) WriteHeader(statusCode int) {
	l.responseData.statusCode = statusCode
	l.ResponseWriter.WriteHeader(statusCode)
}

func init() {
	cfg := zap.NewProductionConfig()

	cfg.OutputPaths = []string{
		logPath,
	}

	zlogger, err := cfg.Build()
	if err != nil {
		log.Fatalln("Could not create http logger: ", err)
	}

	logger = zlogger
}

func WithLogging(f http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()

		rd := responseData{}
		lw := logResponseWriter{
			&rd,
			w,
		}

		f.ServeHTTP(lw, r)
		requestTime := time.Since(timeStart)

		logger.Info("Got request",
			zap.String("method", r.Method),
			zap.String("endpoint", r.URL.Path),
			zap.Duration("duration", requestTime),
			zap.Int("status code", rd.statusCode),
			zap.Int("content length", rd.contentLen),
		)
	}

	return http.HandlerFunc(logFn)
}
