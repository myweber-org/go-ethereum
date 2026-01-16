package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "errors"
    "fmt"
    "io"
    "os"
    "path/filepath"
)

func deriveKey(password string, salt []byte) []byte {
    hash := sha256.New()
    hash.Write([]byte(password))
    hash.Write(salt)
    return hash.Sum(nil)
}

func encryptFile(inputPath, outputPath, password string) error {
    plaintext, err := os.ReadFile(inputPath)
    if err != nil {
        return err
    }

    salt := make([]byte, 16)
    if _, err := rand.Read(salt); err != nil {
        return err
    }

    key := deriveKey(password, salt)
    block, err := aes.NewCipher(key)
    if err != nil {
        return err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := rand.Read(nonce); err != nil {
        return err
    }

    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
    ciphertext = append(salt, ciphertext...)

    return os.WriteFile(outputPath, ciphertext, 0644)
}

func decryptFile(inputPath, outputPath, password string) error {
    ciphertext, err := os.ReadFile(inputPath)
    if err != nil {
        return err
    }

    if len(ciphertext) < 16 {
        return errors.New("file too short")
    }

    salt := ciphertext[:16]
    ciphertext = ciphertext[16:]

    key := deriveKey(password, salt)
    block, err := aes.NewCipher(key)
    if err != nil {
        return err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return err
    }

    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return errors.New("ciphertext too short")
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return err
    }

    return os.WriteFile(outputPath, plaintext, 0644)
}

func main() {
    if len(os.Args) < 5 {
        fmt.Println("Usage: go run file_encryption.go <encrypt|decrypt> <input> <output> <password>")
        os.Exit(1)
    }

    mode := os.Args[1]
    inputPath := os.Args[2]
    outputPath := os.Args[3]
    password := os.Args[4]

    var err error
    switch mode {
    case "encrypt":
        err = encryptFile(inputPath, outputPath, password)
    case "decrypt":
        err = decryptFile(inputPath, outputPath, password)
    default:
        fmt.Println("Invalid mode. Use 'encrypt' or 'decrypt'")
        os.Exit(1)
    }

    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Operation completed successfully: %s -> %s\n", filepath.Base(inputPath), filepath.Base(outputPath))
}