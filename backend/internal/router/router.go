package router

import (
	"aimet-test/internal/controllers"
	"aimet-test/internal/databases"
	"aimet-test/internal/repositories"
	"aimet-test/internal/usecases"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello World"})
	})

	mongoClient := databases.GetMongoClient()
	eventRepo := repositories.NewEventRepository(mongoClient)
	eventUsc := usecases.NewEventUsecase(eventRepo)
	eventCtl := controllers.NewEventController(eventUsc)

	r.GET("/events", eventCtl.GetFilteredEvents)

	return r
}
