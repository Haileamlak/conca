package memory

import (
	"content-creator-agent/models"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Store defines the interface for long-term persistence.
type Store interface {
	SavePost(post models.Post) error
	GetHistory(brandID string) ([]models.Post, error)
	GetAnalytics(brandID string) ([]models.Analytics, error)
	UpdateAnalytics(brandID string, postID string, analytics models.Analytics) error
}

// FileStore implements Store using JSON files on disk.
type FileStore struct {
	BaseDir string
	mu      sync.Mutex
}

func NewFileStore(baseDir string) *FileStore {
	if baseDir == "" {
		baseDir = "data"
	}
	return &FileStore{BaseDir: baseDir}
}

func (f *FileStore) brandPath(brandID string) string {
	return filepath.Join(f.BaseDir, brandID)
}

func (f *FileStore) SavePost(post models.Post) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	path := f.brandPath(post.BrandID)
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create brand dir: %w", err)
	}

	historyPath := filepath.Join(path, "history.json")
	var history []models.Post

	// Read existing history
	data, err := os.ReadFile(historyPath)
	if err == nil {
		json.Unmarshal(data, &history)
	}

	// Append new post
	history = append(history, post)

	// Save back
	updatedData, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(historyPath, updatedData, 0644)
}

func (f *FileStore) GetHistory(brandID string) ([]models.Post, error) {
	historyPath := filepath.Join(f.brandPath(brandID), "history.json")
	data, err := os.ReadFile(historyPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.Post{}, nil
		}
		return nil, err
	}

	var history []models.Post
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, err
	}
	return history, nil
}

// GetAnalytics is a placeholder for retrieving aggregated analytics.
func (f *FileStore) GetAnalytics(brandID string) ([]models.Analytics, error) {
	history, err := f.GetHistory(brandID)
	if err != nil {
		return nil, err
	}

	var analytics []models.Analytics
	for _, p := range history {
		analytics = append(analytics, p.Analytics)
	}
	return analytics, nil
}
func (f *FileStore) UpdateAnalytics(brandID string, postID string, analytics models.Analytics) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	historyPath := filepath.Join(f.brandPath(brandID), "history.json")
	data, err := os.ReadFile(historyPath)
	if err != nil {
		return err
	}

	var history []models.Post
	if err := json.Unmarshal(data, &history); err != nil {
		return err
	}

	found := false
	for i := range history {
		if history[i].ID == postID {
			history[i].Analytics = analytics
			history[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("post %s not found in history", postID)
	}

	updatedData, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(historyPath, updatedData, 0644)
}
