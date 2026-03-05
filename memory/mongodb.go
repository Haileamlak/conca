package memory

import (
	"content-creator-agent/models"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoStore struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewMongoStore(uri, dbName string) (*MongoStore, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping mongodb: %w", err)
	}

	return &MongoStore{
		client: client,
		db:     client.Database(dbName),
	}, nil
}

func (m *MongoStore) Close() error {
	return m.client.Disconnect(context.Background())
}

func (m *MongoStore) Database() *mongo.Database {
	return m.db
}

// Collections
func (m *MongoStore) posts() *mongo.Collection     { return m.db.Collection("posts") }
func (m *MongoStore) brands() *mongo.Collection    { return m.db.Collection("brands") }
func (m *MongoStore) scheduled() *mongo.Collection { return m.db.Collection("scheduled_posts") }
func (m *MongoStore) users() *mongo.Collection     { return m.db.Collection("users") }

// --- Posts ---

func (m *MongoStore) SavePost(post models.Post) error {
	_, err := m.posts().InsertOne(context.Background(), post)
	return err
}

func (m *MongoStore) GetHistory(brandID string) ([]models.Post, error) {
	cursor, err := m.posts().Find(context.Background(), bson.M{"brand_id": brandID}, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	var posts []models.Post
	err = cursor.All(context.Background(), &posts)
	return posts, err
}

func (m *MongoStore) GetGlobalHistory(userID string, limit int) ([]models.Post, error) {
	// First get brand IDs for this user
	brands, err := m.ListBrands(userID)
	if err != nil {
		return nil, err
	}
	var brandIDs []string
	for _, b := range brands {
		brandIDs = append(brandIDs, b.ID)
	}

	findOptions := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	if limit > 0 {
		findOptions.SetLimit(int64(limit))
	}

	cursor, err := m.posts().Find(context.Background(), bson.M{"brand_id": bson.M{"$in": brandIDs}}, findOptions)
	if err != nil {
		return nil, err
	}
	var posts []models.Post
	err = cursor.All(context.Background(), &posts)
	return posts, err
}

func (m *MongoStore) GetAnalytics(brandID string) ([]models.Analytics, error) {
	posts, err := m.GetHistory(brandID)
	if err != nil {
		return nil, err
	}
	var analytics []models.Analytics
	for _, p := range posts {
		analytics = append(analytics, p.Analytics)
	}
	return analytics, nil
}

func (m *MongoStore) UpdateAnalytics(brandID string, postID string, analytics models.Analytics) error {
	_, err := m.posts().UpdateOne(
		context.Background(),
		bson.M{"_id": postID, "brand_id": brandID},
		bson.M{"$set": bson.M{"analytics": analytics, "updated_at": time.Now()}},
	)
	return err
}

func (m *MongoStore) GetGlobalAnalytics(userID string) (models.GlobalAnalytics, error) {
	brands, err := m.ListBrands(userID)
	if err != nil {
		return models.GlobalAnalytics{}, err
	}
	var brandIDs []string
	for _, b := range brands {
		brandIDs = append(brandIDs, b.ID)
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"brand_id": bson.M{"$in": brandIDs}}}},
		{{Key: "$group", Value: bson.M{
			"_id":               nil,
			"total_impressions": bson.M{"$sum": "$analytics.views"},
			"total_likes":       bson.M{"$sum": "$analytics.likes"},
			"total_shares":      bson.M{"$sum": "$analytics.shares"},
			"total_comments":    bson.M{"$sum": "$analytics.comments"},
		}}},
	}

	cursor, err := m.posts().Aggregate(context.Background(), pipeline)
	if err != nil {
		return models.GlobalAnalytics{}, err
	}

	var results []models.GlobalAnalytics
	if err = cursor.All(context.Background(), &results); err != nil {
		return models.GlobalAnalytics{}, err
	}

	if len(results) == 0 {
		return models.GlobalAnalytics{}, nil
	}
	return results[0], nil
}

func (m *MongoStore) GetBrandPerformance(userID string) ([]models.BrandPerformance, error) {
	brands, err := m.ListBrands(userID)
	if err != nil {
		return nil, err
	}

	var results []models.BrandPerformance
	for _, b := range brands {
		posts, _ := m.GetHistory(b.ID)
		if len(posts) == 0 {
			continue
		}

		perf := models.BrandPerformance{
			BrandID:   b.ID,
			BrandName: b.Name,
			PostCount: len(posts),
		}

		totalLikes := 0
		totalShares := 0
		for _, p := range posts {
			totalLikes += p.Analytics.Likes
			totalShares += p.Analytics.Shares
		}

		perf.AvgLikes = float64(totalLikes) / float64(len(posts))
		perf.AvgShares = float64(totalShares) / float64(len(posts))

		score := int(perf.AvgLikes*2 + perf.AvgShares*5)
		if score > 100 {
			score = 100
		}
		perf.Score = score

		results = append(results, perf)
	}
	return results, nil
}

