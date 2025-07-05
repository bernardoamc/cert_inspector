package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDailyLogger_WriteLine_CreatesFileAndWrites(t *testing.T) {
	tmpDir := t.TempDir()

	dl, err := NewDailyLogger(tmpDir, "testlog")
	if err != nil {
		t.Fatalf("Failed to create DailyLogger: %v", err)
	}
	defer dl.Close()

	err = dl.WriteLine("hello world")
	if err != nil {
		t.Fatalf("WriteLine failed: %v", err)
	}

	files, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to read log directory: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("Expected 1 log file, found %d", len(files))
	}

	name := files[0].Name()
	if !strings.HasPrefix(name, "testlog_") || !strings.HasSuffix(name, ".out") {
		t.Errorf("Unexpected filename: %q", name)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, name))
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	text := string(content)
	if text != "hello world\n" {
		t.Errorf("Log file content missing message, got: %q", text)
	}
}
