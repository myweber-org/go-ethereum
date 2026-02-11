
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
	maxBackups  = 5
)

type RotatingLogger struct {
	mu         sync.Mutex
	file       *os.File
	size       int64
	basePath   string
	currentNum int
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: path,
	}
	if err := rl.openCurrent(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.size+int64(len(p)) > maxFileSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.file.Write(p)
	rl.size += int64(n)
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if rl.file != nil {
		rl.file.Close()
	}

	rl.currentNum++
	if rl.currentNum > maxBackups {
		rl.currentNum = 1
	}

	oldPath := rl.basePath
	if rl.currentNum > 1 {
		oldPath = fmt.Sprintf("%s.%d", rl.basePath, rl.currentNum-1)
	}

	if err := rl.compressFile(oldPath); err != nil {
		return err
	}

	return rl.openCurrent()
}

func (rl *RotatingLogger) compressFile(path string) error {
	src, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer src.Close()

	dstPath := path + ".gz"
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	gz := gzip.NewWriter(dst)
	defer gz.Close()

	if _, err := io.Copy(gz, src); err != nil {
		return err
	}

	os.Remove(path)
	return nil
}

func (rl *RotatingLogger) openCurrent() error {
	dir := filepath.Dir(rl.basePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.file = file
	rl.size = info.Size()
	return nil
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.file != nil {
		return rl.file.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("/var/log/myapp/app.log")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: Application event processed\n",
			time.Now().Format(time.RFC3339), i)
		logger.Write([]byte(msg))
		time.Sleep(10 * time.Millisecond)
	}
}
package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	maxFileSize    = 10 * 1024 * 1024 // 10MB
	backupCount    = 5
	checkInterval  = 30 * time.Second
	logDir         = "./logs"
	currentLogName = "application.log"
)

type LogRotator struct {
	currentFile *os.File
	currentSize int64
	mu          sync.Mutex
	stopChan    chan struct{}
}

func NewLogRotator() (*LogRotator, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	logPath := filepath.Join(logDir, currentLogName)
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to stat log file: %w", err)
	}

	rotator := &LogRotator{
		currentFile: file,
		currentSize: info.Size(),
		stopChan:    make(chan struct{}),
	}

	go rotator.monitor()
	return rotator, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	n, err := lr.currentFile.Write(p)
	if err != nil {
		return n, err
	}

	lr.currentSize += int64(n)
	if lr.currentSize >= maxFileSize {
		if err := lr.rotate(); err != nil {
			log.Printf("rotation failed: %v", err)
		}
	}
	return n, nil
}

func (lr *LogRotator) rotate() error {
	if err := lr.currentFile.Close(); err != nil {
		return fmt.Errorf("failed to close current file: %w", err)
	}

	// Compress and rename current log
	timestamp := time.Now().Format("20060102_150405")
	oldPath := filepath.Join(logDir, currentLogName)
	newPath := filepath.Join(logDir, fmt.Sprintf("application_%s.log.gz", timestamp))

	if err := compressFile(oldPath, newPath); err != nil {
		return fmt.Errorf("compression failed: %w", err)
	}

	// Remove old backups
	if err := lr.cleanupOldBackups(); err != nil {
		log.Printf("backup cleanup failed: %v", err)
	}

	// Create new log file
	file, err := os.OpenFile(oldPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create new log file: %w", err)
	}

	lr.currentFile = file
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

	if _, err := io.Copy(gzWriter, srcFile); err != nil {
		return err
	}

	// Remove original uncompressed file
	return os.Remove(src)
}

func (lr *LogRotator) cleanupOldBackups() error {
	files, err := filepath.Glob(filepath.Join(logDir, "application_*.log.gz"))
	if err != nil {
		return err
	}

	if len(files) <= backupCount {
		return nil
	}

	// Sort by modification time (oldest first)
	for i := 0; i < len(files)-backupCount; i++ {
		if err := os.Remove(files[i]); err != nil {
			return err
		}
	}
	return nil
}

