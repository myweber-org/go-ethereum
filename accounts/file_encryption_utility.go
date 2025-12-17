
package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "encoding/hex"
    "errors"
    "fmt"
    "io"
    "os"
)

const (
    saltSize      = 16
    nonceSize     = 12
    keyIterations = 100000
    keyLength     = 32
)

func deriveKey(password string, salt []byte) []byte {
    hash := sha256.New()
    hash.Write([]byte(password))
    hash.Write(salt)
    for i := 0; i < keyIterations; i++ {
        hash.Write(hash.Sum(nil))
    }
    return hash.Sum(nil)[:keyLength]
}

func encryptFile(inputPath, outputPath, password string) error {
    inputData, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("failed to read input file: %w", err)
    }

    salt := make([]byte, saltSize)
    if _, err := rand.Read(salt); err != nil {
        return fmt.Errorf("failed to generate salt: %w", err)
    }

    nonce := make([]byte, nonceSize)
    if _, err := rand.Read(nonce); err != nil {
        return fmt.Errorf("failed to generate nonce: %w", err)
    }

    key := deriveKey(password, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("failed to create cipher: %w", err)
    }

    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("failed to create GCM: %w", err)
    }

    ciphertext := aesgcm.Seal(nil, nonce, inputData, nil)

    outputData := make([]byte, 0, len(salt)+len(nonce)+len(ciphertext))
    outputData = append(outputData, salt...)
    outputData = append(outputData, nonce...)
    outputData = append(outputData, ciphertext...)

    if err := os.WriteFile(outputPath, outputData, 0644); err != nil {
        return fmt.Errorf("failed to write output file: %w", err)
    }

    return nil
}

func decryptFile(inputPath, outputPath, password string) error {
    inputData, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("failed to read input file: %w", err)
    }

    if len(inputData) < saltSize+nonceSize {
        return errors.New("input file too short")
    }

    salt := inputData[:saltSize]
    nonce := inputData[saltSize : saltSize+nonceSize]
    ciphertext := inputData[saltSize+nonceSize:]

    key := deriveKey(password, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("failed to create cipher: %w", err)
    }

    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("failed to create GCM: %w", err)
    }

    plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return fmt.Errorf("decryption failed: %w", err)
    }

    if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
        return fmt.Errorf("failed to write output file: %w", err)
    }

    return nil
}

func main() {
    if len(os.Args) < 5 {
        fmt.Println("Usage: go run file_encryption_utility.go <encrypt|decrypt> <input> <output> <password>")
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

    fmt.Printf("Operation completed successfully: %s -> %s\n", inputPath, outputPath)
}