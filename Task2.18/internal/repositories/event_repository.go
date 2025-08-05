package repositories

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/PavelBradnitski/WbTechL2/internal/models"
)

type EventRepository struct {
	events map[int]models.Event
	nextID int
	mu     sync.RWMutex
}

func NewEventRepo() *EventRepository {
	return &EventRepository{
		events: make(map[int]models.Event),
		nextID: 1,
		mu:     sync.RWMutex{},
	}
}

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
	fmt.Println(er.events)
	return newEvent.ID, nil
}

func (ed *EventRepository) GetEvent(ctx context.Context, userID int) (models.Event, error) {
	return models.Event{}, nil
}
func (ed *EventRepository) GetEventsForDay(ctx context.Context, userID int, date time.Time) ([]models.Event, error) {
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endDate := startDate.AddDate(0, 0, 1)

	return ed.GetEventsByUserIDAndDateRange(userID, startDate, endDate)
}

func (ed *EventRepository) GetEventsForWeek(ctx context.Context, userID int, date time.Time) ([]models.Event, error) {
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endDate := startDate.AddDate(0, 0, 7)

	return ed.GetEventsByUserIDAndDateRange(userID, startDate, endDate)
}

func (ed *EventRepository) GetEventsForMonth(ctx context.Context, userID int, date time.Time) ([]models.Event, error) {
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endDate := startDate.AddDate(0, 1, 0)

	return ed.GetEventsByUserIDAndDateRange(userID, startDate, endDate)
}

// GetEventsByUserIDAndDateRange возвращает события для указанного user_id в заданном диапазоне дат.
func (ed *EventRepository) GetEventsByUserIDAndDateRange(userID int, startDate, endDate time.Time) ([]models.Event, error) {
	ed.mu.RLock()
	defer ed.mu.RUnlock()

	filteredEvents := []models.Event{}
	for _, event := range ed.events {
		eventTime, err := time.Parse(time.DateOnly, event.Date)
		if err != nil {
			return nil, fmt.Errorf("error parsing date %s: %v", event.Date, err)
		}
		if !eventTime.Before(startDate) && !eventTime.After(endDate) {
			filteredEvents = append(filteredEvents, event)
		}
	}
	if len(filteredEvents) == 0 {
		// return nil, fmt.Errorf("no events")
		return nil, nil
	}
	return filteredEvents, nil
}
func (r *EventRepository) Update(ctx context.Context, Event *models.Event) error {
	return nil
}

func (r *EventRepository) DeleteByGroupAndEvent(ctx context.Context, group, Event string) error {
	return nil
}
