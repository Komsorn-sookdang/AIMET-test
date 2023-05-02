package repositories

import (
	"aimet-test/internal/domains"
	"aimet-test/internal/models"
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type eventRepository struct {
	eventCollection *mongo.Collection
	redisClient     *redis.Client
}

func NewEventRepository(mongoClient *mongo.Client, redisClient *redis.Client) domains.EventRepository {
	collection := mongoClient.Database(viper.GetString("mongo.dbname")).Collection("events")
	return &eventRepository{
		eventCollection: collection,
		redisClient:     redisClient,
	}
}

func (r *eventRepository) FindByMonth(month int, year int) ([]*models.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	key := strconv.Itoa(year) + "-" + strconv.Itoa(month)
	val, err := r.redisClient.Get(ctx, key).Bytes()
	if err == nil {
		// If the key exists in Redis, return the value
		var events []*models.Event
		if err := json.Unmarshal(val, &events); err != nil {
			return nil, err
		}
		return events, nil
	}

	if err != redis.Nil {
		return nil, err
	}

	// Set up a filter to retrieve events in the specified month and year
	startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)
	filter := bson.M{
		"date": bson.M{
			"$gte": startOfMonth,
			"$lte": endOfMonth,
		},
	}
	opts := options.Find().SetSort(bson.M{"date": 1})

	// Retrieve the events that match the filter
	cursor, err := r.eventCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Iterate through the cursor and decode each event
	var events []*models.Event
	for cursor.Next(context.TODO()) {
		var event models.Event
		if err := cursor.Decode(&event); err != nil {
			continue
		}
		events = append(events, &event)
	}

	// Store the result in Redis
	data, err := json.Marshal(events)
	if err != nil {
		return nil, err
	}
	if err := r.redisClient.Set(ctx, key, data, 5*time.Minute).Err(); err != nil {
		return nil, err
	}

	return events, nil
}
