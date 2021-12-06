package compress

import "errors"

var (
	// ErrInvalidResponse uses for invalid unmarshal response from worker
	ErrInvalidResponse = errors.New("invalid response from worker")
	// ErrCompressWorker uses where compress worker got an error
	ErrCompressWorker = errors.New("compress worker got an error")
	// ErrInvalidTypeAssertion uses when interface can't transform for needed type
	ErrInvalidTypeAssertion = errors.New("invalid type assertion")
	// ErrInvalidID uses when try updating record with id <= 0
	ErrInvalidID = errors.New("invalid id")
)
