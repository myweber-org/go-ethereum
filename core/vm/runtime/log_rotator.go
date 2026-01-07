
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
	mu       sync.Mutex
	filename string
	file     *os.File
	size     int64
}

func NewRotatingFile(filename string) (*RotatingFile, error) {
	rf := &RotatingFile{filename: filename}
	if err := rf.openFile(); err != nil {
		return nil, err
	}
	return rf, nil
}

func (rf *RotatingFile) Write(p []byte) (n int, err error) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if rf.size+int64(len(p)) > maxFileSize {
		if err := rf.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = rf.file.Write(p)
	rf.size += int64(n)
	return n, err
}

func (rf *RotatingFile) openFile() error {
	info, err := os.Stat(rf.filename)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if info != nil {
		rf.size = info.Size()
	}

	rf.file, err = os.OpenFile(rf.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return err
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
				os.Rename(oldName, newName)
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
		return rf.filename + ".1"
	}
	ext := filepath.Ext(rf.filename)
	base := rf.filename[:len(rf.filename)-len(ext)]
	return fmt.Sprintf("%s.%d%s.gz", base, index+1, ext)
}

func (rf *RotatingFile) compressOldFiles() error {
	for i := 1; i <= backupCount; i++ {
		filename := fmt.Sprintf("%s.%d", rf.filename, i)
		if _, err := os.Stat(filename); err == nil {
			if err := compressFile(filename); err != nil {
				return err
			}
			os.Remove(filename)
		}
	}
	return nil
}

func compressFile(filename string) error {
	in, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(filename + ".gz")
	if err != nil {
		return err
	}
	defer out.Close()

	gz := gzip.NewWriter(out)
	defer gz.Close()

	_, err = io.Copy(gz, in)
	return err
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

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: This is a test log message\n",
			time.Now().Format("2006-01-02 15:04:05"), i)
		if _, err := logFile.Write([]byte(msg)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	if err := logFile.compressOldFiles(); err != nil {
		fmt.Printf("Compression error: %v\n", err)
	}

	fmt.Println("Log rotation test completed")
}