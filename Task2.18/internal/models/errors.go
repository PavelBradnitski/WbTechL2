package models

import "errors"

var (
	// ErrEventNotFound is returned when the event with the given ID does not exist.
	ErrEventNotFound = errors.New("event not found")

	// ErrEventDoesNotBelongToUser is returned when the event exists but belongs to another user.
	ErrEventDoesNotBelongToUser = errors.New("event does not belong to the user")

	// ErrInvalidEventDateFormat is returned when the provided date does not match YYYY-MM-DD format.
	ErrInvalidEventDateFormat = errors.New("invalid event date format, expected YYYY-MM-DD")

	// ErrEventCreationFailed is returned when the event could not be created.
	ErrEventCreationFailed = errors.New("failed to create event")

	// ErrEventUpdateFailed is returned when the event could not be updated.
	ErrEventUpdateFailed = errors.New("failed to update event")

	// ErrEventDeletionFailed is returned when the event could not be deleted.
	ErrEventDeletionFailed = errors.New("failed to delete event")

	// ErrUserIDRequired is returned when user_id is not provided.
	ErrUserIDRequired = errors.New("user_id must be provided")

	// ErrDateRequired is returned when date is not provided.
	ErrDateRequired = errors.New("date must be provided")

	// ErrInvalidUserID is returned when user_id cannot be converted to a number.
	ErrInvalidUserID = errors.New("fail to convert user_id into a number")

	// ErrFailedToRetrieveEvents is returned when events could not be retrieved.
	ErrFailedToRetrieveEvents = errors.New("failed to retrieve events")
)
