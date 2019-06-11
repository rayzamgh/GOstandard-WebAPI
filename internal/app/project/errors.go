package project

import (
	"github.com/pkg/errors"
)

// Errors
var (
	ErrNotFound            = errors.New("entity not found")
	ErrInvalidField        = errors.New("invalid field")
	ErrMapKeyDoesNotExist  = errors.New("key does not exist")
	ErrUnknownMapValueType = errors.New("unknown map value data type")
)
