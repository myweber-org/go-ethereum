package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	maxLogSize    = 1024 * 1024 // 1MB
	maxBackupLogs = 5
	logFileName   = "app.log"
)

type RotatingLogger struct {
	currentSize int64
	file        *os.File
	logger      *log.Logger
}

func NewRotatingLogger() (*RotatingLogger, error) {
	rl := &RotatingLogger{}
	if err := rl.openLogFile(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openLogFile() error {
	if rl.file != nil {
		rl.file.Close()
	}

	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.file = file
	rl.currentSize = info.Size()
	rl.logger = log.New(rl.file, "", log.LstdFlags)
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	if rl.currentSize+int64(len(p)) > maxLogSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = rl.file.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	rl.file.Close()

	for i := maxBackupLogs - 1; i >= 0; i-- {
		oldName := fmt.Sprintf("%s.%d", logFileName, i)
		newName := fmt.Sprintf("%s.%d", logFileName, i+1)

		if _, err := os.Stat(oldName); err == nil {
			os.Rename(oldName, newName)
		}
	}

	backupName := fmt.Sprintf("%s.0", logFileName)
	os.Rename(logFileName, backupName)

	if err := rl.openLogFile(); err != nil {
		return err
	}

	rl.cleanupOldLogs()
	return nil
}

func (rl *RotatingLogger) cleanupOldLogs() {
	for i := maxBackupLogs; i < 10; i++ {
		filename := fmt.Sprintf("%s.%d", logFileName, i)
		os.Remove(filename)
	}
}

func (rl *RotatingLogger) Log(message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] %s\n", timestamp, message)
	rl.Write([]byte(logEntry))
}

func main() {
	logger, err := NewRotatingLogger()
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.file.Close()

	for i := 0; i < 1000; i++ {
		logger.Log(fmt.Sprintf("Log entry number %d", i))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}