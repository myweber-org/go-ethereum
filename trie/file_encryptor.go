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

func encryptFile(inputPath, outputPath, keyHex string) error {
    key, err := hex.DecodeString(keyHex)
    if err != nil {
        return fmt.Errorf("invalid key: %v", err)
    }
    if len(key) != 32 {
        return errors.New("key must be 32 bytes for AES-256")
    }

    plaintext, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("read file failed: %v", err)
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("cipher creation failed: %v", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("GCM creation failed: %v", err)
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return fmt.Errorf("nonce generation failed: %v", err)
    }

    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

    if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
        return fmt.Errorf("write file failed: %v", err)
    }

    return nil
}

func decryptFile(inputPath, outputPath, keyHex string) error {
    key, err := hex.DecodeString(keyHex)
    if err != nil {
        return fmt.Errorf("invalid key: %v", err)
    }
    if len(key) != 32 {
        return errors.New("key must be 32 bytes for AES-256")
    }

    ciphertext, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("read file failed: %v", err)
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("cipher creation failed: %v", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("GCM creation failed: %v", err)
    }

    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return errors.New("ciphertext too short")
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return fmt.Errorf("decryption failed: %v", err)
    }

    if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
        return fmt.Errorf("write file failed: %v", err)
    }

    return nil
}

func generateRandomKey() (string, error) {
    key := make([]byte, 32)
    if _, err := rand.Read(key); err != nil {
        return "", fmt.Errorf("key generation failed: %v", err)
    }
    return hex.EncodeToString(key), nil
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage:")
        fmt.Println("  Generate key: file_encryptor -genkey")
        fmt.Println("  Encrypt: file_encryptor -encrypt input.txt output.enc KEY_HEX")
        fmt.Println("  Decrypt: file_encryptor -decrypt input.enc output.txt KEY_HEX")
        os.Exit(1)
    }

    switch os.Args[1] {
    case "-genkey":
        key, err := generateRandomKey()
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Generated key: %s\n", key)

    case "-encrypt":
        if len(os.Args) != 5 {
            fmt.Println("Usage: file_encryptor -encrypt input output key")
            os.Exit(1)
        }
        if err := encryptFile(os.Args[2], os.Args[3], os.Args[4]); err != nil {
            fmt.Printf("Encryption failed: %v\n", err)
            os.Exit(1)
        }
        fmt.Println("File encrypted successfully")

    case "-decrypt":
        if len(os.Args) != 5 {
            fmt.Println("Usage: file_encryptor -decrypt input output key")
            os.Exit(1)
        }
        if err := decryptFile(os.Args[2], os.Args[3], os.Args[4]); err != nil {
            fmt.Printf("Decryption failed: %v\n", err)
            os.Exit(1)
        }
        fmt.Println("File decrypted successfully")

    default:
        fmt.Println("Invalid command")
        os.Exit(1)
    }
}