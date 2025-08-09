package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/PavelBradnitski/WbTechL2/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockEventService struct {
	mock.Mock
}

func (m *MockEventService) CreateEvent(ctx context.Context, event *models.Event) (int, error) {
	args := m.Called(ctx, event)
	return args.Get(0).(int), args.Error(1)
}

func (m *MockEventService) GetEventsForDay(ctx context.Context, userID int, date time.Time) ([]models.Event, error) {
	args := m.Called(ctx, userID, date)
	return args.Get(0).([]models.Event), args.Error(1)
}

func (m *MockEventService) GetEventsForWeek(ctx context.Context, userID int, date time.Time) ([]models.Event, error) {
	args := m.Called(ctx, userID, date)
	return args.Get(0).([]models.Event), args.Error(1)
}

func (m *MockEventService) GetEventsForMonth(ctx context.Context, userID int, date time.Time) ([]models.Event, error) {
	args := m.Called(ctx, userID, date)
	return args.Get(0).([]models.Event), args.Error(1)
}

func (m *MockEventService) UpdateEventByUser(ctx context.Context, event *models.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}
func (m *MockEventService) DeleteEventByUser(ctx context.Context, userID, id int) error {
	args := m.Called(ctx, userID, id)
	return args.Error(0)
}

func TestCreateEvent(t *testing.T) {
	event := models.Event{UserID: 53, Date: "2025-08-09", Event: "test event"}

	tests := []struct {
		name         string
		requestBody  interface{}
		mockSetup    func(m *MockEventService)
		expectedCode int
	}{
		{
			name:        "success",
			requestBody: event,
			mockSetup: func(m *MockEventService) {
				m.On("CreateEvent", mock.Anything, &event).Return(int(1), nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid JSON",
			requestBody:  "{bad json}",
			mockSetup:    func(m *MockEventService) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid date",
			requestBody:  models.Event{UserID: 53, Date: "09-08-2025", Event: "bad date"},
			mockSetup:    func(m *MockEventService) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "service error",
			requestBody: event,
			mockSetup: func(m *MockEventService) {
				m.On("CreateEvent", mock.Anything, &event).Return(int(0), errors.New("fail"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockService := newHandlerWithMock()
			tt.mockSetup(mockService)

			resp := performRequest(t, handler.CreateEvent, http.MethodPost, "/Events/create_event", tt.requestBody)
			checkResponseAndMocks(t, resp, tt.expectedCode, mockService)
		})
	}
}

func TestGetEvents(t *testing.T) {
	userID := 53
	date := "2025-08-09"
	parsedDate, _ := time.Parse(time.DateOnly, date)
	mockEvents := []models.Event{
		{UserID: userID, Date: date, Event: "event1"},
		{UserID: userID, Date: date, Event: "event2"},
	}

	tests := []struct {
		name         string
		url          string
		setupMock    func(m *MockEventService)
		expectedCode int
		queryParams  string
	}{
		{
			name: "GetEventsForDay success",
			url:  "/Events/events_for_day",
			setupMock: func(m *MockEventService) {
				m.On("GetEventsForDay", mock.Anything, userID, parsedDate).Return(mockEvents, nil)
			},
			expectedCode: http.StatusOK,
			queryParams:  "?user_id=53&date=2025-08-09",
		},
		{
			name:         "missing user_id for day",
			url:          "/Events/events_for_day",
			setupMock:    func(m *MockEventService) {},
			expectedCode: http.StatusBadRequest,
			queryParams:  "?date=2025-08-09",
		},
		{
			name:         "missing date for day",
			url:          "/Events/events_for_day",
			setupMock:    func(m *MockEventService) {},
			expectedCode: http.StatusBadRequest,
			queryParams:  "?user_id=53",
		},
		{
			name:         "invalid user_id for day",
			url:          "/Events/events_for_day",
			setupMock:    func(m *MockEventService) {},
			expectedCode: http.StatusBadRequest,
			queryParams:  "?user_id=abc&date=2025-08-09",
		},
		{
			name:         "invalid date for day",
			url:          "/Events/events_for_day",
			setupMock:    func(m *MockEventService) {},
			expectedCode: http.StatusBadRequest,
			queryParams:  "?user_id=53&date=09-08-2025",
		},
		{
			name: "service error for day",
			url:  "/Events/events_for_day",
			setupMock: func(m *MockEventService) {
				m.On("GetEventsForDay", mock.Anything, userID, parsedDate).Return([]models.Event{}, errors.New("fail"))
			},
			expectedCode: http.StatusInternalServerError,
			queryParams:  "?user_id=53&date=2025-08-09",
		},

		{
			name: "GetEventsForWeek success",
			url:  "/Events/events_for_week",
			setupMock: func(m *MockEventService) {
				m.On("GetEventsForWeek", mock.Anything, userID, parsedDate).Return(mockEvents, nil)
			},
			expectedCode: http.StatusOK,
			queryParams:  "?user_id=53&date=2025-08-09",
		},
		{
			name:         "missing user_id for week",
			url:          "/Events/events_for_week",
			setupMock:    func(m *MockEventService) {},
			expectedCode: http.StatusBadRequest,
			queryParams:  "?date=2025-08-09",
		},
		{
			name:         "missing date for week",
			url:          "/Events/events_for_week",
			setupMock:    func(m *MockEventService) {},
			expectedCode: http.StatusBadRequest,
			queryParams:  "?user_id=53",
		},
		{
			name:         "invalid user_id for week",
			url:          "/Events/events_for_week",
			setupMock:    func(m *MockEventService) {},
			expectedCode: http.StatusBadRequest,
			queryParams:  "?user_id=abc&date=2025-08-09",
		},
		{
			name:         "invalid date for week",
			url:          "/Events/events_for_week",
			setupMock:    func(m *MockEventService) {},
			expectedCode: http.StatusBadRequest,
			queryParams:  "?user_id=53&date=09-08-2025",
		},
		{
			name: "service error for week",
			url:  "/Events/events_for_week",
			setupMock: func(m *MockEventService) {
				m.On("GetEventsForWeek", mock.Anything, userID, parsedDate).Return([]models.Event{}, errors.New("fail"))
			},
			expectedCode: http.StatusInternalServerError,
			queryParams:  "?user_id=53&date=2025-08-09",
		},

		{
			name: "GetEventsForMonth success",
			url:  "/Events/events_for_month",
			setupMock: func(m *MockEventService) {
				m.On("GetEventsForMonth", mock.Anything, userID, parsedDate).Return(mockEvents, nil)
			},
			expectedCode: http.StatusOK,
			queryParams:  "?user_id=53&date=2025-08-09",
		},
		{
			name:         "missing user_id for month",
			url:          "/Events/events_for_month",
			setupMock:    func(m *MockEventService) {},
			expectedCode: http.StatusBadRequest,
			queryParams:  "?date=2025-08-09",
		},
		{
			name:         "missing date for month",
			url:          "/Events/events_for_month",
			setupMock:    func(m *MockEventService) {},
			expectedCode: http.StatusBadRequest,
			queryParams:  "?user_id=53",
		},
		{
			name:         "invalid user_id for month",
			url:          "/Events/events_for_month",
			setupMock:    func(m *MockEventService) {},
			expectedCode: http.StatusBadRequest,
			queryParams:  "?user_id=abc&date=2025-08-09",
		},
		{
			name:         "invalid date for month",
			url:          "/Events/events_for_month",
			setupMock:    func(m *MockEventService) {},
			expectedCode: http.StatusBadRequest,
			queryParams:  "?user_id=53&date=09-08-2025",
		},
		{
			name: "service error for month",
			url:  "/Events/events_for_month",
			setupMock: func(m *MockEventService) {
				m.On("GetEventsForMonth", mock.Anything, userID, parsedDate).Return([]models.Event{}, errors.New("fail"))
			},
			expectedCode: http.StatusInternalServerError,
			queryParams:  "?user_id=53&date=2025-08-09",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockService := newHandlerWithMock()
			tt.setupMock(mockService)

			router := gin.New()
			// Зарегистрируем все три роутера с соответствующими хендлерами
			router.GET("/Events/events_for_day", handler.GetEventsForDay)
			router.GET("/Events/events_for_week", handler.GetEventsForWeek)
			router.GET("/Events/events_for_month", handler.GetEventsForMonth)

			req, err := http.NewRequest(http.MethodGet, tt.url+tt.queryParams, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestUpdateEvent(t *testing.T) {
	event := models.Event{UserID: 53, Date: "2025-08-09", Event: "updated event"}

	tests := []struct {
		name         string
		requestBody  interface{}
		mockSetup    func(m *MockEventService)
		expectedCode int
	}{
		{
			name:        "success",
			requestBody: event,
			mockSetup: func(m *MockEventService) {
				m.On("UpdateEventByUser", mock.Anything, &event).Return(nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid JSON",
			requestBody:  "{bad json}",
			mockSetup:    func(m *MockEventService) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid date format",
			requestBody:  models.Event{UserID: 53, Date: "09-08-2025", Event: "bad date"},
			mockSetup:    func(m *MockEventService) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "event not found",
			requestBody: event,
			mockSetup: func(m *MockEventService) {
				m.On("UpdateEventByUser", mock.Anything, &event).Return(models.ErrEventNotFound)
			},
			expectedCode: http.StatusServiceUnavailable,
		},
		{
			name:        "event does not belong to user",
			requestBody: event,
			mockSetup: func(m *MockEventService) {
				m.On("UpdateEventByUser", mock.Anything, &event).Return(models.ErrEventDoesNotBelongToUser)
			},
			expectedCode: http.StatusServiceUnavailable,
		},
		{
			name:        "unknown error",
			requestBody: event,
			mockSetup: func(m *MockEventService) {
				m.On("UpdateEventByUser", mock.Anything, &event).Return(errors.New("fail"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockService := newHandlerWithMock()
			tt.mockSetup(mockService)

			resp := performRequest(t, handler.UpdateEvent, http.MethodPost, "/Events/update_event", tt.requestBody)
			checkResponseAndMocks(t, resp, tt.expectedCode, mockService)
		})
	}
}

func TestDeleteEvent(t *testing.T) {
	reqBody := map[string]int{"user_id": 53, "id": 10}

	tests := []struct {
		name         string
		requestBody  interface{}
		mockSetup    func(m *MockEventService)
		expectedCode int
	}{
		{
			name:        "success",
			requestBody: reqBody,
			mockSetup: func(m *MockEventService) {
				m.On("DeleteEventByUser", mock.Anything, 53, 10).Return(nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid JSON",
			requestBody:  "{bad json}",
			mockSetup:    func(m *MockEventService) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "event not found",
			requestBody: reqBody,
			mockSetup: func(m *MockEventService) {
				m.On("DeleteEventByUser", mock.Anything, 53, 10).Return(models.ErrEventNotFound)
			},
			expectedCode: http.StatusServiceUnavailable,
		},
		{
			name:        "event does not belong to user",
			requestBody: reqBody,
			mockSetup: func(m *MockEventService) {
				m.On("DeleteEventByUser", mock.Anything, 53, 10).Return(models.ErrEventDoesNotBelongToUser)
			},
			expectedCode: http.StatusServiceUnavailable,
		},
		{
			name:        "unknown error",
			requestBody: reqBody,
			mockSetup: func(m *MockEventService) {
				m.On("DeleteEventByUser", mock.Anything, 53, 10).Return(errors.New("fail"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockService := newHandlerWithMock()
			tt.mockSetup(mockService)

			resp := performRequest(t, handler.DeleteEvent, http.MethodPost, "/Events/delete_event", tt.requestBody)
			checkResponseAndMocks(t, resp, tt.expectedCode, mockService)
		})
	}
}

func performRequest(t *testing.T, handler gin.HandlerFunc, method, url string, body interface{}) *httptest.ResponseRecorder {
	router := gin.New()
	router.Handle(method, url, handler)

	var reqBody []byte
	var err error
	if body != nil {
		switch v := body.(type) {
		case string:
			reqBody = []byte(v)
		default:
			reqBody, err = json.Marshal(body)
			require.NoError(t, err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func checkResponseAndMocks(t *testing.T, resp *httptest.ResponseRecorder, expectedCode int, mockService *MockEventService) {
	assert.Equal(t, expectedCode, resp.Code)
	mockService.AssertExpectations(t)
}

func newHandlerWithMock() (*EventHandler, *MockEventService) {
	mockService := new(MockEventService)
	handler := &EventHandler{service: mockService}
	return handler, mockService
}
