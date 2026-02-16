package memory

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"sync"
)

// VectorRecord represents a stored embedding with metadata.
type VectorRecord struct {
	ID       string                 `json:"id"`
	Vector   []float32              `json:"vector"`
	Metadata map[string]interface{} `json:"metadata"`
}

// SearchResult represents a single match from the vector store.
type SearchResult struct {
	ID       string                 `json:"id"`
	Score    float32                `json:"score"`
	Metadata map[string]interface{} `json:"metadata"`
}

// VectorStore defines the interface for semantic storage and retrieval.
type VectorStore interface {
	Add(record VectorRecord) error
	Query(queryVector []float32, topK int) ([]SearchResult, error)
}

// LocalVectorStore implements VectorStore using a local JSON file.
type LocalVectorStore struct {
	FilePath string
	records  []VectorRecord
	mu       sync.RWMutex
}

func NewLocalVectorStore(filePath string) *LocalVectorStore {
	store := &LocalVectorStore{FilePath: filePath}
	store.load()
	return store
}

func (l *LocalVectorStore) load() {
	data, err := os.ReadFile(l.FilePath)
	if err == nil {
		json.Unmarshal(data, &l.records)
	}
}

func (l *LocalVectorStore) save() error {
	dir := filepath.Dir(l.FilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(l.records, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(l.FilePath, data, 0644)
}

func (l *LocalVectorStore) Add(record VectorRecord) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.records = append(l.records, record)
	return l.save()
}

func (l *LocalVectorStore) Query(queryVector []float32, topK int) ([]SearchResult, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var results []SearchResult
	for _, rec := range l.records {
		score := l.cosineSimilarity(queryVector, rec.Vector)
		results = append(results, SearchResult{
			ID:       rec.ID,
			Score:    score,
			Metadata: rec.Metadata,
		})
	}

	// Sort by score descending (simple bubble-sort-like or use sort.Slice for topK)
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[i].Score < results[j].Score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	if len(results) > topK {
		results = results[:topK]
	}

	return results, nil
}

func (l *LocalVectorStore) cosineSimilarity(v1, v2 []float32) float32 {
	if len(v1) != len(v2) {
		return 0
	}
	var dotProduct, normV1, normV2 float64
	for i := range v1 {
		dotProduct += float64(v1[i]) * float64(v2[i])
		normV1 += float64(v1[i]) * float64(v1[i])
		normV2 += float64(v2[i]) * float64(v2[i])
	}
	if normV1 == 0 || normV2 == 0 {
		return 0
	}
	return float32(dotProduct / (math.Sqrt(normV1) * math.Sqrt(normV2)))
}
