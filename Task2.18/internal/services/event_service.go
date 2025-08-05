package services

import (
	"context"
	"time"

	"github.com/PavelBradnitski/WbTechL2/internal/models"
	"github.com/PavelBradnitski/WbTechL2/internal/repositories"
)

type EventService struct {
	repo *repositories.EventRepository
}

func NewEventService(repo *repositories.EventRepository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) CreateEvent(ctx context.Context, Event *models.Event) (int, error) {
	return s.repo.Create(ctx, Event)
}
func (s *EventService) GetEvent(ctx context.Context, user_id int, date string) (models.Event, error) {
	return s.repo.GetEvent(ctx, user_id)
}
func (s *EventService) GetEventsForDay(ctx context.Context, user_id int, date time.Time) ([]models.Event, error) {
	return s.repo.GetEventsForDay(ctx, user_id, date)
}

func (s *EventService) GetEventsForWeek(ctx context.Context, user_id int, date time.Time) ([]models.Event, error) {
	return s.repo.GetEventsForWeek(ctx, user_id, date)
}

func (s *EventService) GetEventsForMonth(ctx context.Context, user_id int, date time.Time) ([]models.Event, error) {
	return s.repo.GetEventsForMonth(ctx, user_id, date)
}

func (s *EventService) UpdateEvent(ctx context.Context, Event *models.Event) error {
	return s.repo.Update(ctx, Event)
}

func (s *EventService) DeleteEvent(ctx context.Context, group, Event string) error {
	return s.repo.DeleteByGroupAndEvent(ctx, group, Event)
}
