package collector

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Writer struct {
	mu      sync.Mutex
	file    *os.File
	encoder *json.Encoder
}

func NewWriter(path string) (*Writer, error) {

	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create log directory: %w", err)
	}

	file, err := os.OpenFile(
		path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	return &Writer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (w *Writer) Write(entry Entry) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if err := w.encoder.Encode(entry); err != nil {
		return fmt.Errorf("encode log entry: %w", err)
	}

	return nil
}

func (w *Writer) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.file.Close()
}
