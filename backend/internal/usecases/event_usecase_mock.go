package usecases

import (
	"aimet-test/internal/models"

	"github.com/stretchr/testify/mock"
)

type eventUsecaseMock struct {
	mock.Mock
}

func NewEventUsecaseMock() *eventUsecaseMock {
	return &eventUsecaseMock{}
}

func (u *eventUsecaseMock) GetFilteredEvents(filter *models.GetEventQuery) ([]*models.Event, error) {
	args := u.Called(filter)
	return args.Get(0).([]*models.Event), args.Error(1)
}

func (u *eventUsecaseMock) ValidateGetEventQuery(q *models.GetEventQuery) error {
	args := u.Called(q)
	return args.Error(0)
}
