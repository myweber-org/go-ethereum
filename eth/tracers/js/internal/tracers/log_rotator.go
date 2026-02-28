
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
	maxFileSize = 10 * 1024 * 1024
	backupCount = 5
)

type RotatingFile struct {
	mu         sync.Mutex
	filename   string
	file       *os.File
	currentSize int64
}

func NewRotatingFile(filename string) (*RotatingFile, error) {
	rf := &RotatingFile{
		filename: filename,
	}

	if err := rf.openFile(); err != nil {
		return nil, err
	}

	return rf, nil
}

func (rf *RotatingFile) openFile() error {
	file, err := os.OpenFile(rf.filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rf.file = file
	rf.currentSize = stat.Size()
	return nil
}

func (rf *RotatingFile) Write(p []byte) (int, error) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if rf.currentSize+int64(len(p)) > maxFileSize {
		if err := rf.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rf.file.Write(p)
	if err == nil {
		rf.currentSize += int64(n)
	}
	return n, err
}

func (rf *RotatingFile) rotate() error {
	if rf.file != nil {
		rf.file.Close()
	}

	for i := backupCount - 1; i >= 0; i-- {
		oldName := rf.backupName(i)
		newName := rf.backupName(i + 1)

		if _, err := os.Stat(oldName); err == nil {
			if i == backupCount-1 {
				os.Remove(oldName)
			} else {
				if err := rf.compressAndMove(oldName, newName); err != nil {
					return err
				}
			}
		}
	}

	if err := os.Rename(rf.filename, rf.backupName(0)); err != nil {
		return err
	}

	return rf.openFile()
}

func (rf *RotatingFile) backupName(index int) string {
	if index == 0 {
		return rf.filename
	}
	ext := filepath.Ext(rf.filename)
	base := rf.filename[:len(rf.filename)-len(ext)]
	return fmt.Sprintf("%s.%d%s.gz", base, index, ext)
}

func (rf *RotatingFile) compressAndMove(src, dst string) error {
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
	if err != nil {
		return err
	}

	return os.Remove(src)
}

func (rf *RotatingFile) Close() error {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if rf.file != nil {
		return rf.file.Close()
	}
	return nil
}

func main() {
	logFile, err := NewRotatingFile("application.log")
	if err != nil {
		fmt.Printf("Failed to create log file: %v\n", err)
		return
	}
	defer logFile.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("[%s] Log entry %d: This is a sample log message\n", 
			time.Now().Format(time.RFC3339), i)
		if _, err := logFile.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}