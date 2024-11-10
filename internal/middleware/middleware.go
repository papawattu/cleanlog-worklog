package middleware

import (
	"context"
	"log/slog"
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

		userId := r.Context().Value("user")

		if userId == nil {
			slog.Info("Request", "Method", r.Method, "URL", r.URL.Path, slog.Duration("Taken", time.Since(start)), "Status", wrapped.statusCode, slog.String("User", "Anonymous"))
		} else {
			slog.Info("Request", "Method", r.Method, "URL", r.URL.Path, slog.Duration("Taken", time.Since(start)), "Status", wrapped.statusCode, slog.Int("User", userId.(int)))
		}

	})
}
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("Recovered from panic", "Error", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func Authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")
		if token != "Bearer mytoken" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user", 0)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
