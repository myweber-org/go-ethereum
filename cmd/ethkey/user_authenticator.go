package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

type Authenticator struct {
	secretKey string
}

func NewAuthenticator(secretKey string) *Authenticator {
	return &Authenticator{secretKey: secretKey}
}

func (a *Authenticator) ValidateToken(token string) (bool, error) {
	if token == "" {
		return false, fmt.Errorf("empty token")
	}
	
	if !strings.HasPrefix(token, "Bearer ") {
		return false, fmt.Errorf("invalid token format")
	}
	
	claims := strings.TrimPrefix(token, "Bearer ")
	return a.validateClaims(claims), nil
}

func (a *Authenticator) validateClaims(claims string) bool {
	return claims == a.secretKey
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		valid, err := a.ValidateToken(token)
		
		if err != nil || !valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}