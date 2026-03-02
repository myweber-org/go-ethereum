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

func deriveKey(passphrase string, salt []byte) []byte {
    hash := sha256.New()
    hash.Write([]byte(passphrase))
    hash.Write(salt)
    return hash.Sum(nil)
}

func encryptData(plaintext []byte, passphrase string) ([]byte, error) {
    salt := make([]byte, 16)
    if _, err := rand.Read(salt); err != nil {
        return nil, err
    }

    key := deriveKey(passphrase, salt)
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
    return append(salt, ciphertext...), nil
}

func decryptData(ciphertext []byte, passphrase string) ([]byte, error) {
    if len(ciphertext) < 16 {
        return nil, errors.New("ciphertext too short")
    }

    salt := ciphertext[:16]
    ciphertext = ciphertext[16:]

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
        fmt.Println("Usage: go run file_encryption_utility.go <encrypt|decrypt> <input_file> <output_file>")
        fmt.Println("Set environment variable ENCRYPTION_KEY for passphrase")
        return
    }

    operation := os.Args[1]
    inputFile := os.Args[2]
    outputFile := os.Args[3]
    passphrase := os.Getenv("ENCRYPTION_KEY")

    if passphrase == "" {
        fmt.Println("Error: ENCRYPTION_KEY environment variable not set")
        return
    }

    inputData, err := os.ReadFile(inputFile)
    if err != nil {
        fmt.Printf("Error reading input file: %v\n", err)
        return
    }

    var outputData []byte
    switch operation {
    case "encrypt":
        outputData, err = encryptData(inputData, passphrase)
        if err != nil {
            fmt.Printf("Encryption error: %v\n", err)
            return
        }
        fmt.Printf("Encrypted %d bytes to %s\n", len(outputData), outputFile)
    case "decrypt":
        outputData, err = decryptData(inputData, passphrase)
        if err != nil {
            fmt.Printf("Decryption error: %v\n", err)
            return
        }
        fmt.Printf("Decrypted %d bytes to %s\n", len(outputData), outputFile)
    default:
        fmt.Println("Invalid operation. Use 'encrypt' or 'decrypt'")
        return
    }

    if err := os.WriteFile(outputFile, outputData, 0644); err != nil {
        fmt.Printf("Error writing output file: %v\n", err)
        return
    }

    if operation == "encrypt" {
        fmt.Println("Encrypted hex:", hex.EncodeToString(outputData[:32]))
    }
}