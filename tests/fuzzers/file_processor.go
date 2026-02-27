package main

import (
    "bufio"
    "fmt"
    "os"
    "path/filepath"
    "sync"
)

type FileStats struct {
    Path    string
    Size    int64
    Lines   int
    Words   int
}

func processFile(path string, wg *sync.WaitGroup, results chan<- FileStats) {
    defer wg.Done()

    file, err := os.Open(path)
    if err != nil {
        fmt.Printf("Error opening file %s: %v\n", path, err)
        return
    }
    defer file.Close()

    var lines, words int
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines++
        text := scanner.Text()
        words += countWords(text)
    }

    info, err := file.Stat()
    if err != nil {
        fmt.Printf("Error getting file info %s: %v\n", path, err)
        return
    }

    results <- FileStats{
        Path:  path,
        Size:  info.Size(),
        Lines: lines,
        Words: words,
    }
}

func countWords(text string) int {
    inWord := false
    count := 0
    for _, ch := range text {
        if ch == ' ' || ch == '\t' || ch == '\n' {
            inWord = false
        } else if !inWord {
            inWord = true
            count++
        }
    }
    return count
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: file_processor <directory>")
        os.Exit(1)
    }

    root := os.Args[1]
    var wg sync.WaitGroup
    results := make(chan FileStats, 100)

    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.Mode().IsRegular() {
            return nil
        }

        wg.Add(1)
        go processFile(path, &wg, results)
        return nil
    })

    if err != nil {
        fmt.Printf("Error walking directory: %v\n", err)
        os.Exit(1)
    }

    go func() {
        wg.Wait()
        close(results)
    }()

    totalFiles := 0
    totalSize := int64(0)
    totalLines := 0
    totalWords := 0

    for stats := range results {
        totalFiles++
        totalSize += stats.Size
        totalLines += stats.Lines
        totalWords += stats.Words
        fmt.Printf("%s: %d bytes, %d lines, %d words\n",
            stats.Path, stats.Size, stats.Lines, stats.Words)
    }

    fmt.Printf("\nSummary:\n")
    fmt.Printf("Total files: %d\n", totalFiles)
    fmt.Printf("Total size: %d bytes\n", totalSize)
    fmt.Printf("Total lines: %d\n", totalLines)
    fmt.Printf("Total words: %d\n", totalWords)
}