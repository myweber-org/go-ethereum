
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
    "os"
    "strings"
)

const (
    saltSize        = 16
    nonceSize       = 12
    keyIterations   = 100000
    keyLength       = 32
)

func deriveKey(password string, salt []byte) []byte {
    hash := sha256.New()
    hash.Write([]byte(password))
    hash.Write(salt)
    for i := 0; i < keyIterations-1; i++ {
        hash.Write(hash.Sum(nil))
    }
    return hash.Sum(nil)[:keyLength]
}

func encryptData(plaintext, password string) (string, error) {
    salt := make([]byte, saltSize)
    if _, err := io.ReadFull(rand.Reader, salt); err != nil {
        return "", err
    }

    nonce := make([]byte, nonceSize)
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }

    key := deriveKey(password, salt)
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    ciphertext := aesgcm.Seal(nil, nonce, []byte(plaintext), nil)
    combined := append(salt, nonce...)
    combined = append(combined, ciphertext...)
    return base64.StdEncoding.EncodeToString(combined), nil
}

func decryptData(encrypted, password string) (string, error) {
    data, err := base64.StdEncoding.DecodeString(encrypted)
    if err != nil {
        return "", err
    }

    if len(data) < saltSize+nonceSize {
        return "", errors.New("invalid encrypted data")
    }

    salt := data[:saltSize]
    nonce := data[saltSize : saltSize+nonceSize]
    ciphertext := data[saltSize+nonceSize:]

    key := deriveKey(password, salt)
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", err
    }

    return string(plaintext), nil
}

func main() {
    if len(os.Args) < 4 {
        fmt.Println("Usage: go run file_encryption_utility.go <encrypt|decrypt> <input> <password>")
        os.Exit(1)
    }

    operation := strings.ToLower(os.Args[1])
    input := os.Args[2]
    password := os.Args[3]

    switch operation {
    case "encrypt":
        result, err := encryptData(input, password)
        if err != nil {
            fmt.Printf("Encryption failed: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Encrypted: %s\n", result)

    case "decrypt":
        result, err := decryptData(input, password)
        if err != nil {
            fmt.Printf("Decryption failed: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Decrypted: %s\n", result)

    default:
        fmt.Println("Invalid operation. Use 'encrypt' or 'decrypt'")
        os.Exit(1)
    }
}