// --- Brands ---

func (m *MongoStore) SaveBrand(brand models.BrandProfile, userID string) error {
	brand.UserID = userID
	opts := options.UpdateOne().SetUpsert(true)
	_, err := m.brands().UpdateOne(
		context.Background(),
		bson.M{"_id": brand.ID},
		bson.M{"$set": brand},
		opts,
	)
	return err
}

func (m *MongoStore) GetBrand(id string) (models.BrandProfile, string, error) {
	var brand models.BrandProfile
	err := m.brands().FindOne(context.Background(), bson.M{"_id": id}).Decode(&brand)
	if err != nil {
		return models.BrandProfile{}, "", err
	}
	return brand, brand.UserID, nil
}

func (m *MongoStore) ListBrands(userID string) ([]models.BrandProfile, error) {
	cursor, err := m.brands().Find(context.Background(), bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	var brands []models.BrandProfile
	err = cursor.All(context.Background(), &brands)
	return brands, err
}

func (m *MongoStore) ListAllBrands() ([]models.BrandProfile, error) {
	cursor, err := m.brands().Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	var brands []models.BrandProfile
	err = cursor.All(context.Background(), &brands)
	return brands, err
}

func (m *MongoStore) DeleteBrand(id string) error {
	_, err := m.brands().DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return err
	}
	// Also delete posts and scheduled posts
	m.posts().DeleteMany(context.Background(), bson.M{"brand_id": id})
	m.scheduled().DeleteMany(context.Background(), bson.M{"brand_id": id})
	return nil
}

// --- Calendar ---

func (m *MongoStore) SaveScheduledPost(post models.ScheduledPost) error {
	_, err := m.scheduled().InsertOne(context.Background(), post)
	return err
}

func (m *MongoStore) GetScheduledPosts(brandID string) ([]models.ScheduledPost, error) {
	cursor, err := m.scheduled().Find(context.Background(), bson.M{"brand_id": brandID}, options.Find().SetSort(bson.D{{Key: "scheduled_at", Value: 1}}))
	if err != nil {
		return nil, err
	}
	var posts []models.ScheduledPost
	err = cursor.All(context.Background(), &posts)
	return posts, err
}

func (m *MongoStore) UpdateScheduledPostStatus(postID string, status models.PostStatus) error {
	_, err := m.scheduled().UpdateOne(
		context.Background(),
		bson.M{"_id": postID},
		bson.M{"$set": bson.M{"status": status, "updated_at": time.Now()}},
	)
	return err
}

func (m *MongoStore) UpdateScheduledPost(postID string, topic, content string) error {
	_, err := m.scheduled().UpdateOne(
		context.Background(),
		bson.M{"_id": postID},
		bson.M{"$set": bson.M{"topic": topic, "content": content, "updated_at": time.Now()}},
	)
	return err
}

func (m *MongoStore) GetPendingScheduledPosts() ([]models.ScheduledPost, error) {
	cursor, err := m.scheduled().Find(context.Background(), bson.M{
		"status":       models.StatusScheduled,
		"scheduled_at": bson.M{"$lt": time.Now()},
	})
	if err != nil {
		return nil, err
	}
	var posts []models.ScheduledPost
	err = cursor.All(context.Background(), &posts)
	return posts, err
}

// --- Users ---

func (m *MongoStore) CreateUser(email, passwordHash string) (string, error) {
	count, err := m.users().CountDocuments(context.Background(), bson.M{"email": email})
	if err != nil {
		return "", err
	}
	if count > 0 {
		return "", fmt.Errorf("user already exists")
	}

	userID := fmt.Sprintf("u-%d", time.Now().Unix())
	user := models.User{
		ID:           userID,
		Email:        email,
		PasswordHash: passwordHash,
	}

	_, err = m.users().InsertOne(context.Background(), user)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func (m *MongoStore) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := m.users().FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *MongoStore) GetUserByID(id string) (*models.User, error) {
	var user models.User
	err := m.users().FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
