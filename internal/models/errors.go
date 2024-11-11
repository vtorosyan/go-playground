package models

import "errors"

var (
	ErrNoRecord = errors.New("no Snippet with a given ID is found")

	ErrInvalidCredentials = errors.New("models: invalid credentials")

	ErrDuplicateEmail = errors.New("models: duplicate email")
)
