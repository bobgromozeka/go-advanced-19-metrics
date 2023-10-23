package middlewares

import (
	"bytes"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"io"
	"log"
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

				parsedPrivateKey, err := x509.ParsePKCS8PrivateKey(privateKeyBlock.Bytes)
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				data, readErr := io.ReadAll(r.Body)

				if readErr != nil {
					log.Println(readErr)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				decrypted, decryptErr := decrypt(data, parsedPrivateKey.(*rsa.PrivateKey))

				if decryptErr != nil {
					log.Println(decryptErr)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				r.Body = newBytesReaderCloser(decrypted)

				next.ServeHTTP(w, r)
			},
		)
	}
}

func decrypt(data []byte, k *rsa.PrivateKey) ([]byte, error) {
	h := sha256.New()
	step := k.PublicKey.Size()
	res := make([]byte, 0)

	for i := 0; i < len(data); i += step {
		end := i + step
		if end > len(data) {
			end = len(data)
		}

		decrypted, err := rsa.DecryptOAEP(h, nil, k, data[i:end], []byte("data"))
		if err != nil {
			return data, err
		}

		res = append(res, decrypted...)
	}

	return res, nil
}
