package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"

	"github.com/imarrche/jwt-auth-example/internal/logger"
)

// loggerMiddleware is middleware that logs every request with zap logger.
func (s *Server) loggerMiddleware() func(next http.Handler) http.Handler {
	l := logger.Get()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t := time.Now()
			defer func() {
				l.Info("served",
					zap.String("protocol", r.Proto),
					zap.String("path", r.URL.Path),
					zap.Int("status", ww.Status()),
					zap.Int("size", ww.BytesWritten()),
					zap.Duration("latency", time.Since(t)),
				)
			}()

			next.ServeHTTP(ww, r)
		})
	}
}

// authMiddleware is middleware for JWT authorization.
func (s *Server) authMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			if !strings.HasPrefix(h, "Bearer ") {
				s.error(w, r, http.StatusUnauthorized, nil)
				return
			}
			if _, err := s.service.Auth().ValidateJWT(h[7:], "access"); err != nil {
				s.error(w, r, http.StatusUnauthorized, nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
