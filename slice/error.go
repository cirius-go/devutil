package slice

import (
	"errors"
	"strings"
)

// ElemError represents an error related to slice elements.
type ElemError[In any] struct {
	Index int
	Value In
	Err   error
}

// Error implements the error interface for ElemError.
func (e *ElemError[In]) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return ""
}

// Unwrap returns the underlying error.
func (e *ElemError[In]) Unwrap() error {
	return e.Err
}

// SliceError represents an error related to slice operations.
type SliceError[In any] []*ElemError[In]

// Error implements the error interface for SliceError.
func (e SliceError[In]) Error() string {
	if len(e) == 0 {
		return ""
	}
	b := &strings.Builder{}
	for _, err := range e {
		b.WriteString(err.Error())
		b.WriteString("\n")
	}
	return b.String()
}

// Unwrap returns the underlying errors.
func (e SliceError[In]) Unwrap() error {
	if len(e) == 0 {
		return nil
	}
	var errs []error
	for _, err := range e {
		errs = append(errs, err.Err)
	}
	return errors.Join(errs...)
}

// At returns the error at the specified index.
func (e SliceError[In]) At(index int) error {
	if len(e) == 0 {
		return nil
	}
	if index < 0 || index >= len(e) {
		return nil
	}
	return (e)[index]
}

// OriginAt returns the original error at the specified index.
func (e SliceError[In]) OriginAt(index int) error {
	if len(e) == 0 {
		return nil
	}
	if index < 0 || index >= len(e) {
		return nil
	}
	return e[index].Err
}
