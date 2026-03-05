package main

import (
	"content-creator-agent/api"
	"content-creator-agent/memory"
	"content-creator-agent/scheduler"
	"content-creator-agent/tools"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	port := flag.String("port", "8000", "Port to listen on")
	dataDir := flag.String("data", "data", "Data directory")
	flag.Parse()

	godotenv.Load()

	// --- Initialize shared tools ---
	geminiKey := os.Getenv("GEMINI_API_KEY")
	if geminiKey == "" {
		log.Fatal("GEMINI_API_KEY is required")
	}

	// Search
	var search tools.SearchTool
	ddg := tools.NewDuckDuckGoSearch()
	newsKey := os.Getenv("NEWSAPI_KEY")
	if newsKey != "" {
		var primary tools.SearchTool
		if strings.HasPrefix(newsKey, "pub_") {
			primary = tools.NewNewsDataSearch(newsKey)
		} else {
			primary = tools.NewNewsAPISearch(newsKey)
		}
		search = tools.NewResilientSearch(primary, ddg)
	} else {
		search = ddg
	}

	llm := tools.NewGeminiClient(geminiKey, "gemini-3-flash-preview")
	embedding := tools.NewGeminiEmbeddingClient(geminiKey, "gemini-embedding-001")

	// Social
	social := tools.NewMultiSocialClient()
	analytics := &tools.MultiAnalyticsFetcher{Fetchers: make(map[string]tools.AnalyticsFetcher)}

	twitterKey := os.Getenv("TWITTER_API_KEY")
	if twitterKey != "" {
		tc := tools.NewTwitterClient(twitterKey, os.Getenv("TWITTER_API_SECRET"), os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_SECRET"))
		social.AddClient("twitter", tc)
		analytics.Fetchers["twitter"] = &tools.TwitterAnalyticsFetcher{Client: tc}
	}

	liToken := os.Getenv("LINKEDIN_ACCESS_TOKEN")
	if liToken != "" {
		lc := tools.NewLinkedInClient(liToken, os.Getenv("LINKEDIN_PERSON_URN"))
		social.AddClient("linkedin", lc)
		analytics.Fetchers["linkedin"] = &tools.LinkedInAnalyticsFetcher{Client: lc}
	}

	// Meta Platforms (Stubs)
	social.AddClient("instagram", tools.NewInstagramClient())
	social.AddClient("tiktok", tools.NewTikTokClient())
	social.AddClient("threads", tools.NewThreadsClient())

	if len(social.Clients) == 0 {
		social.AddClient("mock", &tools.MockSocialClient{Platform: "Mock"})
	}

	// --- DatabaseSelection & Setup ---
	var store *memory.MongoStore
	mongoURI := os.Getenv("MONGODB_URI")
	mongoDB := os.Getenv("MONGODB_DB")
	if mongoDB == "" {
		mongoDB = "conca"
	}

	if mongoURI == "" {
		fmt.Println("WARNING: MONGODB_URI not set. Attempting to connect to localhost:27017")
		mongoURI = "mongodb://localhost:27017"
	}

	mStore, err := memory.NewMongoStore(mongoURI, mongoDB)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mStore.Close()
	store = mStore
	fmt.Printf("🍃 Using MongoDB database (%s) for storage.\n", mongoDB)

	// --- JWT Secret ---
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-change-in-production"
		fmt.Println("WARNING: Using default JWT_SECRET. Set JWT_SECRET env for production.")
	}

	// --- Job Queue & Workers ---
	queue := scheduler.NewMongoQueue(mStore.Database())

	factory := scheduler.DefaultAgentFactory(store, search, llm, social, embedding, analytics, mStore.Database())
	worker := scheduler.NewWorker(queue, factory)
	go worker.Start(context.Background())

	sched := scheduler.NewScheduler(store, queue)
	go sched.Start()

	// --- Build server ---
	handlers := &api.Handlers{
		Store:     store,
		Queue:     queue,
		JWTSecret: jwtSecret,
		Search:    search,
		LLM:       llm,
		Social:    social,
		Embedding: embedding,
		Analytics: analytics,
		DataDir:   *dataDir,
	}

	server := api.NewServer(handlers, jwtSecret, *port)

	fmt.Printf("🚀 Conca API running on http://localhost:%s\n", *port)
	if err := server.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
