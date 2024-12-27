package constants

import "errors"

var (
	ErrInvalidKey   = errors.New("Key is not valid")
	ErrInvalidValue = errors.New("Value is invalid")
	ErrKeyNotFound  = errors.New("Key not found")
	ErrDuplicateKey = errors.New("Key with same name is already present")
	ErrInternal     = errors.New("Some internal error occurred")
)
