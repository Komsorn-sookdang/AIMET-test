package usecases_test

import (
	"aimet-test/internal/models"
	"aimet-test/internal/repositories"
	"aimet-test/internal/usecases"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetFilteredEvents(t *testing.T) {

	type testCase struct {
		name 				string
		filter 			*models.GetEventQuery
		expected 		[]*models.Event
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "Test Get Events Default",
			filter: &models.GetEventQuery{
				Month: 1,
				Year: 2021,
				Day: 0,
				SortBy: "",
				Keyword: "",
			},
			expected: []*models.Event{
				{
					Name: "Test Event 1",
					Date: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					StartTime: "3 PM",
					EndTime: "8 PM",
				},
				{
					Name: "Test Event 2",
					Date: time.Date(2021, 1, 10, 0, 0, 0, 0, time.UTC),
					StartTime: "1 AM",
					EndTime: "11 AM",
				},
				{
					Name: "Test Event 3 very long name",
					Date: time.Date(2021, 1, 10, 0, 0, 0, 0, time.UTC),
					StartTime: "9 AM",
					EndTime: "11 AM",
				},
				{
					Name: "Test Event 4 very long name",
					Date: time.Date(2021, 1, 20, 0, 0, 0, 0, time.UTC),
					StartTime: "9 PM",
					EndTime: "11 PM",
				},
			},
			expectedErr: nil,
		},
		{
			name: "Test Filter By Day",
			filter: &models.GetEventQuery{
				Month: 1,
				Year: 2021,
				Day: 10,
				SortBy: "",
				Keyword: "",
			},
			expected: []*models.Event{
				{
					Name: "Test Event 2",
					Date: time.Date(2021, 1, 10, 0, 0, 0, 0, time.UTC),
					StartTime: "1 AM",
					EndTime: "11 AM",
				},
				{
					Name: "Test Event 3 very long name",
					Date: time.Date(2021, 1, 10, 0, 0, 0, 0, time.UTC),
					StartTime: "9 AM",
					EndTime: "11 AM",
				},
			},
			expectedErr: nil,
		},
		{
			name: "Test Filter By Keyword",
			filter: &models.GetEventQuery{
				Month: 1,
				Year: 2021,
				Day: 0,
				SortBy: "",
				Keyword: "very long name",
			},
			expected: []*models.Event{
				{
					Name: "Test Event 3 very long name",
					Date: time.Date(2021, 1, 10, 0, 0, 0, 0, time.UTC),
					StartTime: "9 AM",
					EndTime: "11 AM",
				},
				{
					Name: "Test Event 4 very long name",
					Date: time.Date(2021, 1, 20, 0, 0, 0, 0, time.UTC),
					StartTime: "9 PM",
					EndTime: "11 PM",
				},
			},
			expectedErr: nil,
		},
		{
			name: "Test Filter By Keyword and Day",
			filter: &models.GetEventQuery{
				Month: 1,
				Year: 2021,
				Day: 10,
				SortBy: "",
				Keyword: "very long name",
			},
			expected: []*models.Event{
				{
					Name: "Test Event 3 very long name",
					Date: time.Date(2021, 1, 10, 0, 0, 0, 0, time.UTC),
					StartTime: "9 AM",
					EndTime: "11 AM",
				},
			},
			expectedErr: nil,
		},
		{
			name: "Test Sort By Time",
			filter: &models.GetEventQuery{
				Month: 1,
				Year: 2021,
				Day: 0,
				SortBy: "time",
				Keyword: "",
			},
			expected: []*models.Event{
				{
					Name: "Test Event 2",
					Date: time.Date(2021, 1, 10, 0, 0, 0, 0, time.UTC),
					StartTime: "1 AM",
					EndTime: "11 AM",
				},
				{
					Name: "Test Event 3 very long name",
					Date: time.Date(2021, 1, 10, 0, 0, 0, 0, time.UTC),
					StartTime: "9 AM",
					EndTime: "11 AM",
				},
				{
					Name: "Test Event 1",
					Date: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					StartTime: "3 PM",
					EndTime: "8 PM",
				},
				{
					Name: "Test Event 4 very long name",
					Date: time.Date(2021, 1, 20, 0, 0, 0, 0, time.UTC),
					StartTime: "9 PM",
					EndTime: "11 PM",
				},
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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

			events, _ := eventUsc.GetFilteredEvents(tc.filter)
			assert.Equal(t, tc.expected, events)
		})
	}

	// Case Repository Error
	t.Run("Test Repository Error", func(t *testing.T) {
		eventRepo := repositories.NewEventRepositoryMock()
		eventRepo.On("FindByMonth", 0, time.Now().Year()).Return([]*models.Event{}, errors.New("error"))

		eventUsc := usecases.NewEventUsecase(eventRepo)

		_, err := eventUsc.GetFilteredEvents(&models.GetEventQuery{})
		assert.ErrorIs(t, err, usecases.ErrRepository)
	})
}

