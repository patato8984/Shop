package middleware

import (
	"context"
	"net/http"
	"time"
)

type CreatedCtx struct {
	timeoutCtxConfig map[string]time.Duration
}

func NewCreatedCtx(timeoutCtxConfig map[string]time.Duration) *CreatedCtx {
	return &CreatedCtx{timeoutCtxConfig: timeoutCtxConfig}
}
func (c CreatedCtx) CreatedCtxTimeout(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeOut := time.Second * 3
		ctx := context.Background()
		if customConfig, ok := c.timeoutCtxConfig[r.Pattern]; ok {
			timeOut = customConfig
		}
		ctx, cancel := context.WithTimeout(ctx, timeOut)
		defer cancel()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
