package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/PavelBradnitski/WbTechL2/internal/models"
	"github.com/PavelBradnitski/WbTechL2/internal/services"

	"github.com/gin-gonic/gin"
)

// EventHandler handles HTTP requests related to events.
type EventHandler struct {
	service *services.EventService
}

// APIResponse is a standard response structure for API responses.
type APIResponse struct {
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// NewEventHandler creates a new EventHandler with the provided service.
func NewEventHandler(service *services.EventService) *EventHandler {
	return &EventHandler{service: service}
}

// SuccessResponse creates a successful API response.
func SuccessResponse(result interface{}) APIResponse {
	return APIResponse{Result: result}
}

// ErrorResponse creates an error API response.
func ErrorResponse(err error) APIResponse {
	return APIResponse{Error: err.Error()}
}

// RegisterRoutes registers the event-related routes with the provided router.
func (h *EventHandler) RegisterRoutes(router *gin.Engine) {
	EventGroup := router.Group("/Events")
	{
		EventGroup.POST("/create_event", h.CreateEvent)
		EventGroup.POST("/update_event/", h.UpdateEvent)
		EventGroup.POST("/delete_event/", h.DeleteEvent)
		EventGroup.GET("/events_for_day", h.GetEventsForDay)
		EventGroup.GET("/events_for_week", h.GetEventsForWeek)
		EventGroup.GET("/events_for_month", h.GetEventsForMonth)
	}
}

// CreateEvent handles the creation of a new event.
func (h *EventHandler) CreateEvent(c *gin.Context) {
	var event models.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
	_, err := time.Parse(time.DateOnly, event.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(fmt.Errorf("invalid date format. Use YYYY-MM-DD")))
		return
	}

	id, err := h.service.CreateEvent(c, &event)
	if err != nil {
		log.Printf("createEvent error: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse(fmt.Errorf("failed to create Event")))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(gin.H{"id": id}))
}

// GetEventsForDay handles the retrieval of a specific event by user ID and date for one day.
func (h *EventHandler) GetEventsForDay(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusInternalServerError, ErrorResponse(fmt.Errorf("user_id must be entered")))
		return
	}
	date := c.Query("date")
	if date == "" {
		c.JSON(http.StatusInternalServerError, ErrorResponse(fmt.Errorf("date must be entered")))
		return
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(fmt.Errorf("fail to convert user_id into a number")))
		return
	}
	d, err := time.Parse(time.DateOnly, date)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(fmt.Errorf("invalid date format. Use YYYY-MM-DD")))
		return
	}
	events, err := h.service.GetEventsForDay(c, userID, d)
	if err != nil {
		log.Printf("getEventsForDay error: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse(fmt.Errorf("failed to retrieve Events")))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(events))
}

// GetEventsForWeek handles the retrieval of events for a specific user for a week.
func (h *EventHandler) GetEventsForWeek(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(fmt.Errorf("fail to convert user_id into a number")))
		return
	}
	dateStr := c.Query("date")
	date, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(fmt.Errorf("invalid date format. Use YYYY-MM-DD")))
		return
	}
	Events, err := h.service.GetEventsForWeek(c, userID, date)
	if err != nil {
		log.Printf("getEventsForWeek error: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse(fmt.Errorf("failed to retrieve Events")))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(Events))
}

// GetEventsForMonth handles the retrieval of events for a specific user for a month.
func (h *EventHandler) GetEventsForMonth(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(fmt.Errorf("fail to convert user_id into a number")))
		return
	}
	date := c.Query("date")
	d, err := time.Parse(time.DateOnly, date)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(fmt.Errorf("invalid date format. Use YYYY-MM-DD")))
		return
	}
	Events, err := h.service.GetEventsForMonth(c, userID, d)
	if err != nil {
		log.Printf("getEventsForMonth error: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse(fmt.Errorf("failed to retrieve Events")))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(Events))
}

// UpdateEvent handles the update of an existing event.
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	var event models.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	_, err := time.Parse(time.DateOnly, event.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(fmt.Errorf("invalid date format. Use YYYY-MM-DD")))
		return
	}

	ctx := c.Request.Context()
	err = h.service.UpdateEventByUser(ctx, &event)
	if err != nil {
		log.Printf("updateEvent error: %v", err)
		switch err.Error() {
		case "event not found":
			c.JSON(http.StatusServiceUnavailable, ErrorResponse(fmt.Errorf("event not found")))
			return
		case "event does not belong to the user":
			c.JSON(http.StatusServiceUnavailable, ErrorResponse(fmt.Errorf("event does not belong to the user")))
			return
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse(fmt.Errorf("failed to update Event")))
			return
		}
	}
	c.JSON(http.StatusOK, SuccessResponse("updated"))
}

// DeleteEvent handles the deletion of an event by user ID and event ID.
func (h *EventHandler) DeleteEvent(c *gin.Context) {
	var req struct {
		UserID int `json:"user_id"`
		ID     int `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(fmt.Errorf("invalid request body: %v", err)))
		return
	}

	err := h.service.DeleteEventByUser(c, req.UserID, req.ID)
	if err != nil {
		log.Printf("deleteEvent error: %v", err)
		switch err.Error() {
		case "event not found":
			c.JSON(http.StatusServiceUnavailable, ErrorResponse(fmt.Errorf("event not found")))
			return
		case "event does not belong to the user":
			c.JSON(http.StatusServiceUnavailable, ErrorResponse(fmt.Errorf("event does not belong to the user")))
			return
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse(fmt.Errorf("failed to delete Event")))
			return
		}
	}

	c.JSON(http.StatusOK, SuccessResponse("deleted"))
}
