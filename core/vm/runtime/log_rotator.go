
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

type RotatingWriter struct {
    filename   string
    current    *os.File
    size       int64
    mu         sync.Mutex
}

func NewRotatingWriter(filename string) (*RotatingWriter, error) {
    w := &RotatingWriter{filename: filename}
    if err := w.openFile(); err != nil {
        return nil, err
    }
    return w, nil
}

func (w *RotatingWriter) Write(p []byte) (int, error) {
    w.mu.Lock()
    defer w.mu.Unlock()

    if w.size+int64(len(p)) > maxFileSize {
        if err := w.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := w.current.Write(p)
    w.size += int64(n)
    return n, err
}

func (w *RotatingWriter) openFile() error {
    file, err := os.OpenFile(w.filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }
    w.current = file
    w.size = info.Size()
    return nil
}

func (w *RotatingWriter) rotate() error {
    if w.current != nil {
        w.current.Close()
    }

    timestamp := time.Now().Format("20060102-150405")
    backupName := fmt.Sprintf("%s.%s.gz", w.filename, timestamp)

    if err := compressFile(w.filename, backupName); err != nil {
        return err
    }

    if err := cleanupOldBackups(w.filename); err != nil {
        return err
    }

    os.Remove(w.filename)
    return w.openFile()
}

func compressFile(source, target string) error {
    src, err := os.Open(source)
    if err != nil {
        return err
    }
    defer src.Close()

    dst, err := os.Create(target)
    if err != nil {
        return err
    }
    defer dst.Close()

    gz := gzip.NewWriter(dst)
    defer gz.Close()

    _, err = io.Copy(gz, src)
    return err
}

func cleanupOldBackups(baseName string) error {
    pattern := fmt.Sprintf("%s.*.gz", filepath.Base(baseName))
    matches, err := filepath.Glob(filepath.Join(filepath.Dir(baseName), pattern))
    if err != nil {
        return err
    }

    if len(matches) > maxBackups {
        toDelete := matches[:len(matches)-maxBackups]
        for _, file := range toDelete {
            os.Remove(file)
        }
    }
    return nil
}

func (w *RotatingWriter) Close() error {
    w.mu.Lock()
    defer w.mu.Unlock()
    if w.current != nil {
        return w.current.Close()
    }
    return nil
}

func main() {
    writer, err := NewRotatingWriter("app.log")
    if err != nil {
        panic(err)
    }
    defer writer.Close()

    for i := 0; i < 100; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: This is a sample log message\n",
            time.Now().Format(time.RFC3339), i)
        writer.Write([]byte(logEntry))
        time.Sleep(100 * time.Millisecond)
    }
}