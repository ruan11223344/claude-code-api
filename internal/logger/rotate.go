package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// RotatingFileWriter implements io.Writer with daily rotation
type RotatingFileWriter struct {
	baseDir     string
	currentFile *os.File
	currentDate string
	maxFiles    int
}

// NewRotatingFileWriter creates a new rotating file writer
func NewRotatingFileWriter(baseDir string, maxFiles int) (*RotatingFileWriter, error) {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}
	
	writer := &RotatingFileWriter{
		baseDir:  baseDir,
		maxFiles: maxFiles,
	}
	
	// Open initial file
	if err := writer.rotate(); err != nil {
		return nil, err
	}
	
	// Clean old files
	writer.cleanOldFiles()
	
	return writer, nil
}

// Write implements io.Writer
func (w *RotatingFileWriter) Write(p []byte) (n int, err error) {
	// Check if we need to rotate
	currentDate := time.Now().Format("2006-01-02")
	if currentDate != w.currentDate {
		if err := w.rotate(); err != nil {
			return 0, err
		}
		w.cleanOldFiles()
	}
	
	if w.currentFile == nil {
		return 0, fmt.Errorf("log file is not open")
	}
	
	return w.currentFile.Write(p)
}

// Close closes the current file
func (w *RotatingFileWriter) Close() error {
	if w.currentFile != nil {
		return w.currentFile.Close()
	}
	return nil
}

// rotate opens a new log file for the current date
func (w *RotatingFileWriter) rotate() error {
	// Close current file if open
	if w.currentFile != nil {
		w.currentFile.Close()
	}
	
	// Update current date
	w.currentDate = time.Now().Format("2006-01-02")
	
	// Open new file
	fileName := fmt.Sprintf("server_%s.log", w.currentDate)
	filePath := filepath.Join(w.baseDir, fileName)
	
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	
	w.currentFile = file
	return nil
}

// cleanOldFiles removes old log files beyond maxFiles limit
func (w *RotatingFileWriter) cleanOldFiles() {
	if w.maxFiles <= 0 {
		return // No limit
	}
	
	// List all log files
	files, err := filepath.Glob(filepath.Join(w.baseDir, "server_*.log"))
	if err != nil {
		return
	}
	
	// Sort files by name (which includes date)
	sort.Strings(files)
	
	// Remove old files if we exceed the limit
	if len(files) > w.maxFiles {
		for _, file := range files[:len(files)-w.maxFiles] {
			os.Remove(file)
		}
	}
}

// MultiWriter creates a writer that writes to both stdout and rotating file
func MultiWriter(stdout io.Writer, baseDir string, maxFiles int) (io.Writer, error) {
	fileWriter, err := NewRotatingFileWriter(baseDir, maxFiles)
	if err != nil {
		return stdout, err
	}
	
	return io.MultiWriter(stdout, fileWriter), nil
}