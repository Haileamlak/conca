package tools

import (
	"content-creator-agent/models"
	"fmt"
	"strings"
	"time"
)

// SocialClient defines the interface for posting content.
type SocialClient interface {
	Post(post *models.Post) error
}

// MultiSocialClient routes posts to multiple social platforms.
type MultiSocialClient struct {
	Clients map[string]SocialClient
}

func NewMultiSocialClient() *MultiSocialClient {
	return &MultiSocialClient{
		Clients: make(map[string]SocialClient),
	}
}

func (m *MultiSocialClient) AddClient(platform string, client SocialClient) {
	m.Clients[platform] = client
}

func (m *MultiSocialClient) Post(post *models.Post) error {
	client, ok := m.Clients[post.Platform]
	if !ok {
		// Fallback to all clients if platform is "LinkedIn/X" or similar generic string
		if post.Platform == "LinkedIn/X" || post.Platform == "" {
			var errs []string
			for p, c := range m.Clients {
				if err := c.Post(post); err != nil {
					errs = append(errs, fmt.Sprintf("%s: %v", p, err))
				}
			}
			if len(errs) > 0 {
				return fmt.Errorf("multi-post failed: %s", strings.Join(errs, "; "))
			}
			return nil
		}
		return fmt.Errorf("no client configured for platform: %s", post.Platform)
	}
	return client.Post(post)
}

// MockSocialClient simulates posting by printing to console and updating status.
type MockSocialClient struct {
	Platform string
}

func (m *MockSocialClient) Post(post *models.Post) error {
	fmt.Printf("\n--- [MOCK POSTING to %s] ---\n", m.Platform)
	fmt.Printf("Brand:    %s\n", post.BrandID)
	fmt.Printf("Content:  %s\n", post.Content)
	fmt.Printf("---------------------------\n")

	post.Status = models.StatusPublished
	post.UpdatedAt = time.Now()

	// Simulate some initial analytics for the "feedback loop" later
	post.Analytics = models.Analytics{
		Views:    100,
		Likes:    10,
		Shares:   2,
		Comments: 1,
	}

	return nil
}

type InstagramClient struct{ MockSocialClient }
type TikTokClient struct{ MockSocialClient }
type ThreadsClient struct{ MockSocialClient }

func NewInstagramClient() *InstagramClient {
	return &InstagramClient{MockSocialClient{Platform: "instagram"}}
}

func NewTikTokClient() *TikTokClient {
	return &TikTokClient{MockSocialClient{Platform: "tiktok"}}
}

func NewThreadsClient() *ThreadsClient {
	return &ThreadsClient{MockSocialClient{Platform: "threads"}}
}
