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
	"path/filepath"
)

const (
	saltSize      = 16
	keyIterations = 100000
	keyLength     = 32
)

func deriveKey(password, salt []byte) []byte {
	hash := sha256.New()
	hash.Write(password)
	hash.Write(salt)
	for i := 0; i < keyIterations-1; i++ {
		hash.Write(hash.Sum(nil))
	}
	return hash.Sum(nil)[:keyLength]
}

func encryptFile(inputPath, outputPath, password string) error {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return fmt.Errorf("salt generation failed: %w", err)
	}

	key := deriveKey([]byte(password), salt)

	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read file failed: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation failed: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("GCM mode failed: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return fmt.Errorf("nonce generation failed: %w", err)
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
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read file failed: %w", err)
	}

	if len(ciphertext) < saltSize {
		return errors.New("file too short to contain salt")
	}

	salt := ciphertext[:saltSize]
	key := deriveKey([]byte(password), salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation failed: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("GCM mode failed: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < saltSize+nonceSize {
		return errors.New("file too short to contain nonce")
	}

	nonce := ciphertext[saltSize : saltSize+nonceSize]
	actualCiphertext := ciphertext[saltSize+nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return errors.New("decryption failed - wrong password or corrupted file")
	}

	if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
		return fmt.Errorf("write file failed: %w", err)
	}

	return nil
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: file_encryptor <encrypt|decrypt> <input> <output>")
		fmt.Println("Password will be prompted securely")
		os.Exit(1)
	}

	operation := os.Args[1]
	inputPath := os.Args[2]
	outputPath := os.Args[3]

	fmt.Print("Enter password: ")
	var password string
	fmt.Scanln(&password)

	var err error
	switch operation {
	case "encrypt":
		err = encryptFile(inputPath, outputPath, password)
	case "decrypt":
		err = decryptFile(inputPath, outputPath, password)
	default:
		fmt.Println("Invalid operation. Use 'encrypt' or 'decrypt'")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Operation completed successfully: %s -> %s\n", inputPath, outputPath)
}