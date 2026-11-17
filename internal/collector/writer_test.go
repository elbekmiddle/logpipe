package collector

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWriter(t *testing.T) {
	path := filepath.Join(t.TempDir(), "logs.jsonl")

	writer, err := NewWriter(path)
	if err != nil {
		t.Fatalf("NewWriter failed: %v", err)
	}

	entry := Entry{
		Timestamp: time.Now(),
		Stream:    "STDOUT",
		Message:   "Server started",
	}

	if err := writer.Write(entry); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	var got Entry

	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	if got.Message != entry.Message {
		t.Errorf("expected %q, got %q", entry.Message, got.Message)
	}

	if got.Stream != entry.Stream {
		t.Errorf("expected %q, got %q", entry.Stream, got.Stream)
	}
}
