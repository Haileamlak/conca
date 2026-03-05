package memory

import (
	"content-creator-agent/models"
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
