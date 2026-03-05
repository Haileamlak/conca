package scheduler

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type JobType string

const (
	JobTypeRun     JobType = "run"
	JobTypeSync    JobType = "sync"
	JobTypePlan    JobType = "plan"    // Batch planning
	JobTypePublish JobType = "publish" // Specific post publication
)

type JobStatus string

const (
	StatusPending   JobStatus = "pending"
	StatusRunning   JobStatus = "running"
	StatusCompleted JobStatus = "completed"
	StatusFailed    JobStatus = "failed"
)

type Job struct {
	ID        int64     `bson:"_id"`
	BrandID   string    `bson:"brand_id"`
	Type      JobType   `bson:"type"`
	Status    JobStatus `bson:"status"`
	Retries   int       `bson:"retries"`
	NextRunAt time.Time `bson:"next_run_at"`
	Payload   string    `bson:"payload"`
	Error     string    `bson:"error"`
}

// Queue defines the job management interface
type Queue interface {
	Enqueue(brandID string, jobType JobType, delay time.Duration, payload string) error
	Dequeue() (*Job, error)
	Ack(jobID int64) error
	Fail(jobID int64, errMsg string, retry bool) error
	HasPendingJob(brandID string) (bool, error)
}

type MongoQueue struct {
	collection *mongo.Collection
}

func NewMongoQueue(db *mongo.Database) *MongoQueue {
	return &MongoQueue{collection: db.Collection("jobs")}
}

func (q *MongoQueue) Enqueue(brandID string, jobType JobType, delay time.Duration, payload string) error {
	nextRun := time.Now().Add(delay)
	job := bson.M{
		"_id":         time.Now().UnixNano(), // Simple unique ID
		"brand_id":    brandID,
		"type":        string(jobType),
		"status":      string(StatusPending),
		"retries":     0,
		"next_run_at": nextRun,
		"payload":     payload,
		"error":       "",
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
	}
	_, err := q.collection.InsertOne(context.Background(), job)
	return err
}

func (q *MongoQueue) Dequeue() (*Job, error) {
	filter := bson.M{
		"status":      string(StatusPending),
		"next_run_at": bson.M{"$lte": time.Now()},
	}
	update := bson.M{
		"$set": bson.M{
			"status":     string(StatusRunning),
			"updated_at": time.Now(),
		},
	}
	opts := options.FindOneAndUpdate().SetSort(bson.D{{Key: "next_run_at", Value: 1}}).SetReturnDocument(options.After)

	var job Job
	err := q.collection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&job)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (q *MongoQueue) Ack(jobID int64) error {
	_, err := q.collection.DeleteOne(context.Background(), bson.M{"_id": jobID})
	return err
}

func (q *MongoQueue) Fail(jobID int64, errMsg string, retry bool) error {
	if retry {
		delay := 5 * time.Minute
		update := bson.M{
			"$set": bson.M{
				"status":      string(StatusPending),
				"next_run_at": time.Now().Add(delay),
				"error":       errMsg,
				"updated_at":  time.Now(),
			},
			"$inc": bson.M{"retries": 1},
		}
		_, err := q.collection.UpdateOne(context.Background(), bson.M{"_id": jobID}, update)
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"status":     string(StatusFailed),
			"error":      errMsg,
			"updated_at": time.Now(),
		},
	}
	_, err := q.collection.UpdateOne(context.Background(), bson.M{"_id": jobID}, update)
	return err
}

func (q *MongoQueue) HasPendingJob(brandID string) (bool, error) {
	filter := bson.M{
		"brand_id": brandID,
		"status":   bson.M{"$in": []string{string(StatusPending), string(StatusRunning)}},
	}
	count, err := q.collection.CountDocuments(context.Background(), filter)
	return count > 0, err
}

func (q *MongoQueue) Close() error {
	return nil
}
