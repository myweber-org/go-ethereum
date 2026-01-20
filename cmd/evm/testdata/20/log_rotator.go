package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type RotatingLogger struct {
	filePath   string
	maxSize    int64
	maxBackups int
	current    *os.File
	size       int64
}

func NewRotatingLogger(filePath string, maxSize int64, maxBackups int) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		filePath:   filePath,
		maxSize:    maxSize,
		maxBackups: maxBackups,
	}
	if err := rl.openFile(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openFile() error {
	info, err := os.Stat(rl.filePath)
	if err == nil {
		rl.size = info.Size()
	} else if os.IsNotExist(err) {
		rl.size = 0
	} else {
		return err
	}

	file, err := os.OpenFile(rl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	rl.current = file
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	if rl.size+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}
	n, err = rl.current.Write(p)
	rl.size += int64(n)
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.current.Close(); err != nil {
		return err
	}

	for i := rl.maxBackups - 1; i >= 0; i-- {
		oldPath := rl.backupPath(i)
		newPath := rl.backupPath(i + 1)
		if _, err := os.Stat(oldPath); err == nil {
			if err := os.Rename(oldPath, newPath); err != nil {
				return err
			}
		}
	}

	if err := os.Rename(rl.filePath, rl.backupPath(0)); err != nil && !os.IsNotExist(err) {
		return err
	}

	return rl.openFile()
}

func (rl *RotatingLogger) backupPath(index int) string {
	if index == 0 {
		return rl.filePath + ".1"
	}
	return fmt.Sprintf("%s.%d", rl.filePath, index+1)
}

func (rl *RotatingLogger) Close() error {
	if rl.current != nil {
		return rl.current.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app.log", 1024*1024, 5)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	log.SetOutput(io.MultiWriter(os.Stdout, logger))

	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d at %v", i, time.Now())
		time.Sleep(10 * time.Millisecond)
	}
}