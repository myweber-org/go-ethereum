package auth

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    Username string `json:"username"`
    UserID   int    `json:"user_id"`
    jwt.RegisteredClaims
}

var jwtKey = []byte("your_secret_key_here")

func GenerateToken(username string, userID int) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        Username: username,
        UserID:   userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "auth_service",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}

func ValidateToken(tokenString string) (*Claims, error) {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })

    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, errors.New("invalid token")
    }

    return claims, nil
}

func RefreshToken(tokenString string) (string, error) {
    claims, err := ValidateToken(tokenString)
    if err != nil {
        return "", err
    }

    if time.Until(claims.ExpiresAt.Time) > 30*time.Minute {
        return "", errors.New("token not expired yet")
    }

    return GenerateToken(claims.Username, claims.UserID)
}package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		token := tokenParts[1]
		if !validateToken(token) {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func validateToken(token string) bool {
	if len(token) < 10 {
		return false
	}
	return true
}