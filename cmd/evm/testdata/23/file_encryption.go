
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

const saltSize = 16

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

    salt := make([]byte, saltSize)
    if _, err := io.ReadFull(rand.Reader, salt); err != nil {
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
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return err
    }

    ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
    finalData := append(salt, nonce...)
    finalData = append(finalData, ciphertext...)

    return os.WriteFile(outputPath, finalData, 0644)
}

func decryptFile(inputPath, outputPath, password string) error {
    data, err := os.ReadFile(inputPath)
    if err != nil {
        return err
    }

    if len(data) < saltSize {
        return errors.New("file too short")
    }

    salt := data[:saltSize]
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
    if len(data) < saltSize+nonceSize {
        return errors.New("file too short")
    }

    nonce := data[saltSize : saltSize+nonceSize]
    ciphertext := data[saltSize+nonceSize:]

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
    input := os.Args[2]
    output := os.Args[3]
    password := os.Args[4]

    switch mode {
    case "encrypt":
        err := encryptFile(input, output, password)
        if err != nil {
            fmt.Printf("Encryption failed: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("File encrypted successfully. Salt: %s\n", hex.EncodeToString([]byte{}))
    case "decrypt":
        err := decryptFile(input, output, password)
        if err != nil {
            fmt.Printf("Decryption failed: %v\n", err)
            os.Exit(1)
        }
        fmt.Println("File decrypted successfully.")
    default:
        fmt.Println("Invalid mode. Use 'encrypt' or 'decrypt'.")
        os.Exit(1)
    }
}