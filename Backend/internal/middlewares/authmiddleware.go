package middlewares

import (
	"context"
	"net/http"
	"strings"

	"ticketapp/internal/services"
)

type ctxKey string

const RoleKey ctxKey = "role"
const UserIDKey ctxKey = "user_id"

func AuthMiddleware(jwtSvc *services.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := jwtSvc.Validate(token)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), RoleKey, claims["role"])
			ctx = context.WithValue(ctx, UserIDKey, claims["sub"])

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
