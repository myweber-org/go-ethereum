package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "errors"
    "fmt"
    "io"
    "os"
)

func encryptFile(inputPath, outputPath string, key []byte) error {
    plaintext, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("failed to read input file: %w", err)
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("failed to create cipher: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("failed to create GCM: %w", err)
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return fmt.Errorf("failed to generate nonce: %w", err)
    }

    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

    if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
        return fmt.Errorf("failed to write encrypted file: %w", err)
    }

    return nil
}

func decryptFile(inputPath, outputPath string, key []byte) error {
    ciphertext, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("failed to read encrypted file: %w", err)
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("failed to create cipher: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("failed to create GCM: %w", err)
    }

    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return errors.New("ciphertext too short")
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return fmt.Errorf("failed to decrypt: %w", err)
    }

    if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
        return fmt.Errorf("failed to write decrypted file: %w", err)
    }

    return nil
}

func main() {
    key := make([]byte, 32)
    if _, err := rand.Read(key); err != nil {
        fmt.Printf("Failed to generate key: %v\n", err)
        return
    }

    testFile := "test_data.txt"
    encryptedFile := "test_data.enc"
    decryptedFile := "test_data_decrypted.txt"

    if err := os.WriteFile(testFile, []byte("Sensitive information here"), 0644); err != nil {
        fmt.Printf("Failed to create test file: %v\n", err)
        return
    }
    defer os.Remove(testFile)
    defer os.Remove(encryptedFile)
    defer os.Remove(decryptedFile)

    fmt.Println("Testing encryption...")
    if err := encryptFile(testFile, encryptedFile, key); err != nil {
        fmt.Printf("Encryption failed: %v\n", err)
        return
    }

    fmt.Println("Testing decryption...")
    if err := decryptFile(encryptedFile, decryptedFile, key); err != nil {
        fmt.Printf("Decryption failed: %v\n", err)
        return
    }

    original, _ := os.ReadFile(testFile)
    decrypted, _ := os.ReadFile(decryptedFile)

    if string(original) == string(decrypted) {
        fmt.Println("Success: Encryption and decryption verified")
    } else {
        fmt.Println("Error: Decrypted content does not match original")
    }
}