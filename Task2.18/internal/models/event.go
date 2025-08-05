package models

// Event структура данных для события
type Event struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Date   string `json:"date"`
	Event  string `json:"event"`
}
