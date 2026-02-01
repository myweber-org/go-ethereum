
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
)

func deriveKey(passphrase string, salt []byte) []byte {
    hash := sha256.New()
    hash.Write([]byte(passphrase))
    hash.Write(salt)
    return hash.Sum(nil)
}

func encryptData(plaintext []byte, passphrase string) (string, error) {
    salt := make([]byte, 16)
    if _, err := rand.Read(salt); err != nil {
        return "", err
    }

    key := deriveKey(passphrase, salt)
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
    result := make([]byte, len(salt)+len(ciphertext))
    copy(result[:16], salt)
    copy(result[16:], ciphertext)

    return base64.StdEncoding.EncodeToString(result), nil
}

func decryptData(encrypted string, passphrase string) ([]byte, error) {
    data, err := base64.StdEncoding.DecodeString(encrypted)
    if err != nil {
        return nil, err
    }

    if len(data) < 16 {
        return nil, errors.New("invalid encrypted data")
    }

    salt := data[:16]
    ciphertext := data[16:]

    key := deriveKey(passphrase, salt)
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
    if len(os.Args) < 4 {
        fmt.Println("Usage: go run file_encryption.go <encrypt|decrypt> <input> <passphrase>")
        os.Exit(1)
    }

    mode := os.Args[1]
    input := os.Args[2]
    passphrase := os.Args[3]

    switch mode {
    case "encrypt":
        data, err := os.ReadFile(input)
        if err != nil {
            fmt.Printf("Error reading file: %v\n", err)
            os.Exit(1)
        }

        encrypted, err := encryptData(data, passphrase)
        if err != nil {
            fmt.Printf("Encryption failed: %v\n", err)
            os.Exit(1)
        }

        outputFile := input + ".enc"
        if err := os.WriteFile(outputFile, []byte(encrypted), 0644); err != nil {
            fmt.Printf("Error writing file: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Encrypted data saved to: %s\n", outputFile)

    case "decrypt":
        data, err := os.ReadFile(input)
        if err != nil {
            fmt.Printf("Error reading file: %v\n", err)
            os.Exit(1)
        }

        decrypted, err := decryptData(string(data), passphrase)
        if err != nil {
            fmt.Printf("Decryption failed: %v\n", err)
            os.Exit(1)
        }

        outputFile := input + ".dec"
        if err := os.WriteFile(outputFile, decrypted, 0644); err != nil {
            fmt.Printf("Error writing file: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Decrypted data saved to: %s\n", outputFile)

    default:
        fmt.Println("Invalid mode. Use 'encrypt' or 'decrypt'")
        os.Exit(1)
    }
}