package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "errors"
    "fmt"
    "io"
    "strings"

    "golang.org/x/crypto/pbkdf2"
)

const (
    saltSize   = 16
    nonceSize  = 12
    keySize    = 32
    iterations = 100000
)

type EncryptionResult struct {
    Ciphertext string
    Salt       string
    Nonce      string
}

func deriveKey(password string, salt []byte) []byte {
    return pbkdf2.Key([]byte(password), salt, iterations, keySize, sha256.New)
}

func Encrypt(plaintext, password string) (*EncryptionResult, error) {
    salt := make([]byte, saltSize)
    if _, err := io.ReadFull(rand.Reader, salt); err != nil {
        return nil, fmt.Errorf("salt generation failed: %w", err)
    }

    key := deriveKey(password, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, fmt.Errorf("cipher creation failed: %w", err)
    }

    nonce := make([]byte, nonceSize)
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, fmt.Errorf("nonce generation failed: %w", err)
    }

    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("GCM mode initialization failed: %w", err)
    }

    ciphertext := aesgcm.Seal(nil, nonce, []byte(plaintext), nil)

    return &EncryptionResult{
        Ciphertext: base64.StdEncoding.EncodeToString(ciphertext),
        Salt:       base64.StdEncoding.EncodeToString(salt),
        Nonce:      base64.StdEncoding.EncodeToString(nonce),
    }, nil
}

func Decrypt(result *EncryptionResult, password string) (string, error) {
    salt, err := base64.StdEncoding.DecodeString(result.Salt)
    if err != nil {
        return "", fmt.Errorf("salt decoding failed: %w", err)
    }

    nonce, err := base64.StdEncoding.DecodeString(result.Nonce)
    if err != nil {
        return "", fmt.Errorf("nonce decoding failed: %w", err)
    }

    ciphertext, err := base64.StdEncoding.DecodeString(result.Ciphertext)
    if err != nil {
        return "", fmt.Errorf("ciphertext decoding failed: %w", err)
    }

    key := deriveKey(password, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return "", fmt.Errorf("cipher creation failed: %w", err)
    }

    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", fmt.Errorf("GCM mode initialization failed: %w", err)
    }

    plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", errors.New("decryption failed - incorrect password or corrupted data")
    }

    return string(plaintext), nil
}

func main() {
    secretMessage := "Confidential data: API keys, tokens, and sensitive configuration"
    password := "StrongPassw0rd!2024"

    fmt.Println("Original:", secretMessage)

    encrypted, err := Encrypt(secretMessage, password)
    if err != nil {
        fmt.Printf("Encryption error: %v\n", err)
        return
    }

    fmt.Printf("Encrypted: %s\n", strings.Join([]string{
        encrypted.Ciphertext[:30] + "...",
        encrypted.Salt,
        encrypted.Nonce,
    }, " | "))

    decrypted, err := Decrypt(encrypted, password)
    if err != nil {
        fmt.Printf("Decryption error: %v\n", err)
        return
    }

    fmt.Println("Decrypted:", decrypted)

    wrongPassword := "WrongPassword123"
    _, err = Decrypt(encrypted, wrongPassword)
    if err != nil {
        fmt.Println("Expected error with wrong password:", err)
    }
}