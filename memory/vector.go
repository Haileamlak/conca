package memory

import (
	"context"
	"math"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// VectorRecord represents a stored embedding with metadata.
type VectorRecord struct {
	ID       string                 `json:"id" bson:"_id"`
	Vector   []float32              `json:"vector" bson:"vector"`
	Metadata map[string]interface{} `json:"metadata" bson:"metadata"`
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
	UpdateMetadata(id string, metadata map[string]interface{}) error
}

// MongoVectorStore implements VectorStore using MongoDB.
type MongoVectorStore struct {
	collection *mongo.Collection
}

func NewMongoVectorStore(db *mongo.Database, collectionName string) *MongoVectorStore {
	return &MongoVectorStore{collection: db.Collection(collectionName)}
}

func (m *MongoVectorStore) Add(record VectorRecord) error {
	_, err := m.collection.InsertOne(context.Background(), record)
	return err
}

func (m *MongoVectorStore) Query(queryVector []float32, topK int) ([]SearchResult, error) {
	cursor, err := m.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	var results []SearchResult
	var records []VectorRecord
	if err := cursor.All(context.Background(), &records); err != nil {
		return nil, err
	}

	for _, rec := range records {
		score := m.cosineSimilarity(queryVector, rec.Vector)
		results = append(results, SearchResult{
			ID:       rec.ID,
			Score:    score,
			Metadata: rec.Metadata,
		})
	}

	// Sort by score descending
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

func (m *MongoVectorStore) cosineSimilarity(v1, v2 []float32) float32 {
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

func (m *MongoVectorStore) UpdateMetadata(id string, metadata map[string]interface{}) error {
	update := bson.M{"$set": bson.M{}}
	for k, v := range metadata {
		update["$set"].(bson.M)["metadata."+k] = v
	}

	_, err := m.collection.UpdateOne(context.Background(), bson.M{"_id": id}, update)
	return err
}
