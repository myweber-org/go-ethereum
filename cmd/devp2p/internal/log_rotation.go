package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

const maxLogSize = 1024 * 1024 // 1MB
const backupCount = 5

type RotatingWriter struct {
	currentSize int64
	basePath    string
	file        *os.File
}

func NewRotatingWriter(path string) (*RotatingWriter, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	return &RotatingWriter{
		currentSize: info.Size(),
		basePath:    path,
		file:        file,
	}, nil
}

func (w *RotatingWriter) Write(p []byte) (int, error) {
	if w.currentSize+int64(len(p)) > maxLogSize {
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := w.file.Write(p)
	if err == nil {
		w.currentSize += int64(n)
	}
	return n, err
}

func (w *RotatingWriter) rotate() error {
	if err := w.file.Close(); err != nil {
		return err
	}

	for i := backupCount - 1; i >= 0; i-- {
		oldPath := w.backupPath(i)
		newPath := w.backupPath(i + 1)

		if _, err := os.Stat(oldPath); err == nil {
			if err := os.Rename(oldPath, newPath); err != nil {
				return err
			}
		}
	}

	if err := os.Rename(w.basePath, w.backupPath(0)); err != nil && !os.IsNotExist(err) {
		return err
	}

	file, err := os.OpenFile(w.basePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	w.file = file
	w.currentSize = 0
	return nil
}

func (w *RotatingWriter) backupPath(index int) string {
	if index == 0 {
		return w.basePath
	}
	return w.basePath + "." + string(rune('0'+index))
}

func (w *RotatingWriter) Close() error {
	return w.file.Close()
}

func main() {
	writer, err := NewRotatingWriter("logs/app.log")
	if err != nil {
		log.Fatal(err)
	}
	defer writer.Close()

	log.SetOutput(io.MultiWriter(os.Stdout, writer))

	for i := 0; i < 10000; i++ {
		log.Printf("Log entry number %d", i)
	}
}