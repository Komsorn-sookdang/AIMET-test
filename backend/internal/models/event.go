package models

import (
	"time"
)

type GetEventQuery struct {
	Month   int    `form:"month" binding:"required"`
	Year    int    `form:"year"`
	Day     int    `form:"day"`
	SortBy  string `form:"sort_by"`
	Keyword string `form:"keyword"`
}

type Event struct {
	Name      string    `json:"name" bson:"name"`
	Date      time.Time `json:"date" bson:"date"`
	StartTime string    `json:"start_time" bson:"start_time"`
	EndTime   string    `json:"end_time" bson:"end_time"`
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
