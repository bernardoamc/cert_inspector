package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type DailyLogger struct {
	mu         sync.Mutex
	logDir     string
	baseName   string
	currentDay time.Time
	file       *os.File
}

func NewDailyLogger(logDir, baseName string) (*DailyLogger, error) {
	dl := &DailyLogger{logDir: logDir, baseName: baseName}
	return dl, dl.rotateIfNeeded()
}

func (dl *DailyLogger) rotateIfNeeded() error {
	dl.mu.Lock()
	defer dl.mu.Unlock()

	today := time.Now().Truncate(24 * time.Hour)
	if dl.file != nil && today.Equal(dl.currentDay) {
		return nil
	}

	if dl.file != nil {
		_ = dl.file.Close()
	}

	if err := os.MkdirAll(dl.logDir, 0755); err != nil {
		return fmt.Errorf("creating log dir: %w", err)
	}

	filename := fmt.Sprintf("%s_%04d_%02d_%02d.out", dl.baseName, today.Year(), today.Month(), today.Day())
	fullPath := filepath.Join(dl.logDir, filename)

	f, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("opening log file: %w", err)
	}

	dl.file = f
	dl.currentDay = today
	return nil
}

func (dl *DailyLogger) WriteLine(line string) error {
	if err := dl.rotateIfNeeded(); err != nil {
		return err
	}

	timestamped := fmt.Sprintf("[%s] %s\n", time.Now().Format(time.RFC3339), line)

	dl.mu.Lock()
	defer dl.mu.Unlock()
	_, err := dl.file.WriteString(timestamped)
	return err
}

func (dl *DailyLogger) Close() error {
	dl.mu.Lock()
	defer dl.mu.Unlock()
	if dl.file != nil {
		err := dl.file.Close()
		dl.file = nil
		return err
	}
	return nil
}
