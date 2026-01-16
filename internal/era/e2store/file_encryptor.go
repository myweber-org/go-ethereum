
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

	return os.WriteFile(outputPath, ciphertext, 0644)
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

	return os.WriteFile(outputPath, plaintext, 0644)
}

func generateKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_encryptor <encrypt|decrypt|genkey> [arguments]")
		return
	}

	switch os.Args[1] {
	case "genkey":
		key, err := generateKey()
		if err != nil {
			fmt.Printf("Key generation failed: %v\n", err)
			return
		}
		fmt.Printf("Generated key: %x\n", key)

	case "encrypt":
		if len(os.Args) != 5 {
			fmt.Println("Usage: file_encryptor encrypt <input> <output> <hex_key>")
			return
		}
		var key []byte
		fmt.Sscanf(os.Args[4], "%x", &key)
		if err := encryptFile(os.Args[2], os.Args[3], key); err != nil {
			fmt.Printf("Encryption failed: %v\n", err)
		} else {
			fmt.Println("File encrypted successfully")
		}

	case "decrypt":
		if len(os.Args) != 5 {
			fmt.Println("Usage: file_encryptor decrypt <input> <output> <hex_key>")
			return
		}
		var key []byte
		fmt.Sscanf(os.Args[4], "%x", &key)
		if err := decryptFile(os.Args[2], os.Args[3], key); err != nil {
			fmt.Printf("Decryption failed: %v\n", err)
		} else {
			fmt.Println("File decrypted successfully")
		}

	default:
		fmt.Println("Unknown command")
	}
}