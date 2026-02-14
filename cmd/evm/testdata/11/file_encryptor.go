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

func generateKey() ([]byte, error) {
    key := make([]byte, 32)
    _, err := rand.Read(key)
    if err != nil {
        return nil, err
    }
    return key, nil
}

func encryptData(plaintext []byte, key []byte) (string, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }

    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decryptData(encrypted string, key []byte) ([]byte, error) {
    ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
    if err != nil {
        return nil, err
    }

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
    return gcm.Open(nil, nonce, ciphertext, nil)
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run file_encryptor.go <text_to_encrypt>")
        return
    }

    plaintext := []byte(os.Args[1])
    key, err := generateKey()
    if err != nil {
        fmt.Printf("Key generation error: %v\n", err)
        return
    }

    encrypted, err := encryptData(plaintext, key)
    if err != nil {
        fmt.Printf("Encryption error: %v\n", err)
        return
    }

    fmt.Printf("Generated Key (base64): %s\n", base64.StdEncoding.EncodeToString(key))
    fmt.Printf("Encrypted Data: %s\n", encrypted)

    decrypted, err := decryptData(encrypted, key)
    if err != nil {
        fmt.Printf("Decryption error: %v\n", err)
        return
    }

    fmt.Printf("Decrypted Text: %s\n", string(decrypted))
}