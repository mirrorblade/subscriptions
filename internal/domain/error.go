package domain

import "errors"

var (
	ErrSubscriptionNotFound = errors.New("subscription was not found")
	ErrUserNotFound         = errors.New("user was not found")
	ErrInvalidID            = errors.New("id is not valid")
	ErrNoUpdateParameters   = errors.New("no update paramaters was chose")
	ErrInvalidPrice         = errors.New("price is not valid")
	ErrInvalidDate          = errors.New("date is not valid")
)
