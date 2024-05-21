package middleware

import (
	"atur-dana/internal/auth"
	"atur-dana/internal/common"
	"net/http"
	"strings"

	"github.com/gorilla/context"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check headers
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			common.JSONError(w, http.StatusUnauthorized, "Missing Authorization header")
			return
		}

		// validate JWT
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		metadata, err := auth.ValidateJWT(tokenString)
		if err != nil {
			common.JSONError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		context.Set(r, "metadata", metadata)
		next.ServeHTTP(w, r)
	})
}
