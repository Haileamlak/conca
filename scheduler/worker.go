package scheduler

import (
	"content-creator-agent/agent"
	"content-creator-agent/memory"
	"content-creator-agent/models"
	"content-creator-agent/tools"
	"content-creator-agent/tools/logger"
	"context"
	"fmt"
	"path/filepath"
	"time"
)

// AgentFactory creates a new agent for a specific brand.
type AgentFactory func(brandID string) (*agent.Agent, error)

type Worker struct {
	Queue        Queue
	AgentFactory AgentFactory
	Quit         chan bool
}

func NewWorker(q Queue, factory AgentFactory) *Worker {
	return &Worker{
		Queue:        q,
		AgentFactory: factory,
		Quit:         make(chan bool),
	}
}

// Start runs the worker loop.
func (w *Worker) Start(ctx context.Context) {
	logger.GlobalBuffer.Info("👷 Worker started. Waiting for jobs...")
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.GlobalBuffer.Info("Worker shutting down...")
			return
		case <-w.Quit:
			return
		case <-ticker.C:
			job, err := w.Queue.Dequeue()
			if err != nil {
				logger.GlobalBuffer.Error("Worker error dequeuing: %v", err)
				continue
			}
			if job == nil {
				continue
			}

			w.Process(job)
		}
	}
}

func (w *Worker) Process(job *Job) {
	logger.GlobalBuffer.Info("🚀 Processing job %d (Brand: %s, Type: %s)", job.ID, job.BrandID, job.Type)

	agentInstance, err := w.AgentFactory(job.BrandID)
	if err != nil {
		logger.GlobalBuffer.Error("Worker failed to create agent for brand %s: %v", job.BrandID, err)
		w.Queue.Fail(job.ID, err.Error(), false)
		return
	}

	var runErr error
	switch job.Type {
	case JobTypeRun:
		runErr = agentInstance.Run()
	case JobTypeSync:
		runErr = agentInstance.SyncAnalytics()
	case JobTypePlan:
		runErr = agentInstance.PlanBatch(5) // Default to 5 posts for now
	case JobTypePublish:
		// Payload contains the ScheduledPostID
		posts, err := agentInstance.Store.GetScheduledPosts(job.BrandID)
		if err != nil {
			runErr = err
		} else {
			var target *models.ScheduledPost
			for _, p := range posts {
				if p.ID == job.Payload {
					target = &p
					break
				}
			}
			if target != nil {
				runErr = agentInstance.PublishScheduledPost(*target)
			} else {
				runErr = fmt.Errorf("scheduled post %s not found", job.Payload)
			}
		}
	default:
		runErr = fmt.Errorf("unknown job type: %s", job.Type)
	}

	if runErr != nil {
		logger.GlobalBuffer.Error("Job %d failed: %v", job.ID, runErr)
		// Retry if it's the first few failures
		shouldRetry := job.Retries < 3
		w.Queue.Fail(job.ID, runErr.Error(), shouldRetry)
	} else {
		logger.GlobalBuffer.Info("✅ Job %d completed successfully!", job.ID)
		w.Queue.Ack(job.ID)
	}
}

// DefaultAgentFactory helper to create the factory.
func DefaultAgentFactory(store memory.Store, search tools.SearchTool, llm tools.LLMTool, social tools.SocialClient, embedding tools.EmbeddingTool, analytics tools.AnalyticsFetcher, dataDir string) AgentFactory {
	return func(brandID string) (*agent.Agent, error) {
		brand, _, err := store.GetBrand(brandID)
		if err != nil {
			return nil, err
		}

		vectorStore := memory.NewLocalVectorStore(filepath.Join(dataDir, brandID, "vectors.json"))
		return agent.NewAgent(brand, search, llm, social, store, vectorStore, embedding, analytics), nil
	}
}
