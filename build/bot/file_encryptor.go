
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
    "strings"
)

func deriveKey(password string, salt []byte) []byte {
    hash := sha256.New()
    hash.Write([]byte(password))
    hash.Write(salt)
    return hash.Sum(nil)
}

func encryptData(plaintext []byte, password string) (string, error) {
    salt := make([]byte, 16)
    if _, err := rand.Read(salt); err != nil {
        return "", err
    }

    key := deriveKey(password, salt)
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := rand.Read(nonce); err != nil {
        return "", err
    }

    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
    combined := append(salt, ciphertext...)
    return base64.StdEncoding.EncodeToString(combined), nil
}

func decryptData(encodedCiphertext string, password string) ([]byte, error) {
    combined, err := base64.StdEncoding.DecodeString(encodedCiphertext)
    if err != nil {
        return nil, err
    }

    if len(combined) < 16 {
        return nil, errors.New("invalid ciphertext length")
    }

    salt := combined[:16]
    ciphertext := combined[16:]

    key := deriveKey(password, salt)
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
        return nil, errors.New("invalid ciphertext")
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    return gcm.Open(nil, nonce, ciphertext, nil)
}

func main() {
    if len(os.Args) < 4 {
        fmt.Println("Usage: go run file_encryptor.go <encrypt|decrypt> <input_file> <password> [output_file]")
        os.Exit(1)
    }

    mode := strings.ToLower(os.Args[1])
    inputFile := os.Args[2]
    password := os.Args[3]
    outputFile := ""

    if len(os.Args) > 4 {
        outputFile = os.Args[4]
    } else {
        if mode == "encrypt" {
            outputFile = inputFile + ".enc"
        } else {
            outputFile = strings.TrimSuffix(inputFile, ".enc")
        }
    }

    data, err := os.ReadFile(inputFile)
    if err != nil {
        fmt.Printf("Error reading file: %v\n", err)
        os.Exit(1)
    }

    switch mode {
    case "encrypt":
        encrypted, err := encryptData(data, password)
        if err != nil {
            fmt.Printf("Encryption error: %v\n", err)
            os.Exit(1)
        }
        if err := os.WriteFile(outputFile, []byte(encrypted), 0644); err != nil {
            fmt.Printf("Error writing file: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("File encrypted successfully: %s\n", outputFile)

    case "decrypt":
        decrypted, err := decryptData(string(data), password)
        if err != nil {
            fmt.Printf("Decryption error: %v\n", err)
            os.Exit(1)
        }
        if err := os.WriteFile(outputFile, decrypted, 0644); err != nil {
            fmt.Printf("Error writing file: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("File decrypted successfully: %s\n", outputFile)

    default:
        fmt.Println("Invalid mode. Use 'encrypt' or 'decrypt'")
        os.Exit(1)
    }
}