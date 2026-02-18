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
	// Posts
	SavePost(post models.Post) error
	GetHistory(brandID string) ([]models.Post, error)
	GetGlobalHistory(userID string, limit int) ([]models.Post, error)
	GetAnalytics(brandID string) ([]models.Analytics, error)
	UpdateAnalytics(brandID string, postID string, analytics models.Analytics) error
	GetGlobalAnalytics(userID string) (models.GlobalAnalytics, error)
	GetBrandPerformance(userID string) ([]models.BrandPerformance, error)

	// Brands
	SaveBrand(brand models.BrandProfile, userID string) error
	GetBrand(id string) (models.BrandProfile, string, error) // Returns brand + userID
	ListBrands(userID string) ([]models.BrandProfile, error)
	ListAllBrands() ([]models.BrandProfile, error)
	DeleteBrand(id string) error

	// Calendar & Approval
	SaveScheduledPost(post models.ScheduledPost) error
	GetScheduledPosts(brandID string) ([]models.ScheduledPost, error)
	UpdateScheduledPostStatus(postID string, status models.PostStatus) error
	UpdateScheduledPost(postID string, topic, content string) error
	GetPendingScheduledPosts() ([]models.ScheduledPost, error) // For the scheduler to publish

	// User management
	CreateUser(email, passwordHash string) (string, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
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

func (f *FileStore) GetGlobalHistory(userID string, limit int) ([]models.Post, error) {
	brands, err := f.ListBrands(userID)
	if err != nil {
		return nil, err
	}

	var allPosts []models.Post
	for _, b := range brands {
		posts, _ := f.GetHistory(b.ID)
		allPosts = append(allPosts, posts...)
	}

	// Sort by created_at desc
	// (Actually for simplicity I'll just return the last N)
	if limit > 0 && len(allPosts) > limit {
		allPosts = allPosts[len(allPosts)-limit:]
	}
	return allPosts, nil
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

func (f *FileStore) GetGlobalAnalytics(userID string) (models.GlobalAnalytics, error) {
	brands, err := f.ListBrands(userID)
	if err != nil {
		return models.GlobalAnalytics{}, err
	}

	var global models.GlobalAnalytics
	for _, b := range brands {
		analytics, _ := f.GetAnalytics(b.ID)
		for _, a := range analytics {
			global.TotalImpressions += a.Views
			global.TotalLikes += a.Likes
			global.TotalShares += a.Shares
			global.TotalComments += a.Comments
		}
	}
	return global, nil
}

func (f *FileStore) GetBrandPerformance(userID string) ([]models.BrandPerformance, error) {
	brands, err := f.ListBrands(userID)
	if err != nil {
		return nil, err
	}

	var results []models.BrandPerformance
	for _, b := range brands {
		posts, _ := f.GetHistory(b.ID)
		if len(posts) == 0 {
			continue
		}

		perf := models.BrandPerformance{
			BrandID:   b.ID,
			BrandName: b.Name,
			PostCount: len(posts),
		}

		totalLikes := 0
		totalShares := 0
		for _, p := range posts {
			totalLikes += p.Analytics.Likes
			totalShares += p.Analytics.Shares
		}

		perf.AvgLikes = float64(totalLikes) / float64(len(posts))
		perf.AvgShares = float64(totalShares) / float64(len(posts))

		// Simple score calculation: avg likes * 2 + avg shares * 5
		score := int(perf.AvgLikes*2 + perf.AvgShares*5)
		if score > 100 {
			score = 100
		}
		perf.Score = score

		results = append(results, perf)
	}
	return results, nil
}

// --- Brand Management (FileStore Impl) ---

func (f *FileStore) SaveBrand(brand models.BrandProfile, userID string) error {
	brand.UserID = userID
	path := filepath.Join(f.BaseDir, brand.ID)
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	configPath := filepath.Join(path, "config.json")
	data, _ := json.MarshalIndent(brand, "", "  ")
	return os.WriteFile(configPath, data, 0644)
}

func (f *FileStore) GetBrand(id string) (models.BrandProfile, string, error) {
	configPath := filepath.Join(f.BaseDir, id, "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return models.BrandProfile{}, "", err
	}
	var brand models.BrandProfile
	if err := json.Unmarshal(data, &brand); err != nil {
		return models.BrandProfile{}, "", err
	}
	return brand, brand.UserID, nil
}

func (f *FileStore) ListBrands(userID string) ([]models.BrandProfile, error) {
	all, err := f.ListAllBrands()
	if err != nil {
		return nil, err
	}
	var filtered []models.BrandProfile
	for _, b := range all {
		if b.UserID == userID {
			filtered = append(filtered, b)
		}
	}
	return filtered, nil
}

func (f *FileStore) ListAllBrands() ([]models.BrandProfile, error) {
	entries, err := os.ReadDir(f.BaseDir)
	if err != nil {
		return nil, err
	}
	var brands []models.BrandProfile
	for _, entry := range entries {
		if entry.IsDir() {
			brand, _, err := f.GetBrand(entry.Name())
			if err == nil {
				brands = append(brands, brand)
			}
		}
	}
	return brands, nil
}

func (f *FileStore) DeleteBrand(id string) error {
	path := filepath.Join(f.BaseDir, id)
	return os.RemoveAll(path)
}

// --- Calendar & Approval (FileStore Impl) ---

func (f *FileStore) SaveScheduledPost(post models.ScheduledPost) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	path := f.brandPath(post.BrandID)
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	calendarPath := filepath.Join(path, "calendar.json")
	var calendar []models.ScheduledPost

	data, err := os.ReadFile(calendarPath)
	if err == nil {
		json.Unmarshal(data, &calendar)
	}

	calendar = append(calendar, post)
	updatedData, _ := json.MarshalIndent(calendar, "", "  ")
	return os.WriteFile(calendarPath, updatedData, 0644)
}

func (f *FileStore) GetScheduledPosts(brandID string) ([]models.ScheduledPost, error) {
	calendarPath := filepath.Join(f.brandPath(brandID), "calendar.json")
	data, err := os.ReadFile(calendarPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.ScheduledPost{}, nil
		}
		return nil, err
	}

	var calendar []models.ScheduledPost
	if err := json.Unmarshal(data, &calendar); err != nil {
		return nil, err
	}
	return calendar, nil
}

