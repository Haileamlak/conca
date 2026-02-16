package tools

import (
	"content-creator-agent/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// SearchTool defines the interface for research capabilities.
type SearchTool interface {
	Search(query string) ([]models.Trend, error)
}

// NewsAPISearch implements SearchTool using newsapi.org.
type NewsAPISearch struct {
	APIKey string
	client *http.Client
}

func NewNewsAPISearch(apiKey string) *NewsAPISearch {
	return &NewsAPISearch{
		APIKey: apiKey,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type newsAPIResponse struct {
	Status   string `json:"status"`
	Articles []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		URL         string `json:"url"`
	} `json:"articles"`
}

func (n *NewsAPISearch) Search(query string) ([]models.Trend, error) {
	if n.APIKey == "" {
		return nil, fmt.Errorf("newsapi key is required")
	}

	u := fmt.Sprintf("https://newsapi.org/v2/everything?q=%s&sortBy=publishedAt&pageSize=5&apiKey=%s", url.QueryEscape(query), n.APIKey)

	resp, err := n.client.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("newsapi error %d", resp.StatusCode)
	}

	var newsResp newsAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&newsResp); err != nil {
		return nil, err
	}

	var trends []models.Trend
	for _, art := range newsResp.Articles {
		trends = append(trends, models.Trend{
			Query:     query,
			Title:     art.Title,
			Snippet:   art.Description,
			URL:       art.URL,
			Timestamp: time.Now(),
		})
	}

	return trends, nil
}

// DuckDuckGoSearch implements SearchTool using HTML scraping.
type DuckDuckGoSearch struct {
	client *http.Client
}

func NewDuckDuckGoSearch() *DuckDuckGoSearch {
	return &DuckDuckGoSearch{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Search performs a query on DuckDuckGo and parses the HTML results.
// Note: This is fragile as it depends on DDG's HTML structure.
func (s *DuckDuckGoSearch) Search(query string) ([]models.Trend, error) {
	// DDG HTML endpoint
	baseURL := "https://html.duckduckgo.com/html/"

	vals := url.Values{}
	vals.Add("q", query)

	req, err := http.NewRequest("POST", baseURL, strings.NewReader(vals.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("search request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search returned status: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var trends []models.Trend

	// Parse results. DDG HTML usually has results in .result elements
	doc.Find(".result").Each(func(i int, s *goquery.Selection) {
		if i >= 5 { // Limit to top 5 results
			return
		}

		title := s.Find(".result__a").Text()
		link, _ := s.Find(".result__a").Attr("href")
		snippet := s.Find(".result__snippet").Text()

		if title != "" && link != "" {
			trends = append(trends, models.Trend{
				Query:     query,
				Title:     strings.TrimSpace(title),
				Snippet:   strings.TrimSpace(snippet),
				URL:       link,
				Timestamp: time.Now(),
			})
		}
	})

	return trends, nil
}

// ResilientSearch tries a primary tool and falls back to a backup on failure.
type ResilientSearch struct {
	Primary SearchTool
	Backup  SearchTool
}

func NewResilientSearch(primary, backup SearchTool) *ResilientSearch {
	return &ResilientSearch{Primary: primary, Backup: backup}
}

func (r *ResilientSearch) Search(query string) ([]models.Trend, error) {
	fmt.Printf("Searching with primary tool...\n")
	results, err := r.Primary.Search(query)
	if err == nil {
		return results, nil
	}

	fmt.Printf("Primary search failed (%v). Falling back to backup...\n", err)
	if r.Backup == nil {
		return nil, err
	}
	return r.Backup.Search(query)
}