func TestValidateGetEventQuery(t *testing.T) {
	type testCase struct {
		name string
		query *models.GetEventQuery
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "Test Sort By Date",
			query: &models.GetEventQuery{
				Month: 1,
				Year: 2021,
				Day: 1,
				SortBy: "date",
				Keyword: "test",
			},
			expectedErr: nil,
		},
		{
			name: "Test Sort By Time",
			query: &models.GetEventQuery{
				Month: 1,
				Year: 2021,
				Day: 1,
				SortBy: "time",
				Keyword: "test",
			},
			expectedErr: nil,
		},
		{
			name: "Test Sort By Invalid",
			query: &models.GetEventQuery{
				Month: 1,
				Year: 2021,
				Day: 1,
				SortBy: "invalid",
				Keyword: "test",
			},
			expectedErr: usecases.ErrInvalidSortBy,
		},
		{
			name: "Test Sort By Default",
			query: &models.GetEventQuery{
				Month: 1,
				Year: 2021,
				Day: 1,
				SortBy: "",
				Keyword: "test",
			},
			expectedErr: nil,
		},
		{
			name: "Test Invalid Month 1",
			query: &models.GetEventQuery{
				Month: 0,
				Year: 2021,
				Day: 1,
				SortBy: "date",
				Keyword: "test",
			},
			expectedErr: usecases.ErrInvalidMonth,
		},
		{
			name: "Test Invalid Month 2",
			query: &models.GetEventQuery{
				Month: 13,
				Year: 2021,
				Day: 1,
				SortBy: "date",
				Keyword: "test",
			},
			expectedErr: usecases.ErrInvalidMonth,
		},
		{
			name: "Test Invalid Year 1",
			query: &models.GetEventQuery{
				Month: 1,
				Year: -1,
				Day: 1,
				SortBy: "date",
				Keyword: "test",
			},
			expectedErr: usecases.ErrInvalidYear,
		},
		{
			name: "Test Invalid Year 2",
			query: &models.GetEventQuery{
				Month: 1,
				Year: 10000,
				Day: 1,
				SortBy: "date",
				Keyword: "test",
			},
			expectedErr: usecases.ErrInvalidYear,
		},
		{
			name: "Test Default Date",
			query: &models.GetEventQuery{
				Month: 1,
				Year: 2021,
				Day: 0,
				SortBy: "date",
				Keyword: "test",
			},
			expectedErr: nil,
		},
		{
			name: "Test Invalid Date 1",
			query: &models.GetEventQuery{
				Month: 1,
				Year: 2021,
				Day: -1,
				SortBy: "date",
				Keyword: "test",
			},
			expectedErr: usecases.ErrInvalidDate,
		},
		{
			name: "Test Invalid Date 2",
			query: &models.GetEventQuery{
				Month: 1,
				Year: 2021,
				Day: 32,
				SortBy: "date",
				Keyword: "test",
			},
			expectedErr: usecases.ErrInvalidDate,
		},
		{
			name: "Test Invalid Date 3",
			query: &models.GetEventQuery{
				Month: 2,
				Year: 2021,
				Day: 30,
				SortBy: "date",
				Keyword: "test",
			},
			expectedErr: usecases.ErrInvalidDate,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := usecases.ValidateGetEventQuery(tc.query)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestSortEventsByTime(t *testing.T) {
	type testCase struct {
		name string
		events []*models.Event
		expected []*models.Event
	}

	testCases := []testCase{
		{
			name: "Test 1",
			events: []*models.Event{
				{
					Name: "Event 1",
					Date: time.Date(2021, 1, 1, 9, 0, 0, 0, time.UTC),
					StartTime: "09:00",
					EndTime: "10:00",
				},
				{
					Name: "Event 2",
					Date: time.Date(2021, 1, 1, 9, 0, 0, 0, time.UTC),
					StartTime: "10:00",
					EndTime: "11:00",
				},
				{
					Name: "Event 3",
					Date: time.Date(2021, 1, 1, 9, 0, 0, 0, time.UTC),
					StartTime: "11:00",
					EndTime: "12:00",
				},
			},
			expected: []*models.Event{
				{
					Name: "Event 1",
					Date: time.Date(2021, 1, 1, 9, 0, 0, 0, time.UTC),
					StartTime: "09:00",
					EndTime: "10:00",
				},
				{
					Name: "Event 2",
					Date: time.Date(2021, 1, 1, 9, 0, 0, 0, time.UTC),
					StartTime: "10:00",
					EndTime: "11:00",
				},
				{
					Name: "Event 3",
					Date: time.Date(2021, 1, 1, 9, 0, 0, 0, time.UTC),
					StartTime: "11:00",
					EndTime: "12:00",
				},
			},
		},
		{
			name: "Test 2",
			events: []*models.Event{
				{
					Name: "Event 1",
					Date: time.Date(2021, 1, 1, 9, 0, 0, 0, time.UTC),
					StartTime: "11:00",
					EndTime: "12:00",
				},
				{
					Name: "Event 2",
					Date: time.Date(2021, 1, 1, 9, 0, 0, 0, time.UTC),
					StartTime: "10:00",
					EndTime: "11:00",
				},
				{
					Name: "Event 3",
					Date: time.Date(2021, 1, 1, 9, 0, 0, 0, time.UTC),
					StartTime: "09:00",
					EndTime: "10:00",
				},
			},
			expected: []*models.Event{
				{
					Name: "Event 3",
					Date: time.Date(2021, 1, 1, 9, 0, 0, 0, time.UTC),
					StartTime: "09:00",
					EndTime: "10:00",
				},
				{
					Name: "Event 2",
					Date: time.Date(2021, 1, 1, 9, 0, 0, 0, time.UTC),
					StartTime: "10:00",
					EndTime: "11:00",
				},
				{
					Name: "Event 1",
					Date: time.Date(2021, 1, 1, 9, 0, 0, 0, time.UTC),
					StartTime: "11:00",
					EndTime: "12:00",
				},
			},
		},
		{
			name: "Test 3",
			events: []*models.Event{
				{
					Name: "Event 1",
					Date: time.Date(2021, 1, 1, 9, 0, 0, 0, time.UTC),
					StartTime: "11:00",
					EndTime: "12:00",
				},
				{
					Name: "Event 2",
					Date: time.Date(2021, 1, 1, 9, 0, 0, 0, time.UTC),
					StartTime: "09:00",
					EndTime: "10:00",
				},
				{
					Name: "Event 3",
					Date: time.Date(2021, 1, 1, 9, 0, 0, 0, time.UTC),
					StartTime: "10:00",
					EndTime: "11:00",
				},
			},
			expected: []*models.Event{
				{
					Name: "Event 2",
					Date: time.Date(2021, 1, 1, 9, 0, 0, 0, time.UTC),
					StartTime: "09:00",
					EndTime: "10:00",
				},
				{
					Name: "Event 3",
					Date: time.Date(2021, 1, 1, 9, 0, 0, 0, time.UTC),
					StartTime: "10:00",
					EndTime: "11:00",
				},
				{
					Name: "Event 1",
					Date: time.Date(2021, 1, 1, 9, 0, 0, 0, time.UTC),
					StartTime: "11:00",
					EndTime: "12:00",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			usecases.SortEventsByTime(tc.events)
			assert.Equal(t, tc.expected, tc.events)
		})
	}
}

func TestFormatTimeTo12Hour(t *testing.T) {
	type testCase struct {
		name string
		time24Str string
		expected string
	}

	testCases := []testCase{
		{
			name: "Test AM 1",
			time24Str: "09:10",
			expected: "9:10 AM",
		},
		{
			name: "Test AM 2",
			time24Str: "00:10",
			expected: "12:10 AM",
		},
		{
			name: "Test AM 3",
		time24Str: "11:00",
			expected: "11 AM",
		},
		{
			name: "Test PM 1",
			time24Str: "13:10",
			expected: "1:10 PM",
		},
		{
			name: "Test PM 2",
			time24Str: "12:10",
			expected: "12:10 PM",
		},
		{
			name: "Test PM 3",
			time24Str: "23:00",
			expected: "11 PM",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := usecases.FormatTimeTo12Hour(tc.time24Str)

			assert.Equal(t, tc.expected, actual)
		})
	}
}	
