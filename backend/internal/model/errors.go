package model

import "fmt"

var (
	ErrNoResults = Errorf("error: no results")
	ErrExists    = Errorf("error: record exists")
	ErrRequest   = Errorf("error: bad request")
	ErrServer    = Errorf("error: internal server error")
	ErrDuplicate = Errorf("error: duplicated record")
)

type Error struct {
	s string
}

func Errorf(s string, args ...interface{}) Error {
	return Error{s: fmt.Sprintf(s, args...)}
}

func (err Error) Error() string {
	return err.s
}
