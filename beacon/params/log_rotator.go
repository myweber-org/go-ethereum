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

type RotatingLogger struct {
	mu          sync.Mutex
	currentFile *os.File
	basePath    string
	maxSize     int64
	fileCount   int
	maxFiles    int
	currentSize int64
}

func NewRotatingLogger(basePath string, maxSize int64, maxFiles int) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: basePath,
		maxSize:  maxSize,
		maxFiles: maxFiles,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.currentFile.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	oldPath := rl.currentFilePath()
	archivePath := fmt.Sprintf("%s.%d.gz", oldPath, time.Now().Unix())

	if err := rl.compressFile(oldPath, archivePath); err != nil {
		return err
	}

	os.Remove(oldPath)

	rl.fileCount++
	if rl.fileCount > rl.maxFiles {
		rl.cleanOldFiles()
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	gz := gzip.NewWriter(out)
	defer gz.Close()

	_, err = io.Copy(gz, in)
	return err
}

func (rl *RotatingLogger) openCurrentFile() error {
	file, err := os.OpenFile(rl.currentFilePath(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.currentFile = file
	rl.currentSize = info.Size()
	return nil
}

func (rl *RotatingLogger) currentFilePath() string {
	return rl.basePath
}

func (rl *RotatingLogger) cleanOldFiles() {
	pattern := rl.basePath + ".*.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) > rl.maxFiles {
		for i := 0; i < len(matches)-rl.maxFiles; i++ {
			os.Remove(matches[i])
		}
	}
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
	logger, err := NewRotatingLogger("app.log", 1024*1024, 5)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		logger.Write([]byte(fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))))
		time.Sleep(10 * time.Millisecond)
	}
}
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
	file       *os.File
	size       int64
	basePath   string
	currentNum int
}

func NewRotatingFile(path string) (*RotatingFile, error) {
	rf := &RotatingFile{
		basePath: path,
	}
	if err := rf.rotateIfNeeded(); err != nil {
		return nil, err
	}
	return rf, nil
}

func (rf *RotatingFile) Write(p []byte) (int, error) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if err := rf.rotateIfNeeded(); err != nil {
		return 0, err
	}

	n, err := rf.file.Write(p)
	if err == nil {
		rf.size += int64(n)
	}
	return n, err
}

func (rf *RotatingFile) rotateIfNeeded() error {
	if rf.file != nil && rf.size < maxFileSize {
		return nil
	}

	if rf.file != nil {
		if err := rf.file.Close(); err != nil {
			return err
		}
		if err := rf.compressCurrent(); err != nil {
			return err
		}
		rf.cleanOldBackups()
	}

	filename := rf.basePath
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rf.file = file
	rf.size = stat.Size()
	return nil
}

func (rf *RotatingFile) compressCurrent() error {
	srcPath := rf.basePath
	dstPath := fmt.Sprintf("%s.%d.gz", rf.basePath, rf.currentNum)

	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)
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

	if err := os.Remove(srcPath); err != nil {
		return err
	}

	rf.currentNum = (rf.currentNum + 1) % backupCount
	return nil
}

func (rf *RotatingFile) cleanOldBackups() {
	for i := 0; i < backupCount; i++ {
		pattern := fmt.Sprintf("%s.*.gz", rf.basePath)
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}

		if len(matches) <= backupCount {
			return
		}

		oldest := ""
		var oldestTime time.Time
		for _, match := range matches {
			info, err := os.Stat(match)
			if err != nil {
				continue
			}
			if oldest == "" || info.ModTime().Before(oldestTime) {
				oldest = match
				oldestTime = info.ModTime()
			}
		}

		if oldest != "" {
			os.Remove(oldest)
		}
	}
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
	logFile, err := NewRotatingFile("app.log")
	if err != nil {
		fmt.Printf("Failed to create log file: %v\n", err)
		return
	}
	defer logFile.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("[%s] Log entry %d: Sample log message\n",
			time.Now().Format(time.RFC3339), i)
		if _, err := logFile.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}