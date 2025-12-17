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

const (
	maxFileSize    = 10 * 1024 * 1024 // 10MB
	maxBackupFiles = 5
	logDir         = "./logs"
)

type RotatingLogger struct {
	currentFile *os.File
	filePath    string
	baseName    string
	fileIndex   int
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	rl := &RotatingLogger{
		baseName:  strings.TrimSuffix(baseName, ".log"),
		fileIndex: 0,
	}

	rl.filePath = rl.generateFilePath()
	if err := rl.openFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) generateFilePath() string {
	return filepath.Join(logDir, fmt.Sprintf("%s_%d.log", rl.baseName, rl.fileIndex))
}

func (rl *RotatingLogger) openFile() error {
	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	rl.currentFile = file
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	fileInfo, err := rl.currentFile.Stat()
	if err != nil {
		return 0, err
	}

	if fileInfo.Size()+int64(len(p)) > maxFileSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	return rl.currentFile.Write(p)
}

func (rl *RotatingLogger) rotate() error {
	rl.currentFile.Close()
	rl.fileIndex++

	if rl.fileIndex > maxBackupFiles {
		rl.fileIndex = 0
		if err := rl.cleanupOldFiles(); err != nil {
			return err
		}
	}

	rl.filePath = rl.generateFilePath()
	return rl.openFile()
}

func (rl *RotatingLogger) cleanupOldFiles() error {
	files, err := filepath.Glob(filepath.Join(logDir, rl.baseName+"_*.log"))
	if err != nil {
		return err
	}

	if len(files) > maxBackupFiles {
		for i := 0; i < len(files)-maxBackupFiles; i++ {
			if err := os.Remove(files[i]); err != nil {
				return err
			}
		}
	}

	return nil
}

func (rl *RotatingLogger) Close() error {
	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	customLog := log.New(io.MultiWriter(os.Stdout, logger), "", log.LstdFlags)

	for i := 0; i < 100; i++ {
		customLog.Printf("Log entry %d at %s", i, time.Now().Format(time.RFC3339))
		time.Sleep(100 * time.Millisecond)
	}
}