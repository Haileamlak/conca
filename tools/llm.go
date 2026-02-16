package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// LLMTool defines the interface for text generation.
type LLMTool interface {
	Generate(systemPrompt, userPrompt string) (string, error)
}

// GeminiClient implements LLMTool using Google's Gemini REST API.
type GeminiClient struct {
	APIKey string
	Model  string
	client *http.Client
}

func NewGeminiClient(apiKey, model string) *GeminiClient {
	if model == "" {
		model = "gemini-1.5-flash" // Flash is faster and cheaper/free
	}
	return &GeminiClient{
		APIKey: apiKey,
		Model:  model,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiRequest struct {
	Contents          []geminiContent `json:"contents"`
	SystemInstruction struct {
		Parts []geminiPart `json:"parts"`
	} `json:"system_instruction,omitempty"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func (g *GeminiClient) Generate(systemPrompt, userPrompt string) (string, error) {
	if g.APIKey == "" {
		return "", fmt.Errorf("gemini api key is required")
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", g.Model, g.APIKey)

	req := geminiRequest{}
	req.Contents = append(req.Contents, geminiContent{
		Parts: []geminiPart{
			{Text: userPrompt},
		},
	})

	if systemPrompt != "" {
		req.SystemInstruction.Parts = append(req.SystemInstruction.Parts, geminiPart{
			Text: systemPrompt,
		})
	}

	jsonBody, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	resp, err := g.client.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("gemini api error %d: %s", resp.StatusCode, string(body))
	}

	var gemResp geminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&gemResp); err != nil {
		return "", err
	}

	if len(gemResp.Candidates) == 0 || len(gemResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from gemini")
	}

	return gemResp.Candidates[0].Content.Parts[0].Text, nil
}

// OllamaClient implements LLMTool connecting to a local Ollama instance.
type OllamaClient struct {
	Model   string
	BaseURL string
	client  *http.Client
}

func NewOllamaClient(model string) *OllamaClient {
	if model == "" {
		model = "mistral"
	}
	return &OllamaClient{
		Model:   model,
		BaseURL: "http://localhost:11434/api/generate",
		client: &http.Client{
			Timeout: 120 * time.Second, // Long timeout for generation
		},
	}
}

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	System string `json:"system"`
	Stream bool   `json:"stream"`
}

type ollamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func (o *OllamaClient) Generate(systemPrompt, userPrompt string) (string, error) {
	reqBody := ollamaRequest{
		Model:  o.Model,
		Prompt: userPrompt,
		System: systemPrompt,
		Stream: false,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := o.client.Post(o.BaseURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("ollama request failed (is ollama running?): %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama error %d: %s", resp.StatusCode, string(body))
	}

	var ollamaResp ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return ollamaResp.Response, nil
}
