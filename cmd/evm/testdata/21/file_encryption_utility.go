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
        return fmt.Errorf("failed to write output file: %w", err)
    }

    return nil
}

func decryptFile(inputPath, outputPath string, key []byte) error {
    ciphertext, err := os.ReadFile(inputPath)
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
        return fmt.Errorf("failed to write output file: %w", err)
    }

    return nil
}

func main() {
    key := []byte("examplekey123456") // 16 bytes for AES-128
    inputFile := "test.txt"
    encryptedFile := "test.enc"
    decryptedFile := "test_decrypted.txt"

    // Create a test file
    if err := os.WriteFile(inputFile, []byte("Secret data for encryption test"), 0644); err != nil {
        fmt.Printf("Failed to create test file: %v\n", err)
        return
    }
    defer os.Remove(inputFile)
    defer os.Remove(encryptedFile)
    defer os.Remove(decryptedFile)

    fmt.Println("Encrypting file...")
    if err := encryptFile(inputFile, encryptedFile, key); err != nil {
        fmt.Printf("Encryption failed: %v\n", err)
        return
    }

    fmt.Println("Decrypting file...")
    if err := decryptFile(encryptedFile, decryptedFile, key); err != nil {
        fmt.Printf("Decryption failed: %v\n", err)
        return
    }

    content, err := os.ReadFile(decryptedFile)
    if err != nil {
        fmt.Printf("Failed to read decrypted file: %v\n", err)
        return
    }

    fmt.Printf("Decrypted content: %s\n", string(content))
    fmt.Println("Operation completed successfully")
}