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
		return fmt.Errorf("read file error: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation error: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("GCM mode error: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("nonce generation error: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
		return fmt.Errorf("write file error: %w", err)
	}

	return nil
}

func decryptFile(inputPath, outputPath string, key []byte) error {
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read file error: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation error: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("GCM mode error: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("decryption error: %w", err)
	}

	if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
		return fmt.Errorf("write file error: %w", err)
	}

	return nil
}

func main() {
	key := []byte("32-byte-long-key-here-123456789012")
	
	err := encryptFile("plain.txt", "encrypted.bin", key)
	if err != nil {
		fmt.Printf("Encryption failed: %v\n", err)
		return
	}
	fmt.Println("File encrypted successfully")

	err = decryptFile("encrypted.bin", "decrypted.txt", key)
	if err != nil {
		fmt.Printf("Decryption failed: %v\n", err)
		return
	}
	fmt.Println("File decrypted successfully")
}
package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "errors"
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
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
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
        return nil, errors.New("ciphertext too short")
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
        fmt.Println("Usage: go run file_encryptor.go <encrypt|decrypt|keygen>")
        return
    }

    switch os.Args[1] {
    case "encrypt":
        if len(os.Args) != 4 {
            fmt.Println("Usage: go run file_encryptor.go encrypt <input_file> <key_base64>")
            return
        }

        data, err := os.ReadFile(os.Args[2])
        if err != nil {
            fmt.Printf("Error reading file: %v\n", err)
            return
        }

        key, err := base64.StdEncoding.DecodeString(os.Args[3])
        if err != nil {
            fmt.Printf("Error decoding key: %v\n", err)
            return
        }

        encrypted, err := encryptData(data, key)
        if err != nil {
            fmt.Printf("Error encrypting: %v\n", err)
            return
        }

        outputFile := os.Args[2] + ".enc"
        if err := os.WriteFile(outputFile, encrypted, 0644); err != nil {
            fmt.Printf("Error writing file: %v\n", err)
            return
        }

        fmt.Printf("Encrypted file saved as: %s\n", outputFile)

    case "decrypt":
        if len(os.Args) != 4 {
            fmt.Println("Usage: go run file_encryptor.go decrypt <encrypted_file> <key_base64>")
            return
        }

        data, err := os.ReadFile(os.Args[2])
        if err != nil {
            fmt.Printf("Error reading file: %v\n", err)
            return
        }

        key, err := base64.StdEncoding.DecodeString(os.Args[3])
        if err != nil {
            fmt.Printf("Error decoding key: %v\n", err)
            return
        }

        decrypted, err := decryptData(data, key)
        if err != nil {
            fmt.Printf("Error decrypting: %v\n", err)
            return
        }

        outputFile := os.Args[2] + ".dec"
        if err := os.WriteFile(outputFile, decrypted, 0644); err != nil {
            fmt.Printf("Error writing file: %v\n", err)
            return
        }

        fmt.Printf("Decrypted file saved as: %s\n", outputFile)

    case "keygen":
        key, err := generateRandomKey()
        if err != nil {
            fmt.Printf("Error generating key: %v\n", err)
            return
        }

        keyBase64 := base64.StdEncoding.EncodeToString(key)
        fmt.Printf("Generated key: %s\n", keyBase64)
        fmt.Println("Save this key securely!")

    default:
        fmt.Println("Invalid command. Use: encrypt, decrypt, or keygen")
    }
}