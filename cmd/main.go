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
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// CLI Flags
	configPath := flag.String("config", "config/tech_startup.json", "Path to brand config JSON")
	useDDG := flag.Bool("ddg", false, "Use DuckDuckGo instead of NewsAPI")
	syncOnly := flag.Bool("sync", false, "Only sync analytics for past posts")
	daemon := flag.Bool("daemon", false, "Run in autonomous daemon mode")
	interval := flag.Duration("interval", 4*time.Hour, "Interval between cycles in daemon mode (e.g. 1h, 30m)")
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
	godotenv.Load()
	var search tools.SearchTool
	ddg := tools.NewDuckDuckGoSearch()

	apiKey := os.Getenv("NEWSAPI_KEY")
	if *useDDG || apiKey == "" {
		if apiKey == "" && !*useDDG {
			fmt.Println("WARNING: NEWSAPI_KEY not set. Using DuckDuckGo.")
		}
		search = ddg
	} else {
		var primary tools.SearchTool
		if strings.HasPrefix(apiKey, "pub_") {
			primary = tools.NewNewsDataSearch(apiKey)
		} else {
			primary = tools.NewNewsAPISearch(apiKey)
		}
		search = tools.NewResilientSearch(primary, ddg)
	}

	geminiKey := os.Getenv("GEMINI_API_KEY")
	if geminiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable is required.")
	}
	llm := tools.NewGeminiClient(geminiKey, "gemini-2.5-flash")
	embedding := tools.NewGeminiEmbeddingClient(geminiKey, "gemini-embedding-001")

	// Multi-Social Client
	social := tools.NewMultiSocialClient()
	analytics := &tools.MultiAnalyticsFetcher{Fetchers: make(map[string]tools.AnalyticsFetcher)}

	twitterKey := os.Getenv("TWITTER_API_KEY")
	if twitterKey != "" {
		tc := tools.NewTwitterClient(twitterKey, os.Getenv("TWITTER_API_SECRET"), os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_SECRET"))
		social.AddClient("twitter", tc)
		analytics.Fetchers["twitter"] = &tools.TwitterAnalyticsFetcher{Client: tc}
	}

	linkedInToken := os.Getenv("LINKEDIN_ACCESS_TOKEN")
	if linkedInToken != "" {
		lc := tools.NewLinkedInClient(linkedInToken, os.Getenv("LINKEDIN_PERSON_URN"))
		social.AddClient("linkedin", lc)
		analytics.Fetchers["linkedin"] = &tools.LinkedInAnalyticsFetcher{Client: lc}
	}

	if len(social.Clients) == 0 {
		social.AddClient("mock", &tools.MockSocialClient{Platform: "Mock"})
	}

	// --- Database Selection ---
	var store memory.Store
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL != "" {
		pgStore, err := memory.NewPostgresStore(dbURL)
		if err != nil {
			log.Fatalf("Failed to connect to PostgreSQL: %v", err)
		}
		defer pgStore.Close()
		store = pgStore
		fmt.Println("✅ CLI using PostgreSQL database.")
	} else {
		store = memory.NewFileStore("data")
		fmt.Println("📁 CLI using local JSON files.")
	}

	vector := memory.NewLocalVectorStore(filepath.Join("data", brand.ID, "vectors.json"))

	// 3. Initialize Agent
	creator := agent.NewAgent(brand, search, llm, social, store, vector, embedding, analytics)

	// 4. Run Logic
	if *syncOnly {
		fmt.Println("Running analytics sync...")
		if err := creator.SyncAnalytics(); err != nil {
			log.Fatalf("Sync failed: %v", err)
		}
		return
	}

	if *daemon {
		fmt.Printf("Starting agent in daemon mode with interval %s...\n", *interval)
		creator.Start(*interval)
		return
	}

	if err := creator.Run(); err != nil {
		log.Fatalf("Agent run failed: %v", err)
	}
}
