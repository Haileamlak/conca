package api

import (
	"net/http"
	"os"
)

type IntegrationStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"` // "connected" or "disconnected"
	Config string `json:"config,omitempty"`
}

type SystemStatus struct {
	Integrations []IntegrationStatus `json:"integrations"`
	Config       map[string]string   `json:"config"`
}

func (h *Handlers) GetSystemStatus(w http.ResponseWriter, r *http.Request) {
	status := SystemStatus{
		Integrations: []IntegrationStatus{
			{
				Name:   "MongoDB",
				Status: "connected", // If we're here, DB is working
				Config: os.Getenv("MONGODB_DB"),
			},
			{
				Name:   "Gemini API",
				Status: checkEnv("GEMINI_API_KEY"),
			},
			{
				Name:   "NewsAPI",
				Status: checkEnv("NEWSAPI_KEY"),
			},
			{
				Name:   "Twitter / X API",
				Status: checkEnv("TWITTER_API_KEY"),
			},
			{
				Name:   "LinkedIn API",
				Status: checkEnv("LINKEDIN_ACCESS_TOKEN"),
			},
		},
		Config: map[string]string{
			"port": "8000",
			"env":  "development",
		},
	}

	JSON(w, http.StatusOK, status)
}

func checkEnv(key string) string {
	if os.Getenv(key) != "" {
		return "connected"
	}
	return "disconnected"
}
