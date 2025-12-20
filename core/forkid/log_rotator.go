
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type RotatingLogger struct {
	filePath    string
	maxSize     int64
	currentSize int64
	file        *os.File
	compression bool
}

func NewRotatingLogger(filePath string, maxSizeMB int, compression bool) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	return &RotatingLogger{
		filePath:    filePath,
		maxSize:     maxSize,
		currentSize: info.Size(),
		file:        file,
		compression: compression,
	}, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	if rl.currentSize+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.file.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	rl.file.Close()

	ext := filepath.Ext(rl.filePath)
	baseName := strings.TrimSuffix(rl.filePath, ext)
	timestamp := time.Now().Format("20060102_150405")
	archivePath := fmt.Sprintf("%s_%s%s", baseName, timestamp, ext)

	if err := os.Rename(rl.filePath, archivePath); err != nil {
		return err
	}

	if rl.compression {
		go compressFile(archivePath)
	}

	file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	rl.file = file
	rl.currentSize = 0
	return nil
}

func compressFile(source string) {
	dest := source + ".gz"
	log.Printf("Compressing %s to %s", source, dest)
	// Compression implementation would go here
	// For now just simulate by renaming
	if err := os.Rename(source, dest); err != nil {
		log.Printf("Failed to compress %s: %v", source, err)
	}
}

func (rl *RotatingLogger) Close() error {
	return rl.file.Close()
}

func main() {
	logger, err := NewRotatingLogger("app.log", 10, true)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	customLog := log.New(io.MultiWriter(os.Stdout, logger), "", log.LstdFlags)

	for i := 0; i < 100; i++ {
		customLog.Printf("Log entry number %d", i)
		time.Sleep(100 * time.Millisecond)
	}
}