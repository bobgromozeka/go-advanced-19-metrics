package helpers

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

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

func SetupGracefulShutdown(cancel context.CancelFunc) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig
		fmt.Println("Stopping application.....")
		cancel()
	}()
}
