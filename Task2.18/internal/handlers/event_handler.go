package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/PavelBradnitski/WbTechL2/internal/models"
	"github.com/PavelBradnitski/WbTechL2/internal/services"

	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	service *services.EventService
}

func NewEventHandler(service *services.EventService) *EventHandler {
	return &EventHandler{service: service}
}

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

func (h *EventHandler) CreateEvent(c *gin.Context) {
	var event models.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := time.Parse(time.DateOnly, event.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	ctx := c.Request.Context()
	nextID, err := h.service.CreateEvent(ctx, &event)
	if err != nil {
		log.Printf("createEvent error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Event"})
		return
	}
	c.JSON(http.StatusCreated, nextID)
}

func (h *EventHandler) GetEventsForDay(c *gin.Context) {
	userIdStr := c.Query("user_id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fail to convert user_id into a number"})
		return
	}
	date := c.Query("date")
	d, err := time.Parse(time.DateOnly, date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}
	Events, err := h.service.GetEventsForDay(c, userId, d)
	// добавить обработку случая когда нету записей
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Events"})
		return
	}

	c.JSON(http.StatusOK, Events)
}

func (h *EventHandler) GetEventsForWeek(c *gin.Context) {
	userIdStr := c.Query("user_id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fail to convert user_id into a number"})
		return
	}
	date := c.Query("date")
	d, err := time.Parse(time.DateOnly, date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}
	Events, err := h.service.GetEventsForWeek(c, userId, d)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Events"})
		return
	}

	c.JSON(http.StatusOK, Events)
}

func (h *EventHandler) GetEventsForMonth(c *gin.Context) {
	userIdStr := c.Query("user_id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fail to convert user_id into a number"})
		return
	}
	date := c.Query("date")
	d, err := time.Parse(time.DateOnly, date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}
	Events, err := h.service.GetEventsForMonth(c, userId, d)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Events"})
		return
	}

	c.JSON(http.StatusOK, Events)
}

func (h *EventHandler) UpdateEvent(c *gin.Context) {
	userIdStr := c.Param("user_id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "fail to convert user_id into a number"})
		return
	}
	date := c.Param("date")

	var updatedEvent models.Event
	if err := c.ShouldBindJSON(&updatedEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingEvent, err := h.service.GetEvent(c, userId, date)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	existingEvent.Event = updatedEvent.Event

	if err := h.service.UpdateEvent(c, &existingEvent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Event"})
		return
	}

	c.JSON(http.StatusOK, existingEvent)
}

func (h *EventHandler) DeleteEvent(c *gin.Context) {
	user_id := c.Param("user_id")
	date := c.Param("date")

	if err := h.service.DeleteEvent(c, user_id, date); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Event"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
