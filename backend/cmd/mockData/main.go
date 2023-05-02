package main

import (
	"aimet-test/configs"
	"aimet-test/internal/databases"
	"aimet-test/internal/models"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

func main() {
	configs.InitConfig()

	databases.InitMongoClient()
	defer databases.CloseMongoClient()

	mongoClient := databases.GetMongoClient()
	eventCol := mongoClient.Database(viper.GetString("mongo.dbname")).Collection("events")

	rand.Seed(time.Now().UnixNano())

	events := make([]interface{}, 100000)
	for i := 0; i < 100000; i++ {
		// generate random start time
		startHour := rand.Intn(24)
		startMinute := rand.Intn(60)

		// generate random end time
		endHour := startHour + rand.Intn(24-startHour)
		endMinute := rand.Intn(60)

		// format time
		startTime := fmt.Sprintf("%02d:%02d", startHour, startMinute)
		endTime := fmt.Sprintf("%02d:%02d", endHour, endMinute)

		events[i] = models.Event{
			Name:      "Event " + strconv.Itoa(i),
			Date:      randomDate(),
			StartTime: startTime,
			EndTime:   endTime,
		}
		if i%1000 == 0 {
			fmt.Println(i)
		}
	}

	_, err := eventCol.InsertMany(nil, events)
	if err != nil {
		panic(err)
	}
	fmt.Println("Complete insert event data (~100,000 records)")
}

func randomDate() time.Time {
	// Set the range of years.
	minYear := 2000
	maxYear := 2999

	// Generate a random year.
	year := rand.Intn(maxYear-minYear+1) + minYear

	// Generate a random day of the year.
	daysInYear := time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC).YearDay()
	dayOfYear := rand.Intn(daysInYear) + 1

	// Create a time.Time object for the random date.
	date := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, dayOfYear-1)

	return date
}
