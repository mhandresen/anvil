package engine

import "errors"

var (
	ErrNoEngine       = errors.New("no container engine")
	ErrNotFound       = errors.New("not found")
	ErrNotImplemented = errors.New("not implemented")
)