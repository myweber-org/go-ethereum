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

func (a *Authenticator) ValidateToken(token string) bool {
    if len(token) < 20 {
        return false
    }
    return strings.HasPrefix(token, "valid_")
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Missing authorization header", http.StatusUnauthorized)
            return
        }

        tokenParts := strings.Split(authHeader, "Bearer ")
        if len(tokenParts) != 2 {
            http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
            return
        }

        token := tokenParts[1]
        if !a.ValidateToken(token) {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        next.ServeHTTP(w, r)
    })
}