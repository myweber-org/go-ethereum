package middleware

import (
	"net/http"
	"strings"
)

type Authenticator struct {
	secretKey []byte
}

func NewAuthenticator(secret string) *Authenticator {
	return &Authenticator{secretKey: []byte(secret)}
}

func (a *Authenticator) ValidateToken(tokenString string) (bool, error) {
	if tokenString == "" {
		return false, nil
	}
	
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return false, nil
	}
	
	return true, nil
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Bearer token required", http.StatusUnauthorized)
			return
		}

		valid, err := a.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Token validation error", http.StatusInternalServerError)
			return
		}
		if !valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}