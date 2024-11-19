package customErrors

import "errors"

var (
	ErrInvalidInput     = errors.New("invalid input: missing required field")
	ErrJsonOpen         = errors.New("error opening JSON file")
	ErrJsonRead         = errors.New("error reading JSON file")
	ErrJsonWrite        = errors.New("error writing JSON file")
	ErrJsonUnmarshal    = errors.New("error unmarshalling Json")
	ErrJsonMarshal      = errors.New("error marshalling JSON")
	ErrExistConflict    = errors.New("already exist")
	ErrNotExistConflict = errors.New("doesn't exist")
	ErrOrderClosed      = errors.New("the order is already closed")
)
