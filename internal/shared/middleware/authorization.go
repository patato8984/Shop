package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/patato8984/Shop/internal/shared/dto"
	"github.com/patato8984/Shop/internal/shared/logger"
	"github.com/patato8984/Shop/pkg/ctxkey"
	"go.uber.org/zap"
)

type Authorization struct {
	key string
}

func NewAuthorization(key string) *Authorization {
	return &Authorization{key: key}
}
func (a Authorization) Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromContext(r.Context())
		auth := r.Header.Get("Authorization")
		parts := strings.Split(auth, " ")
		if parts[0] != "Bearer" || len(parts) != 2 {
			log.Error("error jwt",
				zap.Error(errors.New("there is no token in the header")),
			)
			dto.Response(w, "error", http.StatusBadRequest, nil)
			return
		}
		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("error jwt")
			}
			return []byte(a.key), nil
		})
		if err != nil || token == nil {
			log.Warn("error jwt parse")
			http.Error(w, "error jwt token", http.StatusUnauthorized)
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			idUser, ok := claims["id"].(float64)
			if !ok {
				log.Warn("error jwt parse",
					zap.Error(errors.New("getting an ID")),
				)
				http.Error(w, "invalid token claims", http.StatusUnauthorized)
				return
			}
			role, ok := claims["role"].(string)
			if !ok || role == "" {
				log.Warn("error jwt parse",
					zap.Error(errors.New("getting an ROLE")),
				)
				http.Error(w, "invalid token claims", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), ctxkey.UserIDKey, int(idUser))
			ctx = context.WithValue(ctx, ctxkey.Role, role)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		log.Error("error jwt",
			zap.Error(errors.New("invalid token")),
		)
		dto.Response(w, "error", http.StatusBadRequest, nil)
		return
	})
}
func (a Authorization) AuthenticationAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromContext(r.Context())
		role := r.Context().Value(ctxkey.Role)
		if role != "admin" {
			log.Error("error jwt",
				zap.Error(errors.New("attempt to gain access without the admin role")),
				zap.String("ip", r.RemoteAddr),
			)
			dto.Response(w, "error role", http.StatusForbidden, nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}
