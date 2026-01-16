
package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type RotatingLogger struct {
	mu          sync.Mutex
	currentFile *os.File
	basePath    string
	maxSize     int64
	currentSize int64
	fileIndex   int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath:  basePath,
		maxSize:   int64(maxSizeMB) * 1024 * 1024,
		fileIndex: 0,
	}

	if err := rl.openNextFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openNextFile() error {
	rl.fileIndex++
	filename := rl.basePath + "_" + strconv.Itoa(rl.fileIndex) + ".log"
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	rl.currentFile = file
	info, err := file.Stat()
	if err != nil {
		rl.currentSize = 0
	} else {
		rl.currentSize = info.Size()
	}

	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > rl.maxSize {
		if err := rl.openNextFile(); err != nil {
			return 0, err
		}
	}

	n, err := rl.currentFile.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app_log", 10)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	customLog := log.New(logger, "", log.LstdFlags)

	for i := 0; i < 1000; i++ {
		customLog.Printf("Log entry number %d at %v", i, time.Now())
		time.Sleep(10 * time.Millisecond)
	}
}