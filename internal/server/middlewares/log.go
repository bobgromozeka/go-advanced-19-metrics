package middlewares

import (
	"log"
	"net/http"
	"time"

	"go.uber.org/zap"
)

var logger = zap.NewNop()
var sugaredLogger *zap.SugaredLogger
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
	sugaredLogger = logger.Sugar()
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
		timeEnd := time.Since(timeStart)

		sugaredLogger.Infof("Got [%s] request on endpoint %s and it took %s. Status code = %d, content length = %d", r.Method, r.URL.Path, timeEnd, rd.statusCode, rd.contentLen)
	}

	return http.HandlerFunc(logFn)
}
