
package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	backupCount = 5
)

type RotatingLogger struct {
	filename    string
	currentSize int64
	file        *os.File
	mu          sync.Mutex
}

func NewRotatingLogger(filename string) (*RotatingLogger, error) {
	rl := &RotatingLogger{filename: filename}
	if err := rl.openFile(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openFile() error {
	info, err := os.Stat(rl.filename)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if info != nil {
		rl.currentSize = info.Size()
	}

	file, err := os.OpenFile(rl.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	rl.file = file
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > maxFileSize {
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
	if err := rl.file.Close(); err != nil {
		return err
	}

	for i := backupCount - 1; i >= 0; i-- {
		src := rl.backupFilename(i)
		dst := rl.backupFilename(i + 1)

		if _, err := os.Stat(src); err == nil {
			if i == backupCount-1 {
				os.Remove(src)
			} else {
				if err := rl.compressAndMove(src, dst); err != nil {
					return err
				}
			}
		}
	}

	if err := os.Rename(rl.filename, rl.backupFilename(0)); err != nil {
		return err
	}

	return rl.openFile()
}

func (rl *RotatingLogger) backupFilename(index int) string {
	if index == 0 {
		return rl.filename + ".1"
	}
	return fmt.Sprintf("%s.%d.gz", rl.filename, index)
}

func (rl *RotatingLogger) compressAndMove(src, dst string) error {
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

	if _, err := io.Copy(gzWriter, srcFile); err != nil {
		return err
	}

	return os.Remove(src)
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.file.Close()
}

func main() {
	logger, err := NewRotatingLogger("app.log")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: This is a test log message\n",
			time.Now().Format(time.RFC3339), i)
		logger.Write([]byte(msg))
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}