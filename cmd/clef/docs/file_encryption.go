
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

func deriveKey(passphrase string) []byte {
    hash := sha256.Sum256([]byte(passphrase))
    return hash[:]
}

func encryptFile(inputPath, outputPath, passphrase string) error {
    plaintext, err := os.ReadFile(inputPath)
    if err != nil {
        return err
    }

    key := deriveKey(passphrase)
    block, err := aes.NewCipher(key)
    if err != nil {
        return err
    }

    ciphertext := make([]byte, aes.BlockSize+len(plaintext))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return err
    }

    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

    return os.WriteFile(outputPath, ciphertext, 0644)
}

func decryptFile(inputPath, outputPath, passphrase string) error {
    ciphertext, err := os.ReadFile(inputPath)
    if err != nil {
        return err
    }

    if len(ciphertext) < aes.BlockSize {
        return errors.New("ciphertext too short")
    }

    key := deriveKey(passphrase)
    block, err := aes.NewCipher(key)
    if err != nil {
        return err
    }

    iv := ciphertext[:aes.BlockSize]
    ciphertext = ciphertext[aes.BlockSize:]

    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(ciphertext, ciphertext)

    return os.WriteFile(outputPath, ciphertext, 0644)
}

func main() {
    if len(os.Args) < 5 {
        fmt.Println("Usage: go run file_encryption.go <encrypt|decrypt> <input> <output> <passphrase>")
        os.Exit(1)
    }

    mode := os.Args[1]
    input := os.Args[2]
    output := os.Args[3]
    passphrase := os.Args[4]

    var err error
    switch mode {
    case "encrypt":
        err = encryptFile(input, output, passphrase)
    case "decrypt":
        err = decryptFile(input, output, passphrase)
    default:
        fmt.Println("Invalid mode. Use 'encrypt' or 'decrypt'")
        os.Exit(1)
    }

    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Operation completed successfully: %s -> %s\n", input, output)
}