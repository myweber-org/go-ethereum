
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func encryptData(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func decryptData(ciphertext []byte, key []byte) ([]byte, error) {
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
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func generateRandomKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_encryptor <encrypt|decrypt>")
		return
	}

	operation := os.Args[1]
	key, _ := generateRandomKey()
	keyHex := hex.EncodeToString(key)
	fmt.Printf("Generated key: %s\n", keyHex)

	sampleData := []byte("Confidential information requiring encryption")

	switch operation {
	case "encrypt":
		encrypted, err := encryptData(sampleData, key)
		if err != nil {
			fmt.Printf("Encryption error: %v\n", err)
			return
		}
		fmt.Printf("Encrypted data (hex): %s\n", hex.EncodeToString(encrypted))

	case "decrypt":
		encrypted, _ := encryptData(sampleData, key)
		decrypted, err := decryptData(encrypted, key)
		if err != nil {
			fmt.Printf("Decryption error: %v\n", err)
			return
		}
		fmt.Printf("Decrypted data: %s\n", string(decrypted))

	default:
		fmt.Println("Invalid operation. Use 'encrypt' or 'decrypt'")
	}
}