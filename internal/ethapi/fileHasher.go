package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

type FileHash struct {
	Path string
	SHA256 string
	MD5    string
}

func computeFileHash(path string) (*FileHash, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	sha256Hash := sha256.New()
	md5Hash := md5.New()

	multiWriter := io.MultiWriter(sha256Hash, md5Hash)

	if _, err := io.Copy(multiWriter, file); err != nil {
		return nil, err
	}

	return &FileHash{
		Path:   path,
		SHA256: hex.EncodeToString(sha256Hash.Sum(nil)),
		MD5:    hex.EncodeToString(md5Hash.Sum(nil)),
	}, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: fileHasher <filepath>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	hashResult, err := computeFileHash(filePath)
	if err != nil {
		fmt.Printf("Error hashing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("File: %s\n", hashResult.Path)
	fmt.Printf("SHA256: %s\n", hashResult.SHA256)
	fmt.Printf("MD5:    %s\n", hashResult.MD5)
}