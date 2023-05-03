package domains

import "aimet-test/internal/models"

type EventUsecase interface {
	GetFilteredEvents(filter *models.GetEventQuery) ([]*models.Event, error)
}
