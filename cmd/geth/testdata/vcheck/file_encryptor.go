
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
    "golang.org/x/crypto/pbkdf2"
)

const (
    saltSize   = 16
    nonceSize  = 12
    keySize    = 32
    iterations = 100000
)

type EncryptedData struct {
    Ciphertext string `json:"ciphertext"`
    Salt       string `json:"salt"`
    Nonce      string `json:"nonce"`
}

func deriveKey(password string, salt []byte) []byte {
    return pbkdf2.Key([]byte(password), salt, iterations, keySize, sha256.New)
}

func Encrypt(plaintext, password string) (*EncryptedData, error) {
    salt := make([]byte, saltSize)
    if _, err := io.ReadFull(rand.Reader, salt); err != nil {
        return nil, fmt.Errorf("failed to generate salt: %w", err)
    }

    key := deriveKey(password, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, fmt.Errorf("failed to create cipher: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("failed to create GCM: %w", err)
    }

    nonce := make([]byte, nonceSize)
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, fmt.Errorf("failed to generate nonce: %w", err)
    }

    ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)

    return &EncryptedData{
        Ciphertext: base64.StdEncoding.EncodeToString(ciphertext),
        Salt:       base64.StdEncoding.EncodeToString(salt),
        Nonce:      base64.StdEncoding.EncodeToString(nonce),
    }, nil
}

func Decrypt(data *EncryptedData, password string) (string, error) {
    salt, err := base64.StdEncoding.DecodeString(data.Salt)
    if err != nil {
        return "", fmt.Errorf("invalid salt encoding: %w", err)
    }

    nonce, err := base64.StdEncoding.DecodeString(data.Nonce)
    if err != nil {
        return "", fmt.Errorf("invalid nonce encoding: %w", err)
    }

    ciphertext, err := base64.StdEncoding.DecodeString(data.Ciphertext)
    if err != nil {
        return "", fmt.Errorf("invalid ciphertext encoding: %w", err)
    }

    key := deriveKey(password, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return "", fmt.Errorf("failed to create cipher: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", fmt.Errorf("failed to create GCM: %w", err)
    }

    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", errors.New("decryption failed: invalid password or corrupted data")
    }

    return string(plaintext), nil
}

func main() {
    secret := "confidential data"
    password := "securePass123!"

    encrypted, err := Encrypt(secret, password)
    if err != nil {
        fmt.Printf("Encryption error: %v\n", err)
        return
    }

    fmt.Printf("Encrypted: %+v\n", encrypted)

    decrypted, err := Decrypt(encrypted, password)
    if err != nil {
        fmt.Printf("Decryption error: %v\n", err)
        return
    }

    fmt.Printf("Decrypted: %s\n", decrypted)

    if decrypted == secret {
        fmt.Println("Encryption/decryption successful")
    } else {
        fmt.Println("Encryption/decryption failed")
    }
}