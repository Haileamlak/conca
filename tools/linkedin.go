package tools

import (
	"bytes"
	"content-creator-agent/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// LinkedInClient handles posting content to LinkedIn.
type LinkedInClient struct {
	AccessToken string
	PersonURN   string // The profile URN (e.g., "urn:li:person:12345")
}

func NewLinkedInClient(accessToken, personURN string) *LinkedInClient {
	return &LinkedInClient{
		AccessToken: accessToken,
		PersonURN:   personURN,
	}
}

// Post sends a share via the LinkedIn API.
func (l *LinkedInClient) Post(post *models.Post) error {
	apiURL := "https://api.linkedin.com/v2/ugcPosts"

	// Construct LinkedIn UGC Post payload
	payload := map[string]interface{}{
		"author":         l.PersonURN,
		"lifecycleState": "PUBLISHED",
		"specificContent": map[string]interface{}{
			"com.linkedin.ugc.ShareContent": map[string]interface{}{
				"shareCommentary": map[string]interface{}{
					"text": post.Content,
				},
				"shareMediaCategory": "NONE",
			},
		},
		"visibility": map[string]interface{}{
			"com.linkedin.ugc.MemberNetworkVisibility": "PUBLIC",
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal linkedin payload: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+l.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Restli-Protocol-Version", "2.0.0")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request to LinkedIn failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("LinkedIn API returned error (status %d): %s", resp.StatusCode, string(body))
	}

	var liResp struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&liResp); err == nil {
		post.SocialID = liResp.ID
	}

	post.Status = models.StatusPublished
	post.UpdatedAt = time.Now()
	return nil
}
