
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

type FileEncryptor struct {
	key []byte
}

func NewFileEncryptor(key string) (*FileEncryptor, error) {
	if len(key) != 64 {
		return nil, errors.New("key must be 64 hex characters for AES-256")
	}
	decodedKey, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}
	return &FileEncryptor{key: decodedKey}, nil
}

func (fe *FileEncryptor) EncryptFile(inputPath, outputPath string) error {
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(fe.key)
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

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return os.WriteFile(outputPath, ciphertext, 0644)
}

func (fe *FileEncryptor) DecryptFile(inputPath, outputPath string) error {
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(fe.key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, plaintext, 0644)
}

func generateRandomKey() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}

func main() {
	key, err := generateRandomKey()
	if err != nil {
		fmt.Printf("Failed to generate key: %v\n", err)
		return
	}

	encryptor, err := NewFileEncryptor(key)
	if err != nil {
		fmt.Printf("Failed to create encryptor: %v\n", err)
		return
	}

	testData := []byte("This is a secret message for encryption testing.")
	inputFile := "test_input.txt"
	encryptedFile := "test_encrypted.bin"
	decryptedFile := "test_decrypted.txt"

	if err := os.WriteFile(inputFile, testData, 0644); err != nil {
		fmt.Printf("Failed to write test file: %v\n", err)
		return
	}
	defer os.Remove(inputFile)
	defer os.Remove(encryptedFile)
	defer os.Remove(decryptedFile)

	fmt.Printf("Generated encryption key: %s\n", key)

	if err := encryptor.EncryptFile(inputFile, encryptedFile); err != nil {
		fmt.Printf("Encryption failed: %v\n", err)
		return
	}
	fmt.Println("File encrypted successfully")

	if err := encryptor.DecryptFile(encryptedFile, decryptedFile); err != nil {
		fmt.Printf("Decryption failed: %v\n", err)
		return
	}
	fmt.Println("File decrypted successfully")

	decryptedData, err := os.ReadFile(decryptedFile)
	if err != nil {
		fmt.Printf("Failed to read decrypted file: %v\n", err)
		return
	}

	if string(decryptedData) == string(testData) {
		fmt.Println("Encryption/decryption test passed")
	} else {
		fmt.Println("Encryption/decryption test failed")
	}
}