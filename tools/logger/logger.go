package logger

import (
	"fmt"
	"sync"
	"time"
)

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
}

type Buffer struct {
	entries []LogEntry
	maxSize int
	mu      sync.RWMutex
}

func NewBuffer(maxSize int) *Buffer {
	return &Buffer{
		entries: make([]LogEntry, 0, maxSize),
		maxSize: maxSize,
	}
}

func (b *Buffer) Log(level, message string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
	}

	if len(b.entries) >= b.maxSize {
		b.entries = b.entries[1:]
	}
	b.entries = append(b.entries, entry)

	// Also print to stdout for debugging
	fmt.Printf("[%s] %s: %s\n", entry.Timestamp.Format("15:04:05"), level, message)
}

func (b *Buffer) Info(format string, a ...interface{}) {
	b.Log("INFO", fmt.Sprintf(format, a...))
}

func (b *Buffer) Error(format string, a ...interface{}) {
	b.Log("ERROR", fmt.Sprintf(format, a...))
}

func (b *Buffer) Warn(format string, a ...interface{}) {
	b.Log("WARN", fmt.Sprintf(format, a...))
}

func (b *Buffer) GetEntries() []LogEntry {
	b.mu.RLock()
	defer b.mu.RUnlock()

	// Return a copy
	res := make([]LogEntry, len(b.entries))
	copy(res, b.entries)
	return res
}

var GlobalBuffer = NewBuffer(100)
