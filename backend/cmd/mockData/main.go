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
	minDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	maxDate := time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC)
	duration := maxDate.Sub(minDate)
	for i := 0; i < 100000; i++ {
		startHour := rand.Intn(24)
		startMinute := rand.Intn(60)

		endHour := startHour + rand.Intn(24-startHour)
		endMinute := rand.Intn(60)

		startTime := fmt.Sprintf("%02d:%02d", startHour, startMinute)
		endTime := fmt.Sprintf("%02d:%02d", endHour, endMinute)

		events[i] = models.Event{
			Name:      "Event " + strconv.Itoa(i),
			Date:      minDate.Add(time.Duration(rand.Int63n(int64(duration)))),
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
