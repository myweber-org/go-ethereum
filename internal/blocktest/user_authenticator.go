package main

import (
    "fmt"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

type User struct {
    ID       int
    Username string
    Email    string
}

type Claims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    jwt.RegisteredClaims
}

var jwtSecret = []byte("your-secret-key-change-in-production")

func GenerateToken(user User) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)

    claims := &Claims{
        UserID:   user.ID,
        Username: user.Username,
        Email:    user.Email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "auth-service",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (*Claims, error) {
    claims := &Claims{}

    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })

    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }

    return claims, nil
}

func main() {
    testUser := User{
        ID:       1,
        Username: "john_doe",
        Email:    "john@example.com",
    }

    token, err := GenerateToken(testUser)
    if err != nil {
        fmt.Printf("Error generating token: %v\n", err)
        return
    }

    fmt.Printf("Generated Token: %s\n\n", token)

    claims, err := ValidateToken(token)
    if err != nil {
        fmt.Printf("Error validating token: %v\n", err)
        return
    }

    fmt.Printf("Token validated successfully!\n")
    fmt.Printf("User ID: %d\n", claims.UserID)
    fmt.Printf("Username: %s\n", claims.Username)
    fmt.Printf("Email: %s\n", claims.Email)
}