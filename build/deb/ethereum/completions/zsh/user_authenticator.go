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

func (a *Authenticator) ValidateToken(tokenString string) bool {
	if tokenString == "" {
		return false
	}
	
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return false
	}
	
	return verifySignature(parts, a.secretKey)
}

func verifySignature(parts []string, secret []byte) bool {
	// Simplified signature verification logic
	// In production, use proper JWT library
	expectedSig := generateSignature(parts[0]+"."+parts[1], secret)
	return parts[2] == expectedSig
}

func generateSignature(data string, secret []byte) string {
	// Simplified signature generation
	// In production, use proper HMAC or RSA signing
	return "simulated_signature_for_" + data
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if !a.ValidateToken(token) {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}