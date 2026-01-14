package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "errors"
    "io"
    "os"
)

func encryptData(key []byte, plaintext []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    ciphertext := make([]byte, aes.BlockSize+len(plaintext))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return nil, err
    }

    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

    return ciphertext, nil
}

func decryptData(key []byte, ciphertext []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    if len(ciphertext) < aes.BlockSize {
        return nil, errors.New("ciphertext too short")
    }

    iv := ciphertext[:aes.BlockSize]
    ciphertext = ciphertext[aes.BlockSize:]

    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(ciphertext, ciphertext)

    return ciphertext, nil
}

func encryptFile(key []byte, inputPath string, outputPath string) error {
    data, err := os.ReadFile(inputPath)
    if err != nil {
        return err
    }

    encrypted, err := encryptData(key, data)
    if err != nil {
        return err
    }

    encoded := base64.StdEncoding.EncodeToString(encrypted)
    return os.WriteFile(outputPath, []byte(encoded), 0644)
}

func decryptFile(key []byte, inputPath string, outputPath string) error {
    encoded, err := os.ReadFile(inputPath)
    if err != nil {
        return err
    }

    encrypted, err := base64.StdEncoding.DecodeString(string(encoded))
    if err != nil {
        return err
    }

    decrypted, err := decryptData(key, encrypted)
    if err != nil {
        return err
    }

    return os.WriteFile(outputPath, decrypted, 0644)
}