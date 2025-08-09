package repositories

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/PavelBradnitski/WbTechL2/internal/models"
)

// EventRepository provides methods to interact with the event data store.
type EventRepository struct {
	events map[int]models.Event
	nextID int // Next ID to be assigned to a new event
	mu     sync.RWMutex
}

// NewEventRepo creates a new instance of EventRepository.
func NewEventRepo() *EventRepository {
	return &EventRepository{
		events: make(map[int]models.Event),
		nextID: 1,
		mu:     sync.RWMutex{},
	}
}

// Create adds a new event to the repository and returns its ID.
func (er *EventRepository) Create(ctx context.Context, event *models.Event) (int, error) {
	er.mu.Lock()
	defer er.mu.Unlock()
	newEvent := models.Event{
		ID:     er.nextID,
		UserID: event.UserID,
		Date:   event.Date,
		Event:  event.Event,
	}
	er.events[newEvent.ID] = newEvent
	er.nextID++
	return newEvent.ID, nil
}

// GetEventsForDay retrieves an event by user ID and date for a day.
func (er *EventRepository) GetEventsForDay(ctx context.Context, userID int, date time.Time) ([]models.Event, error) {
	er.mu.RLock()
	defer er.mu.RUnlock()
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endDate := startDate.AddDate(0, 0, 1)

	return er.GetEventsByUserIDAndDateRange(userID, startDate, endDate)
}

// GetEventsForWeek retrieves events for a user for a week starting from the given date.
func (er *EventRepository) GetEventsForWeek(ctx context.Context, userID int, date time.Time) ([]models.Event, error) {
	er.mu.RLock()
	defer er.mu.RUnlock()
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endDate := startDate.AddDate(0, 0, 7)

	return er.GetEventsByUserIDAndDateRange(userID, startDate, endDate)
}

// GetEventsForMonth retrieves events for a user for a month starting from the given date.
func (er *EventRepository) GetEventsForMonth(ctx context.Context, userID int, date time.Time) ([]models.Event, error) {
	er.mu.RLock()
	defer er.mu.RUnlock()
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endDate := startDate.AddDate(0, 1, 0)

	return er.GetEventsByUserIDAndDateRange(userID, startDate, endDate)
}

// GetEventsByUserIDAndDateRange retrieves events for a user within a specified date range.
func (er *EventRepository) GetEventsByUserIDAndDateRange(userID int, startDate, endDate time.Time) ([]models.Event, error) {
	er.mu.RLock()
	defer er.mu.RUnlock()

	filteredEvents := []models.Event{}
	for _, event := range er.events {
		if userID == event.UserID {
			eventTime, err := time.Parse(time.DateOnly, event.Date)
			if err != nil {
				return nil, fmt.Errorf("error parsing date %s: %v", event.Date, err)
			}
			if !eventTime.Before(startDate) && !eventTime.After(endDate) {
				filteredEvents = append(filteredEvents, event)
			}
		}
	}
	return filteredEvents, nil
}

// UpdateByUser updates an existing event by user ID.
func (er *EventRepository) UpdateByUser(ctx context.Context, event *models.Event) error {
	er.mu.Lock()
	defer er.mu.Unlock()

	existing, exists := er.events[event.ID]
	if !exists {
		return models.ErrEventNotFound
	}
	if existing.UserID != event.UserID {
		return models.ErrEventDoesNotBelongToUser
	}

	er.events[event.ID] = *event
	return nil
}

// DeleteByUser deletes an event by user ID and event ID.
func (er *EventRepository) DeleteByUser(ctx context.Context, userID, id int) error {
	er.mu.Lock()
	defer er.mu.Unlock()

	existing, exists := er.events[id]
	if !exists {
		return models.ErrEventNotFound
	}
	if existing.UserID != userID {
		return models.ErrEventDoesNotBelongToUser
	}

	delete(er.events, id)
	return nil
}
