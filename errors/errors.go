package powErrors

import (
	"errors"
)

var (
	ErrIO                = errors.New("I/O error")
	ErrTaskInterruption  = errors.New("task interruption")
	ErrValidationFailure = errors.New("validation failure")
)
