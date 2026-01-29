package main

import (
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

func NewRotatingLogger(path string, maxSize int64, maxBackups int) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		filePath:   path,
		maxSize:    maxSize,
		maxBackups: maxBackups,
	}
	if err := rl.openCurrent(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openCurrent() error {
	if rl.current != nil {
		rl.current.Close()
	}
	f, err := os.OpenFile(rl.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	info, err := f.Stat()
	if err != nil {
		f.Close()
		return err
	}
	rl.current = f
	rl.size = info.Size()
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	if rl.size+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			log.Printf("rotate failed: %v", err)
		}
	}
	n, err := rl.current.Write(p)
	if err == nil {
		rl.size += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.current.Close(); err != nil {
		return err
	}
	ext := filepath.Ext(rl.filePath)
	base := rl.filePath[:len(rl.filePath)-len(ext)]
	timestamp := time.Now().Format("20060102_150405")
	backupPath := base + "_" + timestamp + ext
	if err := os.Rename(rl.filePath, backupPath); err != nil {
		return err
	}
	if err := rl.openCurrent(); err != nil {
		return err
	}
	return rl.cleanupOld()
}

func (rl *RotatingLogger) cleanupOld() error {
	if rl.maxBackups <= 0 {
		return nil
	}
	ext := filepath.Ext(rl.filePath)
	base := rl.filePath[:len(rl.filePath)-len(ext)]
	pattern := base + "_*" + ext
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	if len(matches) <= rl.maxBackups {
		return nil
	}
	toDelete := matches[:len(matches)-rl.maxBackups]
	for _, path := range toDelete {
		if err := os.Remove(path); err != nil {
			return err
		}
	}
	return nil
}

func (rl *RotatingLogger) Close() error {
	if rl.current != nil {
		return rl.current.Close()
	}
	return nil
}