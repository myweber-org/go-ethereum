package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "errors"
    "fmt"
    "io"
    "os"
)

func encryptData(plaintext []byte, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }

    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
    return ciphertext, nil
}

func decryptData(ciphertext []byte, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return nil, errors.New("ciphertext too short")
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, err
    }

    return plaintext, nil
}

func generateKey() ([]byte, error) {
    key := make([]byte, 32)
    if _, err := rand.Read(key); err != nil {
        return nil, err
    }
    return key, nil
}

func main() {
    key, err := generateKey()
    if err != nil {
        fmt.Println("Error generating key:", err)
        return
    }

    fmt.Println("Generated key (base64):", base64.StdEncoding.EncodeToString(key))

    originalText := "Sensitive data requiring encryption"
    fmt.Println("Original text:", originalText)

    encrypted, err := encryptData([]byte(originalText), key)
    if err != nil {
        fmt.Println("Encryption error:", err)
        return
    }

    fmt.Println("Encrypted (base64):", base64.StdEncoding.EncodeToString(encrypted))

    decrypted, err := decryptData(encrypted, key)
    if err != nil {
        fmt.Println("Decryption error:", err)
        return
    }

    fmt.Println("Decrypted text:", string(decrypted))

    if string(decrypted) == originalText {
        fmt.Println("Encryption/decryption successful")
    } else {
        fmt.Println("Encryption/decryption failed")
        os.Exit(1)
    }
}