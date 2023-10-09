package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/bobgromozeka/metrics/internal/log"
)

// WithLogging Adds access middleware for http handlers. Logs handler method, endpoint path, request duration, status code and content length.
func WithLogging(logPaths []string) func(handler http.Handler) http.Handler {
	if len(logPaths) < 1 {
		fmt.Println("No log paths specified, skipping WithLogging middleware")
		return func(f http.Handler) http.Handler {
			return f
		}
	}

	logger := log.NewLogger(logPaths)
	return func(next http.Handler) http.Handler {
		return loggingHandler(next, logger)
	}
}

func loggingHandler(next http.Handler, logger *zap.Logger) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()

		lw := log.NewResponseWriter(w)

		next.ServeHTTP(lw, r)
		requestTime := time.Since(timeStart)

		logger.Info(
			"Got request",
			zap.String("method", r.Method),
			zap.String("endpoint", r.URL.Path),
			zap.Duration("duration", requestTime),
			zap.Int("status code", lw.GetStatusCode()),
			zap.Int("content length", lw.GetContentLen()),
		)
	}

	return http.HandlerFunc(logFn)
}
