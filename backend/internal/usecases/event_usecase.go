package usecases

import (
	"aimet-test/internal/domains"
	"aimet-test/internal/models"
	"strings"
)

type eventUsecase struct {
	eventRepository domains.EventRepository
}

func NewEventUsecase(eventRepository domains.EventRepository) domains.EventUsecase {
	return &eventUsecase{
		eventRepository: eventRepository,
	}
}

func (u *eventUsecase) GetFilteredEvents(filter models.GetEventQuery) ([]*models.Event, error) {
	allEvents, err := u.eventRepository.FindByMonth(filter.Month, filter.Year)
	if err != nil {
		return nil, err
	}

	if filter.SortBy == "time" {
		models.SortEventsByTime(allEvents)
	} 

	var filteredEvents []*models.Event
	for _, event := range allEvents {
		if filter.Day != 0 && event.Date.Day() != filter.Day {
			continue
		}
		if filter.Keyword != "" && !strings.Contains(event.Name, filter.Keyword) {
			continue
		}
		event.ConvertTimeTo12HourFormat()
		filteredEvents = append(filteredEvents, event)
	}

	return filteredEvents, nil
}
