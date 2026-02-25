package middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/patato8984/Shop/internal/shared/logger"
	"go.uber.org/zap"
)

type Logger struct {
	l *zap.Logger
}
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggerMiddleware(l *zap.Logger) *Logger {
	return &Logger{l: l}
}
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
func (l Logger) RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId := uuid.NewString()
		w.Header().Set("X-Request-ID", reqId)
		l := l.l.With(zap.String("request_id", reqId))
		ctx := logger.ToContext(r.Context(), l)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func (l Logger) Response(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		times := time.Now()
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		next.ServeHTTP(rw, r)
		duration := time.Since(times)
		log := logger.FromContext(r.Context())
		switch {
		case rw.statusCode <= 299:
			log.Info("response",
				zap.String("proto", r.Proto),
				zap.String("method", r.Method),
				zap.String("url", r.URL.Path),
				zap.Int("status", rw.statusCode),
				zap.Duration("latency", duration),
			)
		case rw.statusCode <= 499:
			log.Warn("response bad request",
				zap.String("proto", r.Proto),
				zap.String("method", r.Method),
				zap.String("url", r.URL.Path),
				zap.Int("status", rw.statusCode),
				zap.Duration("latency", duration),
			)
		default:
			log.Error("response server error",
				zap.String("proto", r.Proto),
				zap.String("method", r.Method),
				zap.String("url", r.URL.Path),
				zap.Int("status", rw.statusCode),
				zap.Duration("latency", duration),
			)
		}
	})
}