func (lr *LogRotator) monitor() {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			lr.mu.Lock()
			if lr.currentSize >= maxFileSize {
				if err := lr.rotate(); err != nil {
					log.Printf("periodic rotation failed: %v", err)
				}
			}
			lr.mu.Unlock()
		case <-lr.stopChan:
			return
		}
	}
}

func (lr *LogRotator) Close() error {
	close(lr.stopChan)
	lr.mu.Lock()
	defer lr.mu.Unlock()
	return lr.currentFile.Close()
}

func main() {
	rotator, err := NewLogRotator()
	if err != nil {
		log.Fatal(err)
	}
	defer rotator.Close()

	log.SetOutput(rotator)

	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d: %s", i, time.Now().Format(time.RFC3339))
		time.Sleep(100 * time.Millisecond)
	}
}package main

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
	mu         sync.Mutex
	file       *os.File
	currentSize int64
	basePath   string
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
	rl := &RotatingLogger{basePath: path}
	if err := rl.openCurrent(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openCurrent() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.file != nil {
		rl.file.Close()
	}

	file, err := os.OpenFile(rl.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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

	// Compress and rename old logs
	for i := backupCount - 1; i >= 0; i-- {
		oldPath := rl.backupPath(i)
		newPath := rl.backupPath(i + 1)

		if _, err := os.Stat(oldPath); err == nil {
			if i == backupCount-1 {
				os.Remove(oldPath)
			} else {
				if err := rl.compressFile(oldPath, newPath); err != nil {
					return err
				}
				os.Remove(oldPath)
			}
		}
	}

	// Move current to first backup
	if err := rl.compressFile(rl.basePath, rl.backupPath(0)); err != nil {
		return err
	}
	os.Remove(rl.basePath)

	return rl.openCurrent()
}

func (rl *RotatingLogger) backupPath(index int) string {
	if index == 0 {
		return rl.basePath + ".1"
	}
	return fmt.Sprintf("%s.%d.gz", rl.basePath, index)
}

func (rl *RotatingLogger) compressFile(src, dst string) error {
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

	gz := gzip.NewWriter(dstFile)
	defer gz.Close()

	_, err = io.Copy(gz, srcFile)
	return err
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.file != nil {
		return rl.file.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app.log")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: Some sample log data here\n",
			time.Now().Format(time.RFC3339), i)
		logger.Write([]byte(msg))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
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

type LogRotator struct {
	mu            sync.Mutex
	currentFile   *os.File
	basePath      string
	maxSize       int64
	currentSize   int64
	rotationCount int
}

func NewLogRotator(basePath string, maxSizeMB int) (*LogRotator, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024

	rotator := &LogRotator{
		basePath: basePath,
		maxSize:  maxSize,
	}

	if err := rotator.openCurrentFile(); err != nil {
		return nil, err
	}

	return rotator, nil
}

func (lr *LogRotator) openCurrentFile() error {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s.log", lr.basePath, timestamp)
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	lr.currentFile = file
	lr.currentSize = info.Size()
	lr.rotationCount = 0

	return nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lr.currentSize+int64(len(p)) > lr.maxSize {
		if err := lr.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := lr.currentFile.Write(p)
	if err == nil {
		lr.currentSize += int64(n)
	}
	return n, err
}

func (lr *LogRotator) rotate() error {
	if lr.currentFile != nil {
		lr.currentFile.Close()

		oldPath := lr.currentFile.Name()
		compressedPath := oldPath + ".gz"

		if err := compressFile(oldPath, compressedPath); err != nil {
			return err
		}

		if err := os.Remove(oldPath); err != nil {
			return err
		}

		lr.rotationCount++
	}

	return lr.openCurrentFile()
}

func compressFile(source, target string) error {
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(target)
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
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lr.currentFile != nil {
		return lr.currentFile.Close()
	}
	return nil
}

func main() {
	rotator, err := NewLogRotator("application", 10)
	if err != nil {
		fmt.Printf("Failed to create log rotator: %v\n", err)
		return
	}
	defer rotator.Close()

	for i := 0; i < 1000; i++ {
		logEntry := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
		if _, err := rotator.Write([]byte(logEntry)); err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}