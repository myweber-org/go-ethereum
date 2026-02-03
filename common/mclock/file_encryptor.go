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

func generateKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func encryptFile(inputPath, outputPath string, key []byte) error {
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

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

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	err = os.WriteFile(outputPath, ciphertext, 0644)
	if err != nil {
		return err
	}

	return nil
}

func decryptFile(inputPath, outputPath string, key []byte) error {
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
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

	err = os.WriteFile(outputPath, plaintext, 0644)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run file_encryptor.go <encrypt|decrypt> <input> <output>")
		fmt.Println("Example: go run file_encryptor.go encrypt secret.txt encrypted.bin")
		os.Exit(1)
	}

	action := os.Args[1]
	inputFile := os.Args[2]
	outputFile := os.Args[3]

	key, err := generateKey()
	if err != nil {
		fmt.Printf("Error generating key: %v\n", err)
		os.Exit(1)
	}

	switch action {
	case "encrypt":
		err = encryptFile(inputFile, outputFile, key)
		if err != nil {
			fmt.Printf("Encryption failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("File encrypted successfully.\nKey: %s\n", hex.EncodeToString(key))

	case "decrypt":
		fmt.Print("Enter encryption key (hex): ")
		var keyHex string
		_, err = fmt.Scanln(&keyHex)
		if err != nil {
			fmt.Printf("Error reading key: %v\n", err)
			os.Exit(1)
		}

		key, err = hex.DecodeString(keyHex)
		if err != nil {
			fmt.Printf("Invalid key format: %v\n", err)
			os.Exit(1)
		}

		err = decryptFile(inputFile, outputFile, key)
		if err != nil {
			fmt.Printf("Decryption failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("File decrypted successfully.")

	default:
		fmt.Println("Invalid action. Use 'encrypt' or 'decrypt'.")
		os.Exit(1)
	}
}