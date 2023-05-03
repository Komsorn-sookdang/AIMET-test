package usecases

import (
	"aimet-test/internal/domains"
	"aimet-test/internal/models"
	"sort"
	"strings"
	"time"
)

type eventUsecase struct {
	eventRepository domains.EventRepository
}

func NewEventUsecase(eventRepository domains.EventRepository) domains.EventUsecase {
	return &eventUsecase{
		eventRepository: eventRepository,
	}
}

func (u *eventUsecase) GetFilteredEvents(filter *models.GetEventQuery) ([]*models.Event, error) {
	if filter.Year == 0 {
		filter.Year = time.Now().Year()
	}

	allEvents, err := u.eventRepository.FindByMonth(filter.Month, filter.Year)
	if err != nil {
		return nil, err
	}

	if filter.SortBy == "time" {
		SortEventsByTime(allEvents)
	} 

	var filteredEvents []*models.Event
	for _, event := range allEvents {
		if filter.Day != 0 && event.Date.Day() != filter.Day {
			continue
		}
		if filter.Keyword != "" && !strings.Contains(event.Name, filter.Keyword) {
			continue
		}
		event.StartTime = FormatTimeTo12Hour(event.StartTime)
		event.EndTime = FormatTimeTo12Hour(event.EndTime)
		filteredEvents = append(filteredEvents, event)
	}

	return filteredEvents, nil
}

func ValidateGetEventQuery(q *models.GetEventQuery) error {

	if  q.SortBy != "" && q.SortBy != "date" && q.SortBy != "time" {
		return ErrInvalidSortBy
	}

	// Check if Month field is valid
	if q.Month < 1 || q.Month > 12 {
		return ErrInvalidMonth
	}

	// Check if Year field is valid
	if q.Year < 0 || q.Year > 9999 {
		return ErrInvalidYear
	}

	if q.Day != 0 {
		// Check if Date is valid
		daysInMonth := time.Date(q.Year, time.Month(q.Month+1), 0, 0, 0, 0, 0, time.UTC).Day()
		if q.Day < 0 || q.Day > daysInMonth {
			return ErrInvalidDate
		}
	}

	return nil
}

func SortEventsByTime(events []*models.Event) {
	// Sort the events by StartTime
	sort.Slice(events, func(i, j int) bool {
		return events[i].StartTime < events[j].StartTime
	})
}

func FormatTimeTo12Hour(time24Str string) string {
	time24, _ := time.Parse("15:04", time24Str)

	var time12Str string

	if time24.Minute() == 0 {
		time12Str = time24.Format("3 PM")
	} else {
		time12Str = time24.Format("3:04 PM")
	}

	return time12Str
}
