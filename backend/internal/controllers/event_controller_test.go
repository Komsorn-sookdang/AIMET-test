package controllers_test

import (
	"aimet-test/internal/controllers"
	"aimet-test/internal/models"
	"aimet-test/internal/usecases"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetFilteredEvents(t *testing.T) {
	t.Run("should return 200 OK", func(t *testing.T) {

		month := 1
		expected := models.GetEventResponse{
			Count: 2,
			Events: []*models.Event{
				{
					Name: "Event 1",
					Date: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					StartTime: "10 AM",
					EndTime: "11 AM",
				},
				{
					Name: "Event 1",
					Date: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
					StartTime: "10 AM",
					EndTime: "11 AM",
				},
			},
		}
		
		eventUsc := usecases.NewEventUsecaseMock()
		eventUsc.On("GetFilteredEvents", &models.GetEventQuery{
			Month: month,
		}).Return(expected.Events, nil)

		eventCtl := controllers.NewEventController(eventUsc)
		
		gin.SetMode(gin.TestMode)
		app := gin.New()
		app.GET("/events", eventCtl.GetFilteredEvents)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/events?month=%v", month), nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if assert.Equal(t, http.StatusOK, w.Code) {
			expectedJSON, _ := json.Marshal(expected)
			assert.Equal(t, expectedJSON, w.Body.Bytes())
		}
	})

	t.Run("should return 400 Bad Request when month is not provided", func(t *testing.T) {
		eventUsc := usecases.NewEventUsecaseMock()
		mockCtl := controllers.NewEventController(eventUsc)
		
		gin.SetMode(gin.TestMode)
		app := gin.New()
		app.GET("/events", mockCtl.GetFilteredEvents)

		req := httptest.NewRequest(http.MethodGet, "/events", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if assert.Equal(t, http.StatusBadRequest, w.Code) {
			expectedJSON, _ := json.Marshal(gin.H{"error": "Key: 'GetEventQuery.Month' Error:Field validation for 'Month' failed on the 'required' tag"})
			assert.Equal(t, expectedJSON, w.Body.Bytes())
		}
	})
}
