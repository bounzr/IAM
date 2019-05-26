package repository

import "errors"

var (

	ErrClientNotFound       = errors.New("clientID not found")
	ErrGroupNotFound        = errors.New("groupID not found")
	ErrResourceNotAvailable = errors.New("resource not available")
	ErrResourceNotFound     = errors.New("resource not found")
	ErrUsernameNotAvailable = errors.New("username not available")
	ErrUsernameNotFound     = errors.New("username not found")
	ErrInvalidLogin         = errors.New("invalid login")

	//user repository errors
	ErrRepositoryNotAvailable = errors.New("repository not available")

	//SessionsStore errors
	ErrSessionNotFound = errors.New("session not found for user")
	ErrSessionInvalid  = errors.New("session not found for token")
)
