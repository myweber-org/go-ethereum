
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
)

type FileEncryptor struct {
	key []byte
}

func NewFileEncryptor(key []byte) (*FileEncryptor, error) {
	if len(key) != 32 {
		return nil, errors.New("encryption key must be 32 bytes for AES-256")
	}
	return &FileEncryptor{key: key}, nil
}

func (fe *FileEncryptor) EncryptFile(inputPath, outputPath string) error {
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	block, err := aes.NewCipher(fe.key)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
		return fmt.Errorf("failed to write encrypted file: %w", err)
	}

	return nil
}

func (fe *FileEncryptor) DecryptFile(inputPath, outputPath string) error {
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read encrypted file: %w", err)
	}

	block, err := aes.NewCipher(fe.key)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("failed to decrypt: %w", err)
	}

	if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
		return fmt.Errorf("failed to write decrypted file: %w", err)
	}

	return nil
}

func generateRandomKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("failed to generate random key: %w", err)
	}
	return key, nil
}

func main() {
	key, err := generateRandomKey()
	if err != nil {
		fmt.Printf("Error generating key: %v\n", err)
		os.Exit(1)
	}

	encryptor, err := NewFileEncryptor(key)
	if err != nil {
		fmt.Printf("Error creating encryptor: %v\n", err)
		os.Exit(1)
	}

	testData := []byte("This is a secret message for encryption testing.")
	testFile := "test_data.txt"
	encryptedFile := "test_data.enc"
	decryptedFile := "test_data_decrypted.txt"

	if err := os.WriteFile(testFile, testData, 0644); err != nil {
		fmt.Printf("Error creating test file: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(testFile)
	defer os.Remove(encryptedFile)
	defer os.Remove(decryptedFile)

	fmt.Println("Testing file encryption...")

	if err := encryptor.EncryptFile(testFile, encryptedFile); err != nil {
		fmt.Printf("Encryption failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("File encrypted successfully")

	if err := encryptor.DecryptFile(encryptedFile, decryptedFile); err != nil {
		fmt.Printf("Decryption failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("File decrypted successfully")

	decryptedData, err := os.ReadFile(decryptedFile)
	if err != nil {
		fmt.Printf("Error reading decrypted file: %v\n", err)
		os.Exit(1)
	}

	if string(decryptedData) == string(testData) {
		fmt.Println("Encryption/decryption test passed")
	} else {
		fmt.Println("Encryption/decryption test failed")
		os.Exit(1)
	}
}