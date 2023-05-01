package domains

import "github.com/gin-gonic/gin"

type EventController interface {
	GetFilteredEvents(ctx *gin.Context)
}
