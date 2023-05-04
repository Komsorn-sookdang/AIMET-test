package controllers

import (
	"aimet-test/internal/domains"
	"aimet-test/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type eventController struct {
	eventUsecase domains.EventUsecase
}

func NewEventController(eventUsecase domains.EventUsecase) domains.EventController {
	return &eventController{
		eventUsecase: eventUsecase,
	}
}

func (c *eventController) GetFilteredEvents(ctx *gin.Context) {
	var query models.GetEventQuery

	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.eventUsecase.ValidateGetEventQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	events, err := c.eventUsecase.GetFilteredEvents(&query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := models.NewGetEventResponse(events)

	ctx.JSON(http.StatusOK, response)
}