func (f *FileStore) UpdateScheduledPostStatus(postID string, status models.PostStatus) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	brands, _ := f.ListAllBrands()
	for _, b := range brands {
		calendarPath := filepath.Join(f.brandPath(b.ID), "calendar.json")
		data, err := os.ReadFile(calendarPath)
		if err != nil {
			continue
		}

		var calendar []models.ScheduledPost
		json.Unmarshal(data, &calendar)

		found := false
		for i := range calendar {
			if calendar[i].ID == postID {
				calendar[i].Status = status
				calendar[i].UpdatedAt = time.Now()
				found = true
				break
			}
		}

		if found {
			updatedData, _ := json.MarshalIndent(calendar, "", "  ")
			return os.WriteFile(calendarPath, updatedData, 0644)
		}
	}
	return fmt.Errorf("scheduled post %s not found", postID)
}

func (f *FileStore) UpdateScheduledPost(postID string, topic, content string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	brands, _ := f.ListAllBrands()
	for _, b := range brands {
		calendarPath := filepath.Join(f.brandPath(b.ID), "calendar.json")
		data, err := os.ReadFile(calendarPath)
		if err != nil {
			continue
		}

		var calendar []models.ScheduledPost
		json.Unmarshal(data, &calendar)

		found := false
		for i := range calendar {
			if calendar[i].ID == postID {
				calendar[i].Topic = topic
				calendar[i].Content = content
				calendar[i].UpdatedAt = time.Now()
				found = true
				break
			}
		}

		if found {
			updatedData, _ := json.MarshalIndent(calendar, "", "  ")
			return os.WriteFile(calendarPath, updatedData, 0644)
		}
	}
	return fmt.Errorf("scheduled post %s not found", postID)
}

func (f *FileStore) GetPendingScheduledPosts() ([]models.ScheduledPost, error) {
	brands, _ := f.ListAllBrands()
	var pending []models.ScheduledPost
	for _, b := range brands {
		posts, _ := f.GetScheduledPosts(b.ID)
		for _, p := range posts {
			if p.Status == models.StatusScheduled && p.ScheduledAt.Before(time.Now()) {
				pending = append(pending, p)
			}
		}
	}
	return pending, nil
}

// --- User Management (FileStore Impl) ---

func (f *FileStore) CreateUser(email, passwordHash string) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	usersPath := filepath.Join(f.BaseDir, "users.json")
	var users []models.User

	data, err := os.ReadFile(usersPath)
	if err == nil {
		json.Unmarshal(data, &users)
	}

	for _, u := range users {
		if u.Email == email {
			return "", fmt.Errorf("user already exists")
		}
	}

	userID := fmt.Sprintf("u-%d", time.Now().Unix())
	users = append(users, models.User{
		ID:           userID,
		Email:        email,
		PasswordHash: passwordHash,
	})

	updatedData, _ := json.MarshalIndent(users, "", "  ")
	os.WriteFile(usersPath, updatedData, 0644)

	return userID, nil
}

func (f *FileStore) GetUserByEmail(email string) (*models.User, error) {
	usersPath := filepath.Join(f.BaseDir, "users.json")
	data, err := os.ReadFile(usersPath)
	if err != nil {
		return nil, err
	}

	var users []models.User
	json.Unmarshal(data, &users)

	for _, u := range users {
		if u.Email == email {
			return &u, nil
		}
	}

	return nil, fmt.Errorf("user not found")
}

func (f *FileStore) GetUserByID(id string) (*models.User, error) {
	usersPath := filepath.Join(f.BaseDir, "users.json")
	data, err := os.ReadFile(usersPath)
	if err != nil {
		return nil, err
	}

	var users []models.User
	json.Unmarshal(data, &users)

	for _, u := range users {
		if u.ID == id {
			return &u, nil
		}
	}

	return nil, fmt.Errorf("user not found")
}
