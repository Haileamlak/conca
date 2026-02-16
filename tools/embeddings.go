package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// EmbeddingTool defines the interface for generating vector embeddings.
type EmbeddingTool interface {
	Embed(text string) ([]float32, error)
}

// GeminiEmbeddingClient implements EmbeddingTool using Google Gemini API.
type GeminiEmbeddingClient struct {
	APIKey string
	Model  string
	client *http.Client
}

func NewGeminiEmbeddingClient(apiKey, model string) *GeminiEmbeddingClient {
	if model == "" {
		model = "text-embedding-004"
	}
	return &GeminiEmbeddingClient{
		APIKey: apiKey,
		Model:  model,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type geminiEmbeddingRequest struct {
	Model   string `json:"model"`
	Content struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"content"`
}

type geminiEmbeddingResponse struct {
	Embedding struct {
		Values []float32 `json:"values"`
	} `json:"embedding"`
}

func (g *GeminiEmbeddingClient) Embed(text string) ([]float32, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:embedContent?key=%s", g.Model, g.APIKey)

	reqBody := geminiEmbeddingRequest{}
	reqBody.Model = "models/" + g.Model
	reqBody.Content.Parts = append(reqBody.Content.Parts, struct {
		Text string `json:"text"`
	}{Text: text})

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := g.client.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("gemini embedding error %d: %s", resp.StatusCode, string(body))
	}

	var embedResp geminiEmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&embedResp); err != nil {
		return nil, err
	}

	return embedResp.Embedding.Values, nil
}
