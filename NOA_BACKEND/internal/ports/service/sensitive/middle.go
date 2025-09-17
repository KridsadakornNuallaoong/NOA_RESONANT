package sensitive

import (
	"context"
	"net/http"
	"strings"
)

// Define a custom type for the context key

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Extract the token string
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Verify the JWT token
		claims, err := VerifyJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add claims to the context using the custom key
		ctx := context.WithValue(r.Context(), userContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
