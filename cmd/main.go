package main

import (
	"content-creator-agent/agent"
	"content-creator-agent/memory"
	"content-creator-agent/models"
	"content-creator-agent/tools"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// CLI Flags
	configPath := flag.String("config", "config/tech_startup.json", "Path to brand config JSON")
	useDDG := flag.Bool("ddg", false, "Use DuckDuckGo instead of NewsAPI")
	flag.Parse()

	// 1. Load Brand Config
	configFile, err := os.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	var brand models.BrandProfile
	if err := json.Unmarshal(configFile, &brand); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	// 2. Initialize Tools
	var search tools.SearchTool
	if *useDDG {
		search = tools.NewDuckDuckGoSearch()
		fmt.Println("Using DuckDuckGo for research...")
	} else {
		apiKey := os.Getenv("NEWSAPI_KEY")
		if apiKey == "" {
			fmt.Println("WARNING: NEWSAPI_KEY not set. Falling back to DuckDuckGo.")
			search = tools.NewDuckDuckGoSearch()
		} else {
			search = tools.NewNewsAPISearch(apiKey)
			fmt.Println("Using NewsAPI for research...")
		}
	}

	geminiKey := os.Getenv("GEMINI_API_KEY")
	if geminiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable is required.")
	}
	llm := tools.NewGeminiClient(geminiKey, "gemini-1.5-flash")
	embedding := tools.NewGeminiEmbeddingClient(geminiKey, "text-embedding-004")

	// Multi-Social Client
	social := tools.NewMultiSocialClient()

	twitterKey := os.Getenv("TWITTER_API_KEY")
	if twitterKey != "" {
		twitterClient := tools.NewTwitterClient(
			twitterKey,
			os.Getenv("TWITTER_API_SECRET"),
			os.Getenv("TWITTER_ACCESS_TOKEN"),
			os.Getenv("TWITTER_ACCESS_SECRET"),
		)
		social.AddClient("twitter", twitterClient)
		fmt.Println("Twitter client integrated.")
	}

	linkedInToken := os.Getenv("LINKEDIN_ACCESS_TOKEN")
	if linkedInToken != "" {
		linkedInClient := tools.NewLinkedInClient(
			linkedInToken,
			os.Getenv("LINKEDIN_PERSON_URN"),
		)
		social.AddClient("linkedin", linkedInClient)
		fmt.Println("LinkedIn client integrated.")
	}

	// Fallback to Mock if no real clients added
	if len(social.Clients) == 0 {
		fmt.Println("Using MockSocialClient (no real credentials found).")
		social.AddClient("mock", &tools.MockSocialClient{Platform: "Mock"})
	}

	store := memory.NewFileStore("data")
	vector := memory.NewLocalVectorStore(filepath.Join("data", brand.ID, "vectors.json"))

	// 3. Initialize Agent
	creator := agent.NewAgent(brand, search, llm, social, store, vector, embedding)

	// 4. Run Loop
	if err := creator.Run(); err != nil {
		log.Fatalf("Agent run failed: %v", err)
	}
}
