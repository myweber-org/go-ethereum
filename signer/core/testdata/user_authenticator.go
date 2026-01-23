package middleware

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const userIDKey contextKey = "userID"

type Authenticator struct {
	tokenValidator func(string) (string, bool)
}

func NewAuthenticator(validator func(string) (string, bool)) *Authenticator {
	return &Authenticator{
		tokenValidator: validator,
	}
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		userID, valid := a.tokenValidator(token)
		if !valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}package auth

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("your-secret-key-change-in-production")

type Claims struct {
    Username string `json:"username"`
    UserID   int    `json:"user_id"`
    jwt.RegisteredClaims
}

func GenerateToken(username string, userID int) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    
    claims := &Claims{
        Username: username,
        UserID:   userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "myapp",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(secretKey)
}

func ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return secretKey, nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid token")
}

func ExtractUserID(tokenString string) (int, error) {
    claims, err := ValidateToken(tokenString)
    if err != nil {
        return 0, err
    }
    return claims.UserID, nil
}