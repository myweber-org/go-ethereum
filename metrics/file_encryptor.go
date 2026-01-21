
package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "io"
    "os"
    "path/filepath"
)

const (
    saltSize   = 16
    keySize    = 32
    nonceSize  = 12
)

func deriveKey(password, salt []byte) []byte {
    hash := sha256.New()
    hash.Write(password)
    hash.Write(salt)
    return hash.Sum(nil)
}

func encryptFile(inputPath, outputPath, password string) error {
    plaintext, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("read file failed: %w", err)
    }

    salt := make([]byte, saltSize)
    if _, err := rand.Read(salt); err != nil {
        return fmt.Errorf("generate salt failed: %w", err)
    }

    key := deriveKey([]byte(password), salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("create cipher failed: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("create GCM failed: %w", err)
    }

    nonce := make([]byte, nonceSize)
    if _, err := rand.Read(nonce); err != nil {
        return fmt.Errorf("generate nonce failed: %w", err)
    }

    ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

    outputData := append(salt, nonce...)
    outputData = append(outputData, ciphertext...)

    if err := os.WriteFile(outputPath, outputData, 0644); err != nil {
        return fmt.Errorf("write file failed: %w", err)
    }

    return nil
}

func decryptFile(inputPath, outputPath, password string) error {
    data, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("read file failed: %w", err)
    }

    if len(data) < saltSize+nonceSize {
        return fmt.Errorf("invalid encrypted file format")
    }

    salt := data[:saltSize]
    nonce := data[saltSize:saltSize+nonceSize]
    ciphertext := data[saltSize+nonceSize:]

    key := deriveKey([]byte(password), salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("create cipher failed: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("create GCM failed: %w", err)
    }

    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return fmt.Errorf("decryption failed: %w", err)
    }

    if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
        return fmt.Errorf("write file failed: %w", err)
    }

    return nil
}

func main() {
    if len(os.Args) < 5 {
        fmt.Println("Usage: file_encryptor <encrypt|decrypt> <input> <output> <password>")
        os.Exit(1)
    }

    mode := os.Args[1]
    inputPath := os.Args[2]
    outputPath := os.Args[3]
    password := os.Args[4]

    var err error
    switch mode {
    case "encrypt":
        err = encryptFile(inputPath, outputPath, password)
    case "decrypt":
        err = decryptFile(inputPath, outputPath, password)
    default:
        fmt.Printf("Invalid mode: %s\n", mode)
        os.Exit(1)
    }

    if err != nil {
        fmt.Printf("Operation failed: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Operation completed successfully\n")
}package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

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
	data, err := base64.StdEncoding.DecodeString(encrypted)
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
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}