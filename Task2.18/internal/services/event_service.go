package services

import (
	"context"
	"time"

	"github.com/PavelBradnitski/WbTechL2/internal/models"
	"github.com/PavelBradnitski/WbTechL2/internal/repositories"
)

type EventService interface {
	CreateEvent(ctx context.Context, event *models.Event) (int, error)
	GetEventsForDay(ctx context.Context, userID int, date time.Time) ([]models.Event, error)
	GetEventsForWeek(ctx context.Context, userID int, date time.Time) ([]models.Event, error)
	GetEventsForMonth(ctx context.Context, userID int, date time.Time) ([]models.Event, error)
	UpdateEventByUser(ctx context.Context, event *models.Event) error
	DeleteEventByUser(ctx context.Context, userID, id int) error
}

// EventService provides methods to manage events.
type eventService struct {
	repo *repositories.EventRepository
}

// NewEventService creates a new instance of EventService with the provided repository.
func NewEventService(repo *repositories.EventRepository) *eventService {
	return &eventService{repo: repo}
}

// CreateEvent creates a new event and returns its ID.
func (s *eventService) CreateEvent(ctx context.Context, Event *models.Event) (int, error) {
	return s.repo.Create(ctx, Event)
}

// GetEventsForDay retrieves events for a user within a specified date range.
func (s *eventService) GetEventsForDay(ctx context.Context, userID int, date time.Time) ([]models.Event, error) {
	return s.repo.GetEventsForDay(ctx, userID, date)
}

// GetEventsForWeek retrieves events for a user within a specified week.
func (s *eventService) GetEventsForWeek(ctx context.Context, userID int, date time.Time) ([]models.Event, error) {
	return s.repo.GetEventsForWeek(ctx, userID, date)
}

// GetEventsForMonth retrieves events for a user within a specified month.
func (s *eventService) GetEventsForMonth(ctx context.Context, userID int, date time.Time) ([]models.Event, error) {
	return s.repo.GetEventsForMonth(ctx, userID, date)
}

// UpdateEventByUser updates an existing event by user ID.
func (s *eventService) UpdateEventByUser(ctx context.Context, Event *models.Event) error {
	return s.repo.UpdateByUser(ctx, Event)
}

// DeleteEventByUser deletes an event by user ID and event ID.
func (s *eventService) DeleteEventByUser(ctx context.Context, userID, id int) error {
	return s.repo.DeleteByUser(ctx, userID, id)
}
