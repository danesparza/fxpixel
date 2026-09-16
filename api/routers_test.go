package api

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterPreservesPeerAddress(t *testing.T) {
	for _, header := range []string{"True-Client-IP", "X-Real-IP", "X-Forwarded-For"} {
		t.Run(header, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/peer-address", nil)
			req.RemoteAddr = "192.0.2.1:1234"
			req.Header.Set(header, "198.51.100.2")
			router := NewRouter(Service{}).(*chi.Mux)
			var peer string
			router.Get("/peer-address", func(w http.ResponseWriter, r *http.Request) { peer = r.RemoteAddr })
			router.ServeHTTP(httptest.NewRecorder(), req)
			if peer != "192.0.2.1:1234" {
				t.Fatalf("untrusted %s changed peer address to %q", header, peer)
			}
		})
	}
}
