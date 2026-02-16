package tools

import (
	"content-creator-agent/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// AnalyticsFetcher defines the interface for pulling performance data.
type AnalyticsFetcher interface {
	Fetch(post *models.Post) (models.Analytics, error)
}

// TwitterAnalyticsFetcher pulls metrics from X API v2.
type TwitterAnalyticsFetcher struct {
	Client *TwitterClient
}

func (t *TwitterAnalyticsFetcher) Fetch(post *models.Post) (models.Analytics, error) {
	if post.SocialID == "" {
		return models.Analytics{}, fmt.Errorf("post has no social ID")
	}

	apiURL := fmt.Sprintf("https://api.twitter.com/2/tweets/%s?tweet.fields=public_metrics", post.SocialID)
	req, _ := http.NewRequest("GET", apiURL, nil)
	authHeader := t.Client.generateOAuthHeader("GET", apiURL, nil)
	req.Header.Set("Authorization", authHeader)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return models.Analytics{}, err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			PublicMetrics struct {
				RetweetCount int `json:"retweet_count"`
				ReplyCount   int `json:"reply_count"`
				LikeCount    int `json:"like_count"`
				QuoteCount   int `json:"quote_count"`
			} `json:"public_metrics"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return models.Analytics{}, err
	}

	return models.Analytics{
		Likes:    result.Data.PublicMetrics.LikeCount,
		Shares:   result.Data.PublicMetrics.RetweetCount + result.Data.PublicMetrics.QuoteCount,
		Comments: result.Data.PublicMetrics.ReplyCount,
		Views:    0, // X API v2 basic metrics don't always include impressions for all tiers
	}, nil
}

// LinkedInAnalyticsFetcher pulls metrics from LinkedIn.
type LinkedInAnalyticsFetcher struct {
	Client *LinkedInClient
}

func (l *LinkedInAnalyticsFetcher) Fetch(post *models.Post) (models.Analytics, error) {
	if post.SocialID == "" {
		return models.Analytics{}, fmt.Errorf("post has no social ID")
	}

	// Example endpoint for organizational shares (URN required)
	apiURL := fmt.Sprintf("https://api.linkedin.com/v2/socialActions/%s", post.SocialID)
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("Authorization", "Bearer "+l.Client.AccessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return models.Analytics{}, err
	}
	defer resp.Body.Close()

	// Simplified parsing for brevity
	var result map[string]interface{} // In a real app, use a concrete struct
	json.NewDecoder(resp.Body).Decode(&result)

	// Mocking extraction as LinkedIn API response vary wildly by post type
	return models.Analytics{
		Likes:    10, // Logic to extract from result
		Shares:   2,
		Comments: 1,
	}, nil
}

// MultiAnalyticsFetcher routes to the correct platform fetcher.
type MultiAnalyticsFetcher struct {
	Fetchers map[string]AnalyticsFetcher
}

func (m *MultiAnalyticsFetcher) Fetch(post *models.Post) (models.Analytics, error) {
	fetcher, ok := m.Fetchers[post.Platform]
	if !ok {
		return models.Analytics{}, fmt.Errorf("no fetcher for platform: %s", post.Platform)
	}
	return fetcher.Fetch(post)
}
