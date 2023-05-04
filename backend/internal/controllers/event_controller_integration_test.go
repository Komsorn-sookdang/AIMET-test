//go:build integration

package controllers_test

import (
	"aimet-test/internal/controllers"
	"aimet-test/internal/models"
	"aimet-test/internal/repositories"
	"aimet-test/internal/usecases"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetFilteredEventsIntegrationUsecase(t *testing.T) {
	t.Run("should return 200 OK", func(t *testing.T) {

		day := 10
		keyword := "very%20long%20name"
		expected := models.GetEventResponse{
			Count: 1,
			Events: []*models.Event{
				{
					Name: "Test Event 3 very long name",
					Date: time.Date(2021, 1, 10, 0, 0, 0, 0, time.UTC),
					StartTime: "9 AM",
					EndTime: "11 AM",
				},
			},
		}

		eventRepo := repositories.NewEventRepositoryMock()
		eventRepo.On("FindByMonth", 1, 2021).Return([]*models.Event{
			{
				Name: "Test Event 1",
				Date: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				StartTime: "15:00",
				EndTime: "20:00",
			},
			{
				Name: "Test Event 2",
				Date: time.Date(2021, 1, 10, 0, 0, 0, 0, time.UTC),
				StartTime: "01:00",
				EndTime: "11:00",
			},
			{
				Name: "Test Event 3 very long name",
				Date: time.Date(2021, 1, 10, 0, 0, 0, 0, time.UTC),
				StartTime: "09:00",
				EndTime: "11:00",
			},
			{
				Name: "Test Event 4 very long name",
				Date: time.Date(2021, 1, 20, 0, 0, 0, 0, time.UTC),
				StartTime: "21:00",
				EndTime: "23:00",
			},
		}, nil)
		
		eventUsc := usecases.NewEventUsecase(eventRepo)
		eventCtl := controllers.NewEventController(eventUsc)

		gin.SetMode(gin.TestMode)
		app := gin.New()
		app.GET("/events", eventCtl.GetFilteredEvents)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/events?month=1&year=2021&day=%v&keyword=%v", day, keyword), nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if assert.Equal(t, http.StatusOK, w.Code) {
			expectedJSON, _ := json.Marshal(expected)
			assert.Equal(t, expectedJSON, w.Body.Bytes())
		}
	})

	t.Run("should return 400 Bad Request when month is not provided", func(t *testing.T) {

		eventRepo := repositories.NewEventRepositoryMock()
		eventUsc := usecases.NewEventUsecase(eventRepo)
		eventCtl := controllers.NewEventController(eventUsc)

		gin.SetMode(gin.TestMode)
		app := gin.New()
		app.GET("/events", eventCtl.GetFilteredEvents)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/events?year=2021&day=10"), nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return 400 Bad Request when year is invalid", func(t *testing.T) {

		year := -1

		eventRepo := repositories.NewEventRepositoryMock()
		eventUsc := usecases.NewEventUsecase(eventRepo)
		eventCtl := controllers.NewEventController(eventUsc)

		gin.SetMode(gin.TestMode)
		app := gin.New()
		app.GET("/events", eventCtl.GetFilteredEvents)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/events?month=1&year=%v&day=10", year), nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return 400 Bad Request when date is invalid", func(t *testing.T) {

		day := 29
		year := 2023
		month := 2

		eventRepo := repositories.NewEventRepositoryMock()
		eventUsc := usecases.NewEventUsecase(eventRepo)
		eventCtl := controllers.NewEventController(eventUsc)

		gin.SetMode(gin.TestMode)
		app := gin.New()
		app.GET("/events", eventCtl.GetFilteredEvents)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/events?month=%v&year=%v&day=%v", month, year, day), nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return 400 Bad Request when sort_by is invalid", func(t *testing.T) {

		sort_by := "invalid"

		eventRepo := repositories.NewEventRepositoryMock()
		eventUsc := usecases.NewEventUsecase(eventRepo)
		eventCtl := controllers.NewEventController(eventUsc)

		gin.SetMode(gin.TestMode)
		app := gin.New()
		app.GET("/events", eventCtl.GetFilteredEvents)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/events?month=1&year=2021&day=10&sort_by=%v", sort_by), nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return 500 Internal Server Error when repository returns error", func(t *testing.T) {
		
		eventRepo := repositories.NewEventRepositoryMock()

		eventRepo.On("FindByMonth", 1, 2021).Return([]*models.Event{}, errors.New("error"))
		eventUsc := usecases.NewEventUsecase(eventRepo)
		eventCtl := controllers.NewEventController(eventUsc)

		gin.SetMode(gin.TestMode)
		app := gin.New()
		app.GET("/events", eventCtl.GetFilteredEvents)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/events?month=1&year=2021&day=10"), nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
