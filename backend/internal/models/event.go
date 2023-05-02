package models

import (
	"fmt"
	"sort"
	"time"
)

type GetEventQuery struct {
	Month   int    `form:"month" binding:"required"`
	Year    int    `form:"year"`
	Day     int    `form:"day"`
	SortBy  string `form:"sort_by"`
	Keyword string `form:"keyword"`
}

func (q *GetEventQuery) Validate() error {
	// Set SortBy to "date" if not provided
	if q.SortBy == "" {
		q.SortBy = "date"
	}

	// Check if SortBy field is valid
	if q.SortBy != "date" && q.SortBy != "time" {
		return fmt.Errorf("Invalid sort_by: %s", q.SortBy)
	}

	// Check if Month field is valid
	if q.Month < 1 || q.Month > 12 {
		return fmt.Errorf("Invalid month: %d", q.Month)
	}

	// Set Year to current year if not provided
	if q.Year == 0 {
		now := time.Now()
		q.Year = now.Year()
	}

	// Check if Year field is valid
	if q.Year < 0 || q.Year > 9999 {
		return fmt.Errorf("Invalid year: %d", q.Year)
	}

	if q.Day != 0 {
		// Check if Date is valid
		daysInMonth := time.Date(q.Year, time.Month(q.Month+1), 0, 0, 0, 0, 0, time.UTC).Day()
		if q.Day < 0 || q.Day > daysInMonth {
			return fmt.Errorf("Invalid day: %d for month: %d, year: %d", q.Day, q.Month, q.Year)
		}
	}

	return nil
}

type Event struct {
	Name      string    `json:"name" bson:"name"`
	Date      time.Time `json:"date" bson:"date"`
	StartTime string    `json:"start_time" bson:"start_time"`
	EndTime   string    `json:"end_time" bson:"end_time"`
}

func (e *Event) ConvertTimeTo12HourFormat() {
	startTime24, _ := time.Parse("15:04", e.StartTime)
	endTime24, _ := time.Parse("15:04", e.EndTime)

	if startTime24.Minute() == 0 {
		e.StartTime = startTime24.Format("3 PM")
	} else {
		e.StartTime = startTime24.Format("3:04 PM")
	}

	if endTime24.Minute() == 0 {
		e.EndTime = endTime24.Format("3 PM")
	} else {
		e.EndTime = endTime24.Format("3:04 PM")
	}
}

func SortEventsByTime(events []*Event) {
	// Sort the events by StartTime
	sort.Slice(events, func(i, j int) bool {
		return events[i].StartTime < events[j].StartTime
	})
}

type GetEventResponse struct {
	Count  int      `json:"count"`
	Events []*Event `json:"events"`
}

func NewGetEventResponse(events []*Event) *GetEventResponse {
	return &GetEventResponse{
		Count:  len(events),
		Events: events,
	}
}
