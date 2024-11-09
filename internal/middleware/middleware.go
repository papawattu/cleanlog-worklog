package middleware

import (
	"log"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

type wrappedResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func CreateMiddleware(mw ...Middleware) Middleware {
	return func(h http.Handler) http.Handler {
		for i := range mw {
			h = mw[len(mw)-1-i](h)
		}
		return h
	}
}
func (w *wrappedResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedResponseWriter{w, http.StatusOK}

		next.ServeHTTP(wrapped, r)
		log.Printf("Request: %s %s processed in %s : Status code %d", r.Method, r.URL.Path, time.Since(start), wrapped.statusCode)
	})
}
