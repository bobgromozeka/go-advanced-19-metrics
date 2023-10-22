package middlewares

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"net/http"

	"github.com/bobgromozeka/metrics/internal"
)

type bytesReaderCloser struct {
	*bytes.Reader
}

func newBytesReaderCloser(data []byte) *bytesReaderCloser {
	return &bytesReaderCloser{
		bytes.NewReader(data),
	}
}

func (brc *bytesReaderCloser) Close() error {
	return nil
}

func Rsa(key []byte) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get(internal.RSAEncryptedHeader) != "true" || len(key) < 1 {
					next.ServeHTTP(w, r)
					return
				}

				privateKeyBlock, _ := pem.Decode(key)

				parsedPrivateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				data, readErr := io.ReadAll(r.Body)
				if readErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				decrypted, decryptErr := rsa.DecryptPKCS1v15(rand.Reader, parsedPrivateKey, data)
				if decryptErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				r.Body = newBytesReaderCloser(decrypted)

				next.ServeHTTP(w, r)
			},
		)
	}
}
