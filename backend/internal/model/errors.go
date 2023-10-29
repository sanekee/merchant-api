package model

import (
	"errors"
)

var (
	ErrNoResults = errors.New("error: no results")
	ErrExists    = errors.New("error: record exists")
	ErrRequest   = errors.New("error: bad request")
	ErrServer    = errors.New("error: internal server error")
	ErrDuplicate = errors.New("error: duplicated record")
)
