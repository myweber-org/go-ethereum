package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/hex"
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

func generateRandomKey() ([]byte, error) {
    key := make([]byte, 32)
    if _, err := rand.Read(key); err != nil {
        return nil, err
    }
    return key, nil
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run file_encryption_tool.go <encrypt|decrypt|keygen>")
        return
    }

    switch os.Args[1] {
    case "keygen":
        key, err := generateRandomKey()
        if err != nil {
            fmt.Printf("Error generating key: %v\n", err)
            return
        }
        fmt.Printf("Generated key: %s\n", hex.EncodeToString(key))

    case "encrypt":
        if len(os.Args) != 4 {
            fmt.Println("Usage: go run file_encryption_tool.go encrypt <input_file> <key_hex>")
            return
        }

        key, err := hex.DecodeString(os.Args[3])
        if err != nil {
            fmt.Printf("Invalid key: %v\n", err)
            return
        }

        data, err := os.ReadFile(os.Args[2])
        if err != nil {
            fmt.Printf("Error reading file: %v\n", err)
            return
        }

        encrypted, err := encryptData(data, key)
        if err != nil {
            fmt.Printf("Encryption error: %v\n", err)
            return
        }

        outputFile := os.Args[2] + ".enc"
        if err := os.WriteFile(outputFile, encrypted, 0644); err != nil {
            fmt.Printf("Error writing file: %v\n", err)
            return
        }
        fmt.Printf("Encrypted file saved as: %s\n", outputFile)

    case "decrypt":
        if len(os.Args) != 4 {
            fmt.Println("Usage: go run file_encryption_tool.go decrypt <encrypted_file> <key_hex>")
            return
        }

        key, err := hex.DecodeString(os.Args[3])
        if err != nil {
            fmt.Printf("Invalid key: %v\n", err)
            return
        }

        data, err := os.ReadFile(os.Args[2])
        if err != nil {
            fmt.Printf("Error reading file: %v\n", err)
            return
        }

        decrypted, err := decryptData(data, key)
        if err != nil {
            fmt.Printf("Decryption error: %v\n", err)
            return
        }

        outputFile := os.Args[2] + ".dec"
        if err := os.WriteFile(outputFile, decrypted, 0644); err != nil {
            fmt.Printf("Error writing file: %v\n", err)
            return
        }
        fmt.Printf("Decrypted file saved as: %s\n", outputFile)

    default:
        fmt.Println("Invalid command. Use: encrypt, decrypt, or keygen")
    }
}