package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type LogRotator struct {
	filePath    string
	maxSize     int64
	currentSize int64
	file        *os.File
}

func NewLogRotator(path string, maxSizeMB int) (*LogRotator, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	return &LogRotator{
		filePath:    path,
		maxSize:     maxSize,
		currentSize: info.Size(),
		file:        file,
	}, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	if lr.currentSize+int64(len(p)) > lr.maxSize {
		if err := lr.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := lr.file.Write(p)
	lr.currentSize += int64(n)
	return n, err
}

func (lr *LogRotator) rotate() error {
	lr.file.Close()

	timestamp := time.Now().Format("20060102_150405")
	archivePath := fmt.Sprintf("%s.%s.gz", lr.filePath, timestamp)

	if err := compressFile(lr.filePath, archivePath); err != nil {
		return err
	}

	if err := os.Truncate(lr.filePath, 0); err != nil {
		return err
	}

	file, err := os.OpenFile(lr.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	lr.file = file
	lr.currentSize = 0
	return nil
}

func compressFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	gzWriter := gzip.NewWriter(dstFile)
	defer gzWriter.Close()

	_, err = io.Copy(gzWriter, srcFile)
	return err
}

func (lr *LogRotator) Close() error {
	return lr.file.Close()
}

func main() {
	rotator, err := NewLogRotator("app.log", 10)
	if err != nil {
		fmt.Printf("Failed to create log rotator: %v\n", err)
		return
	}
	defer rotator.Close()

	for i := 0; i < 1000; i++ {
		logEntry := fmt.Sprintf("[%s] Log entry %d: Some sample log data here\n",
			time.Now().Format(time.RFC3339), i)
		if _, err := rotator.Write([]byte(logEntry)); err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed. Check app.log and compressed archives.")
}