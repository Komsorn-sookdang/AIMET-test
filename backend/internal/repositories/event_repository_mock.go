package repositories

import (
	"aimet-test/internal/models"

	"github.com/stretchr/testify/mock"
)

type eventRepositoryMock struct {
	mock.Mock
}

func NewEventRepositoryMock() *eventRepositoryMock {
	return &eventRepositoryMock{}
}

func (m *eventRepositoryMock) FindByMonth(month int, year int) ([]*models.Event, error) {
	args := m.Called(month, year)
	return args.Get(0).([]*models.Event), args.Error(1)
}