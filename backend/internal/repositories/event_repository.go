package repositories

import (
	"aimet-test/internal/domains"
	"aimet-test/internal/models"
	"context"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type eventRepository struct {
	collection *mongo.Collection
}

func NewEventRepository(mongoClient *mongo.Client) domains.EventRepository {
	collection := mongoClient.Database(viper.GetString("mongo.dbname")).Collection("events")
	return &eventRepository{collection: collection}
}

func (r *eventRepository) FindByMonth(month int, year int) ([]*models.Event, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

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
	cursor, err := r.collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// Iterate through the cursor and decode each event
	var events []*models.Event
	for cursor.Next(context.TODO()) {
		var event models.Event
		if err := cursor.Decode(&event); err != nil {
			continue
		}
		events = append(events, &event)
	}

	return events, nil
}
