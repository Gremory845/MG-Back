package middleware

import (
	"context"
	"net/http"
	"strings"

	"mg-back/internal/infrastructure/security"
	"mg-back/internal/presentation/http/response"
)

type contextKey string

const userIDContextKey contextKey = "userID"

func unauthorized(w http.ResponseWriter, message string) {
	response.JSON(w, http.StatusUnauthorized, map[string]interface{}{
		"error": map[string]string{
			"code":    "UNAUTHORIZED",
			"message": message,
		},
	})
}

// Auth valida el access token (JWT) del header Authorization: Bearer <token>
// y deja el userID disponible en el contexto de la petición.
func Auth(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				unauthorized(w, "falta el token de autorización")
				return
			}

			claims, err := security.ParseAccessToken(strings.TrimPrefix(header, "Bearer "), secret)
			if err != nil {
				unauthorized(w, "token inválido o expirado")
				return
			}

			ctx := context.WithValue(r.Context(), userIDContextKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDContextKey).(string)
	return id, ok
}
