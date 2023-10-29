package middlewares

import (
	"log"
	"net"
	"net/http"

	"github.com/bobgromozeka/metrics/internal"
)

func TrustedSubnet(subnet string) func(next http.Handler) http.Handler {
	funcWithoutMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			},
		)
	}

	if subnet == "" {
		return funcWithoutMiddleware
	}

	_, ipNet, parseErr := net.ParseCIDR(subnet)
	if parseErr != nil {
		log.Println("Could not parse trusted subnet: ", parseErr)
		return funcWithoutMiddleware
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				realIPs := r.Header.Values(internal.RealIpHeader)
				if len(realIPs) > 0 {
					clientIP := net.ParseIP(realIPs[0])
					if !ipNet.Contains(clientIP) {
						w.WriteHeader(http.StatusForbidden)
						return
					}
				}
				next.ServeHTTP(w, r)
			},
		)
	}
}
