package domains

import "aimet-test/internal/models"

type EventRepository interface {
	FindByMonth(month int, year int) ([]*models.Event, error)
}
