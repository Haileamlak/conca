package tools

import (
	"bytes"
	"content-creator-agent/models"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// TwitterClient handles posting content to X (formerly Twitter) using API v2.
type TwitterClient struct {
	APIKey            string
	APIKeySecret      string
	AccessToken       string
	AccessTokenSecret string
}

func NewTwitterClient(apiKey, apiSecret, accessToken, accessSecret string) *TwitterClient {
	return &TwitterClient{
		APIKey:            apiKey,
		APIKeySecret:      apiSecret,
		AccessToken:       accessToken,
		AccessTokenSecret: accessSecret,
	}
}

// Post sends a tweet via the X API v2.
func (t *TwitterClient) Post(post *models.Post) error {
	apiURL := "https://api.twitter.com/2/tweets"

	payload := map[string]string{
		"text": post.Content,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal tweet payload: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Sign request with OAuth 1.0a (simplified implementation for the artifact)
	authHeader := t.generateOAuthHeader("POST", apiURL, nil)
	req.Header.Set("Authorization", authHeader)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request to X failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("X API returned error (status %d): %s", resp.StatusCode, string(body))
	}

	post.Status = models.StatusPublished
	post.UpdatedAt = time.Now()
	return nil
}

// generateOAuthHeader implements a basic OAuth 1.0a header generator.
func (t *TwitterClient) generateOAuthHeader(method, apiURL string, params map[string]string) string {
	nonce := fmt.Sprintf("%d", time.Now().UnixNano())
	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	oauthParams := map[string]string{
		"oauth_consumer_key":     t.APIKey,
		"oauth_nonce":            nonce,
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        timestamp,
		"oauth_token":            t.AccessToken,
		"oauth_version":          "1.0",
	}

	// In a real implementation, we'd include 'params' and build the signature base string.
	// For this code artifact, we'll use generic signing logic or note it as a critical piece.
	// Note: Post body params are NOT included in the signature for JSON payloads in v2.

	signature := t.sign(method, apiURL, oauthParams)
	oauthParams["oauth_signature"] = signature

	var parts []string
	for k, v := range oauthParams {
		parts = append(parts, fmt.Sprintf("%s=\"%s\"", url.QueryEscape(k), url.QueryEscape(v)))
	}
	sort.Strings(parts)

	return "OAuth " + strings.Join(parts, ", ")
}

func (t *TwitterClient) sign(method, apiURL string, params map[string]string) string {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var paramParts []string
	for _, k := range keys {
		paramParts = append(paramParts, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(params[k])))
	}
	paramString := strings.Join(paramParts, "&")

	signatureBase := fmt.Sprintf("%s&%s&%s",
		strings.ToUpper(method),
		url.QueryEscape(apiURL),
		url.QueryEscape(paramString))

	signingKey := fmt.Sprintf("%s&%s",
		url.QueryEscape(t.APIKeySecret),
		url.QueryEscape(t.AccessTokenSecret))

	hash := hmac.New(sha1.New, []byte(signingKey))
	hash.Write([]byte(signatureBase))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}
